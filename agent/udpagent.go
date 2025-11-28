package agent

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenNHP/opennhp/nhp/core"
	"github.com/OpenNHP/opennhp/nhp/log"
)

type KnockUser struct {
	UserId         string
	OrganizationId string
	UserData       map[string]any
}

type KnockResource struct {
	AuthServiceId  string `json:"aspId"`
	ResourceId     string `json:"resId"`
	ServerHostname string `json:"serverHostname"`
	ServerIp       string `json:"serverIp"`
	ServerPort     int    `json:"serverPort"`
}

func (res *KnockResource) Id() string {
	return res.AuthServiceId + "/" + res.ResourceId
}

func (res *KnockResource) ServerHost() string {
	hostAddr := res.ServerIp
	if len(res.ServerHostname) > 0 {
		hostAddr = res.ServerHostname
	}
	if res.ServerPort == 0 {
		return hostAddr
	}
	return fmt.Sprintf("%s:%d", hostAddr, res.ServerPort)
}

type KnockTarget struct {
	sync.Mutex
	KnockResource
	ServerPeer           *core.UdpPeer
	LastKnockSuccessTime time.Time
}

func (kt *KnockTarget) GetServerPeer() *core.UdpPeer {
	kt.Lock()
	defer kt.Unlock()

	return kt.ServerPeer
}

type UdpAgent struct {
	stats struct {
		totalRecvBytes uint64
		totalSendBytes uint64
	}

	config *Config
	//log    *log.Logger

	remoteConnectionMutex sync.Mutex
	remoteConnectionMap   map[string]*UdpConn // indexed by remote UDP address

	knockTargetMapMutex sync.Mutex
	knockTargetMap      map[string]*KnockTarget // indexed by aspId + resId

	serverPeerMutex sync.Mutex
	serverPeerMap   map[string]*core.UdpPeer // indexed by server's public key

	device  *core.Device
	wg      sync.WaitGroup
	running atomic.Bool

	signals struct {
		stop                  chan struct{}
		knockTargetStop       chan struct{}
		knockTargetMapUpdated chan struct{}
	}

	recvMsgCh <-chan *core.PacketParserData
	sendMsgCh chan *core.MsgData

	// one agent should serve only one specific user at a time
	knockUserMutex sync.RWMutex
	knockUser      *KnockUser
	deviceId       string
	checkResults   map[string]any
}

type UdpConn struct {
	ConnData *core.ConnectionData
	netConn  *net.UDPConn
}

func (c *UdpConn) Close() {
	c.netConn.Close()
	c.ConnData.Close()
}

func (a *UdpAgent) Start() (err error) {
	err = a.loadAgentConfig()
	if err != nil {
		return err
	}

	prk, err := base64.StdEncoding.DecodeString(a.config.PrivateKeyBase64)
	if err != nil {
		log.Error("private key parse error %v\n", err)
		return fmt.Errorf("private key parse error %v", err)
	}

	a.device = core.NewDevice(core.NHP_AGENT, prk, nil)
	if a.device == nil {
		log.Critical("failed to create device %v\n", err)
		return fmt.Errorf("failed to create device %v", err)
	}

	// start device routines
	a.device.Start()

	// load peers
	a.loadPeers()

	a.remoteConnectionMap = make(map[string]*UdpConn)

	a.signals.stop = make(chan struct{})
	a.signals.knockTargetStop = make(chan struct{})
	a.signals.knockTargetMapUpdated = make(chan struct{}, 1)

	// load knock resources
	a.loadResources()

	a.recvMsgCh = a.device.DecryptedMsgQueue
	a.sendMsgCh = make(chan *core.MsgData, core.SendQueueSize)

	// start agent routines
	a.wg.Add(2)
	go a.sendMessageRoutine()
	go a.recvMessageRoutine()

	a.running.Store(true)

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func (a *UdpAgent) SetKnockUser(usrId string, orgId string, userData map[string]any) {
	a.knockUserMutex.Lock()
	a.knockUser.UserId = usrId
	a.knockUser.OrganizationId = orgId
	a.knockUser.UserData = userData
	a.knockUserMutex.Unlock()
}

func (a *UdpAgent) SetDeviceId(devId string) {
	a.deviceId = devId
}

func (a *UdpAgent) SetCheckResults(results map[string]any) {
	a.checkResults = results
}

// export Stop
func (a *UdpAgent) Stop() {
	a.running.Store(false)
	close(a.signals.knockTargetStop)
	close(a.signals.stop)
	a.device.Stop()
	a.StopConfigWatch()
	a.wg.Wait()
	close(a.sendMsgCh)
	close(a.signals.knockTargetMapUpdated)
}

func (a *UdpAgent) IsRunning() bool {
	return a.running.Load()
}

func (a *UdpAgent) newConnection(addr *net.UDPAddr) (conn *UdpConn) {
	conn = &UdpConn{}
	var err error
	// unlike tcp, udp dial is fast (just socket bind), so no need to run in a thread
	conn.netConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Error("could not connect to remote addr %s", addr.String())
		return nil
	}

	// retrieve local port
	laddr := conn.netConn.LocalAddr()
	localAddr, err := net.ResolveUDPAddr(laddr.Network(), laddr.String())
	if err != nil {
		log.Error("resolve local UDPAddr error %v\n", err)
		return nil
	}

	log.Info("Dial up new UDP connection from %s to %s", localAddr.String(), addr.String())

	conn.ConnData = &core.ConnectionData{
		Device:               a.device,
		CookieStore:          &core.CookieStore{},
		RemoteTransactionMap: make(map[uint64]*core.RemoteTransaction),
		LocalAddr:            localAddr,
		RemoteAddr:           addr,
		TimeoutMs:            DefaultConnectionTimeoutMs,
		SendQueue:            make(chan *core.Packet, PacketQueueSizePerConnection),
		RecvQueue:            make(chan *core.Packet, PacketQueueSizePerConnection),
		BlockSignal:          make(chan struct{}),
		SetTimeoutSignal:     make(chan struct{}),
		StopSignal:           make(chan struct{}),
	}

	conn.ConnData.Add(1)
	go a.recvPacketRoutine(conn)

	return conn
}

