package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenNHP/StealthDNS/agent"
	"github.com/OpenNHP/StealthDNS/common"
	com "github.com/OpenNHP/opennhp/nhp/common"
	"github.com/OpenNHP/opennhp/nhp/log"
	"github.com/miekg/dns"
)

// ProxyService dns proxy service
type ProxyService struct {
	nhpAgent agent.NhpAgent

	dnsManager *Manager
	config     *Config
	log        *log.Logger
	listenAddr *net.UDPAddr
	listenConn *net.UDPConn
	localIp    string
	localMac   string

	server *dns.Server

	upstreamDNS string
	dnsCache    *StealthDNSCache

	domainMap     map[string]string
	domainMapLock sync.Mutex

	resourceMap     map[string]*Resource
	resourceMapLock sync.Mutex

	running atomic.Bool
}

func (p *ProxyService) Start(dirPath string, logLevel int) (err error) {
	common.ExeDirPath = dirPath
	p.log = log.NewLogger("StealthDNS", logLevel, filepath.Join(dirPath, "logs"), "proxy-service")
	log.SetGlobalLogger(p.log)
	log.Info("=========================================================")
	log.Info("=============== Stealth DNS started =====================")
	log.Info("=========================================================")

	err = p.loadDNSConfig()
	if err != nil {
		return err
	}
	p.dnsCache = NewStealthDNSCache(time.Duration(10) * time.Second)

	err = p.loadResources()
	if err != nil {
		return err
	}

	p.dnsManager = NewDNSManager()
	if !p.dnsManager.SetStealthDNS() {
		log.Warning("Stealth DNS setup failed. Please ensure the DNS proxy address 127.0.0.1 is set as the alternate DNS.")
	}
	p.upstreamDNS = p.dnsManager.GetUpstreamDNS()
	p.nhpAgent, err = agent.NewNhpAgent(dirPath)
	if err != nil {
		log.Error("init nhp-agent fail: %v", err)
		return err
	}
	err = p.nhpAgent.AgentInit(dirPath, logLevel)
	if err != nil {
		return err
	}
	listenAddr := fmt.Sprintf("%s:%d", common.StealthDnsIp, common.DnsUdpPort)

	p.server = &dns.Server{
		Addr:    listenAddr,
		Net:     "udp",
		Handler: p,
	}
	go p.startServer()
	p.running.Store(true)
	return nil
}

func (p *ProxyService) startServer() {
	err := p.server.ListenAndServe()
	if err != nil {
		log.Error("dns server listen fail: %v", err)
		p.Stop()
	}
}

func (p *ProxyService) Stop() {
	if p.running.Load() {
		if !p.running.CompareAndSwap(true, false) {
			log.Debug("Stealth DNS service is stopping.")
			return
		}
	} else {
		log.Debug("Stealth DNS service has stopped.")
		return
	}

	log.Debug("stop Stealth DNS")
	p.dnsManager.RemoveStealthDNS()
	if p.nhpAgent != nil {
		_ = p.nhpAgent.AgentClose()
	}
	if p.server != nil {
		_ = p.server.Shutdown()
	}
	log.Info("===========================")
	log.Info("=== Stealth DNS stopped ===")
	log.Info("===========================")
	log.Close()
	os.Exit(0)
}

func (p *ProxyService) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	domainName := r.Question[0].Name
	log.Debug("domain name：%s, question type：%s", domainName, dns.TypeToString[r.Question[0].Qtype])
	if strings.Contains(domainName, common.NhpDomainNameSuffix) {
		resId := domainName[:strings.Index(domainName, common.NhpDomainNameSuffix)]
		if r.Question[0].Qtype == common.Type_A || r.Question[0].Qtype == common.Type_AAAA {
			// only process domain name resolution requests of type A or AAAA
			go p.nhpServer(w, r, resId)
		} else {
			go p.noAnswer(w, r)
		}
	} else {
		go p.forwardUpstreamDNS(w, r)
	}
}

func (p *ProxyService) forwardUpstreamDNS(w dns.ResponseWriter, r *dns.Msg) {
	// forward to upstream DNS

	client := &dns.Client{
		Timeout: 5 * time.Second,
	}
	resp, _, err := client.Exchange(r, p.upstreamDNS+":53")
	if err != nil {
		log.Warning("domain :%s request upstream DNS fail: %v", r.Question[0].Name, err)
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}
	resp.Id = r.Id
	_ = w.WriteMsg(resp)
	return
}

