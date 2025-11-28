package dns

import (
	"runtime"

	"github.com/OpenNHP/StealthDNS/dns/handler"
	"github.com/OpenNHP/opennhp/nhp/log"
)

type Manager struct {
	handler handler.Handler
}

func NewDNSManager(removeLocalDNS bool) *Manager {
	m := &Manager{}
	switch runtime.GOOS {
	case "windows":
		m.handler = handler.NewWindowsHandler(removeLocalDNS)
	case "darwin":
		m.handler = handler.NewMacHandler(removeLocalDNS)
	case "linux":
		m.handler = handler.NewLinuxHandler()
	default:
		log.Warning("%s operating system is not supported; unable to create a DNS handler.", runtime.GOOS)
	}
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