func (a *UdpAgent) sendMessageRoutine() {
	defer a.wg.Done()
	defer log.Info("sendMessageRoutine stopped")

	log.Info("sendMessageRoutine started")

	for {
		select {
		case <-a.signals.stop:
			return

		case md, ok := <-a.sendMsgCh:
			if !ok {
				return
			}
			if md == nil || md.RemoteAddr == nil {
				log.Warning("Invalid initiator session starter")
				continue
			}

			addrStr := md.RemoteAddr.String()

			a.remoteConnectionMutex.Lock()
			conn, found := a.remoteConnectionMap[addrStr]
			a.remoteConnectionMutex.Unlock()

			if found {
				md.ConnData = conn.ConnData
			} else {
				conn = a.newConnection(md.RemoteAddr)
				if conn == nil {
					log.Error("Failed to dial to remote address: %s", addrStr)
					continue
				}

				a.remoteConnectionMutex.Lock()
				a.remoteConnectionMap[addrStr] = conn
				a.remoteConnectionMutex.Unlock()

				md.ConnData = conn.ConnData

				// launch connection routine
				a.wg.Add(1)
				go a.connectionRoutine(conn)
			}

			a.device.SendMsgToPacket(md)
		}
	}

}

func (a *UdpAgent) SendPacket(pkt *core.Packet, conn *UdpConn) (n int, err error) {
	defer func() {
		atomic.AddUint64(&a.stats.totalSendBytes, uint64(n))
		atomic.StoreInt64(&conn.ConnData.LastLocalSendTime, time.Now().UnixNano())

		if !pkt.KeepAfterSend {
			a.device.ReleasePoolPacket(pkt)
		}
	}()

	pktType := core.HeaderTypeToString(pkt.HeaderType)
	//log.Debug("Send [%s] packet (%s -> %s): %+v", pktType, conn.ConnData.LocalAddr.String(), conn.ConnData.RemoteAddr.String(), pkt.Content)
	log.Info("Send [%s] packet (%s -> %s), %d bytes", pktType, conn.ConnData.LocalAddr.String(), conn.ConnData.RemoteAddr.String(), len(pkt.Content))
	log.Evaluate("Send [%s] packet (%s -> %s), %d bytes", pktType, conn.ConnData.LocalAddr.String(), conn.ConnData.RemoteAddr.String(), len(pkt.Content))
	return conn.netConn.Write(pkt.Content)
}

func (a *UdpAgent) recvPacketRoutine(conn *UdpConn) {
	addrStr := conn.ConnData.RemoteAddr.String()

	defer conn.ConnData.Done()
	defer log.Debug("recvPacketRoutine for %s stopped", addrStr)

	log.Debug("recvPacketRoutine for %s started", addrStr)

	for {
		select {
		case <-conn.ConnData.StopSignal:
			return

		default:
		}

		// udp recv, blocking until packet arrives or netConn.Close()
		pkt := a.device.AllocatePoolPacket()
		n, err := conn.netConn.Read(pkt.Buf[:])
		if err != nil {
			a.device.ReleasePoolPacket(pkt)
			if n == 0 {
				// udp connection closed, it is not an error
				return
			}
			log.Error("Failed to receive from remote address %s (%v)", addrStr, err)
			continue
		}

		// add total recv bytes
		atomic.AddUint64(&a.stats.totalRecvBytes, uint64(n))

		// check minimal length
		if n < pkt.MinimalLength() {
			a.device.ReleasePoolPacket(pkt)
			log.Error("Received UDP packet from %s is too short, discard", addrStr)
			continue
		}

		pkt.Content = pkt.Buf[:n]
		//log.Trace("receive udp packet (%s -> %s): %+v", conn.ConnData.RemoteAddr.String(), conn.ConnData.LocalAddr.String(), pkt.Content)

		typ, _, err := a.device.RecvPrecheck(pkt)
		msgType := core.HeaderTypeToString(typ)
		log.Info("Receive [%s] packet (%s -> %s), %d bytes", msgType, addrStr, conn.ConnData.LocalAddr.String(), n)
		log.Evaluate("Receive [%s] packet (%s -> %s), %d bytes", msgType, addrStr, conn.ConnData.LocalAddr.String(), n)
		if err != nil {
			a.device.ReleasePoolPacket(pkt)
			log.Warning("Receive [%s] packet (%s -> %s), precheck error: %v", msgType, addrStr, conn.ConnData.LocalAddr.String(), err)
			log.Evaluate("Receive [%s] packet (%s -> %s) precheck error: %v", msgType, addrStr, conn.ConnData.LocalAddr.String(), err)
			continue
		}

		atomic.StoreInt64(&conn.ConnData.LastLocalRecvTime, time.Now().UnixNano())

		conn.ConnData.ForwardInboundPacket(pkt)
	}
}

