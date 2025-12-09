package dns

import (
	"github.com/OpenNHP/StealthDNS/common"
	"github.com/OpenNHP/StealthDNS/dns/handler"
	"github.com/OpenNHP/opennhp/nhp/log"
)

type Manager struct {
	handler handler.Handler
}

func NewDNSManager() *Manager {
	m := &Manager{}
	m.handler = handler.NewDNSHandler()
	return m
}

func (d *Manager) SetStealthDNS() bool {
	log.Debug("start setup Stealth DNS...")
	r, err := d.handler.SetStealthDNS()
	if err != nil {
		log.Error("setup Stealth DNS fail: %v", err)
		log.Warning("manual configuration of stealth DNS is required")
		return false
	}
	log.Debug("setup Stealth DNS success!")
	return r
}

func (d *Manager) RemoveStealthDNS() {
	d.handler.RemoveStealthDNS()
}

func (d *Manager) GetUpstreamDNS() string {
	upstreamDNS := d.handler.GetUpstreamDNS()
	if len(upstreamDNS) == 0 {
		upstreamDNS = common.DefaultUpstreamDNS
	}
	log.Debug("upstream dns is %s", upstreamDNS)
	return upstreamDNS
}