func (p *ProxyService) nhpServer(w dns.ResponseWriter, r *dns.Msg, resId string) {
	if item, found := p.dnsCache.GetCache(resId); found {
		item.value.Id = r.Id
		_ = w.WriteMsg(item.value)
		return
	}

	resultCh := p.dnsCache.group.DoChan(resId, func() (interface{}, error) {
		if item, found := p.dnsCache.GetCache(resId); found {
			item.value.Id = r.Id
			_ = w.WriteMsg(item.value)
			return item.value, nil
		}

		ackMsg, err := p.knock(resId)
		if err != nil {
			p.noAnswer(w, r)
			return nil, err
		}
		if ackMsg == nil {
			p.noAnswer(w, r)
			return nil, errors.New("request nhp-server fail")
		}
		if !strings.EqualFold(ackMsg.ErrCode, "0") {
			p.noAnswer(w, r)
			return nil, com.ErrorCodeToError(ackMsg.ErrCode)
		}
		openTime := ackMsg.OpenTime
		if openTime > 5 {
			// The cache time is reduced by 5 seconds to prevent the port from being closed by nhp-ac right after the domain name resolution result is returned.
			openTime -= 5
		}

		var m *dns.Msg
		for _, host := range ackMsg.ResourceHost {
			if net.ParseIP(host) != nil {
				//
				m, err = p.handleQuery(w, r, host, openTime)
			} else {
				log.Debug("nhp server returns a domain name as the result, which requires further DNS resolution.")
				m, err = p.handleUpstreamQuery(w, r, host, openTime)
			}
			if err != nil {
				return nil, err
			} else {
				break
			}
		}

		if m != nil {
			p.dnsCache.SetCacheWithTTL(resId, m, time.Duration(openTime)*time.Second)
		}

		return nil, err
	})

	// waiting result
	select {
	case result := <-resultCh:
		if result.Err != nil {
			log.Error("query dns Answer fail,err is %v", result.Err)
		} else {
			log.Debug("nhp knock success")
		}
	case <-time.After(5 * time.Second): // time out
		log.Error("timeout waiting for creation")
	}
}

func (p *ProxyService) queryUpstream(domain string, qtype uint16) (*dns.Msg, error) {
	client := &dns.Client{
		Timeout: 5 * time.Second,
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), qtype)

	response, _, err := client.Exchange(msg, p.upstreamDNS+":53")
	if err != nil {
		return nil, fmt.Errorf("upstream query failed: %v", err)
	}

	return response, nil
}

func (p *ProxyService) knock(resId string) (ackMsg *com.ServerKnockAckMsg, err error) {
	if target, ok := p.resourceMap[resId]; !ok {
		log.Warning("unknow resource [%s],", resId)
		return nil, nil
	} else {
		resource, err := p.nhpAgent.AgentKnockResource(target.AuthServiceId, target.ResourceId, target.ServerIp, target.ServerHostname, target.ServerPort)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(resource), &ackMsg)
		if err != nil {
			return nil, err
		}
		//
		return ackMsg, nil
	}
}

func (p *ProxyService) noAnswer(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetRcode(r, dns.RcodeSuccess)
	_ = w.WriteMsg(m)
}

func (p *ProxyService) handleQuery(w dns.ResponseWriter, r *dns.Msg, ip string, ttl uint32) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetReply(r)
	var rr dns.RR
	var err error
	if r.Question[0].Qtype == 1 {
		rr, err = dns.NewRR(fmt.Sprintf("%s A %s", r.Question[0].Name, ip))
	} else {
		rr, err = dns.NewRR(fmt.Sprintf("%s 3600 IN AAAA %s", r.Question[0].Name, ip))
	}

	if err == nil {
		m.Answer = append(m.Answer, rr)
		m.Answer[0].Header().Ttl = ttl
		_ = w.WriteMsg(m)
		return m, nil
	} else {
		log.Error("create dns answer fail, %v", err)
		p.noAnswer(w, r)
		return nil, err
	}
}

func (p *ProxyService) handleUpstreamQuery(w dns.ResponseWriter, r *dns.Msg, host string, ttl uint32) (*dns.Msg, error) {
	question := r.Question[0]
	m := new(dns.Msg)
	m.SetReply(r)
	cname := &dns.CNAME{
		Hdr:    dns.RR_Header{Name: question.Name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: ttl},
		Target: host + ".",
	}
	m = new(dns.Msg)
	m.SetReply(r)
	m.Answer = append(m.Answer, cname)
	client := &dns.Client{
		Timeout: 5 * time.Second,
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(host), r.Question[0].Qtype)

	response, _, err := client.Exchange(msg, p.upstreamDNS+":53")
	if err != nil {
		log.Error("create dns answer fail, %v", err)
		p.noAnswer(w, r)
		return nil, err
	}

	if len(response.Answer) == 0 {
		log.Error("create dns answer fail, %v", err)
		p.noAnswer(w, r)
		return nil, err
	}

	for _, rr := range response.Answer {
		m.Answer = append(m.Answer, rr)
	}
	_ = w.WriteMsg(m)
	return m, nil
}