func (a *UdpAgent) connectionRoutine(conn *UdpConn) {
	addrStr := conn.ConnData.RemoteAddr.String()

	defer a.wg.Done()
	defer log.Debug("Connection routine: %s stopped", addrStr)

	log.Debug("Connection routine: %s started", addrStr)

	// stop receiving packets and clean up
	defer func() {
		a.remoteConnectionMutex.Lock()
		delete(a.remoteConnectionMap, addrStr)
		a.remoteConnectionMutex.Unlock()

		conn.Close()
	}()

	for {
		select {
		case <-a.signals.stop:
			return

		case <-conn.ConnData.SetTimeoutSignal:
			if conn.ConnData.TimeoutMs <= 0 {
				log.Debug("Connection routine closed immediately")
				return
			}

		case <-time.After(time.Duration(conn.ConnData.TimeoutMs) * time.Millisecond):
			// timeout, quit routine
			log.Debug("Connection routine idle timeout")
			return

		case pkt, ok := <-conn.ConnData.SendQueue:
			if !ok {
				return
			}
			if pkt == nil {
				continue
			}
			a.SendPacket(pkt, conn)

		case pkt, ok := <-conn.ConnData.RecvQueue:
			if !ok {
				return
			}
			if pkt == nil {
				continue
			}
			log.Debug("Received udp packet len [%d] from addr: %s\n", len(pkt.Content), addrStr)

			// process keepalive packet
			if pkt.HeaderType == core.NHP_KPL {
				a.device.ReleasePoolPacket(pkt)
				log.Info("Receive [NHP_KPL] message (%s -> %s)", addrStr, conn.ConnData.LocalAddr.String())
				continue
			}

			if a.device.IsTransactionResponse(pkt.HeaderType) {
				// forward to a specific transaction
				transactionId := pkt.Counter()
				transaction := a.device.FindLocalTransaction(transactionId)
				if transaction != nil {
					transaction.NextPacketCh <- pkt
					continue
				}
			}

			pd := &core.PacketData{
				BasePacket: pkt,
				ConnData:   conn.ConnData,
				InitTime:   atomic.LoadInt64(&conn.ConnData.LastLocalRecvTime),
			}
			// generic receive
			a.device.RecvPacketToMsg(pd)

		case <-conn.ConnData.BlockSignal:
			log.Critical("blocking address %s", addrStr)
			return
		}
	}
}

func (a *UdpAgent) recvMessageRoutine() {
	defer a.wg.Done()
	defer log.Info("recvMessageRoutine stopped")

	log.Info("recvMessageRoutine started")

	for {
		select {
		case <-a.signals.stop:
			return

		case ppd, ok := <-a.recvMsgCh:
			if !ok {
				return
			}
			if ppd == nil {
				continue
			}

			switch ppd.HeaderType {
			case core.NHP_COK:
				// synchronously block and deal with cookie message to ensure future messages will be correctly processed. note cookie is not handled as a transaction, so it arrives in here
				a.HandleCookieMessage(ppd)

			}
		}
	}
}

func (a *UdpAgent) FindServerPeerFromResource(res *KnockResource) *core.UdpPeer {
	a.serverPeerMutex.Lock()
	defer a.serverPeerMutex.Unlock()
	for _, peer := range a.serverPeerMap {
		if peer.Host() == res.ServerHost() {
			return peer
		}
	}

	return nil
}

func (a *UdpAgent) GetKnockTargetByResourceId(resId string) *KnockTarget {
	a.knockTargetMapMutex.Lock()
	a.knockTargetMapMutex.Unlock()
	if target, ok := a.knockTargetMap[resId]; ok {
		return target
	}
	return nil
}

func (a *UdpAgent) GetFirstKnockTarget() *KnockTarget {
	a.knockTargetMapMutex.Lock()
	a.knockTargetMapMutex.Unlock()
	if len(a.knockTargetMap) > 0 {
		for _, target := range a.knockTargetMap {
			return target
		}
	}
	return nil
}

func (a *UdpAgent) FindKnockTargetByResourceId(resId string) *KnockTarget {
	a.knockTargetMapMutex.Lock()
	a.knockTargetMapMutex.Unlock()
	for _, target := range a.knockTargetMap {
		if strings.EqualFold(resId, target.ResourceId) {
			return target
		}
	}
	return nil
}
