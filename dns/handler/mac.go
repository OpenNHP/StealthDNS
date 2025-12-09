//go:build darwin

package handler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/OpenNHP/StealthDNS/common"
	"github.com/OpenNHP/opennhp/nhp/log"
)

type MacHandler struct {
	isRoot        bool
	dnsSetupFlag  bool
	interfaceName string
	backupDNS     []string
	upstreamDNS   string
}

func NewDNSHandler() *MacHandler {
	h := &MacHandler{
		dnsSetupFlag: false,
	}
	h.isRootPermission()
	return h
}

func (h *MacHandler) SetStealthDNS() (bool, error) {
	if !h.isRoot {
		log.Warning("The current account does not have root privileges. Please manually configure the alternate DNS.")
		return false, nil
	}
	err := h.getInterfaceName()
	if err != nil {
		log.Error("Failed to retrieve network settings: %v", err)
		return false, err
	}
	err = h.backupAndSet()
	if err != nil {
		return false, err
	}
	h.dnsSetupFlag = true
	return true, nil
}

func (h *MacHandler) RemoveStealthDNS() {
	if !h.dnsSetupFlag {
		log.Warning("The DNS proxy address is configured manually, and the StealthDNS service does not automatically restore the settings.")
		return
	}

	if !h.isRoot {
		log.Warning("The current account does not have root privileges. Please manually configure the alternate DNS.")
		return
	}

	if len(h.interfaceName) == 0 {
		log.Warning("No interfaces found, DNS reset is not required.")
		return
	}
	err := h.setDNS(h.backupDNS)
	if err != nil {
		log.Error("dns network [%s] reset fail: %v", h.interfaceName, err)
	}
}

func (h *MacHandler) getInterfaceName() error {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		service := strings.TrimSpace(scanner.Text())
		if service == "" || strings.HasPrefix(service, "*") {
			continue
		}
		if h.isNetworkServiceActive(service) {
			h.interfaceName = service
		}
	}
	if len(h.interfaceName) == 0 {
		return errors.New("no network service found")
	}
	return nil
}

func (h *MacHandler) isNetworkServiceActive(service string) bool {
	cmd := exec.Command("/usr/sbin/networksetup", "-getinfo", service)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "IP address:") && !strings.Contains(string(output), "IP address: 0.0.0.0")
}

func (h *MacHandler) backupAndSet() error {
	err := h.storeBackupDNS()
	if err != nil {
		return err
	}
	dnsServices := make([]string, 0)
	dnsServices = append(dnsServices, common.StealthDnsIp)
	dnsServices = append(dnsServices, h.backupDNS...)
	err = h.setDNS(dnsServices)
	return err
}

func (h *MacHandler) setDNS(dnsServers []string) error {
	args := []string{"-setdnsservers", h.interfaceName}
	args = append(args, dnsServers...)

	cmd := exec.Command("networksetup", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("set DNS fail: %v, output: %s", err, string(output))
	}
	return nil
}

func (h *MacHandler) storeBackupDNS() error {
	cmd := exec.Command("networksetup", "-getdnsservers", h.interfaceName)
	output, err := cmd.Output()
	if err != nil {
		// If no DNS is set, the command will return an error, but this is normal.
		if strings.Contains(string(output), "aren't any DNS Servers set") {
			return errors.New("no DNS servers set")
		}
		return err
	}

	var dnsServers []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		dns := strings.TrimSpace(scanner.Text())
		if dns != "" && !strings.EqualFold(dns, common.StealthDnsIp) {
			dnsServers = append(dnsServers, dns)
		}
	}
	h.backupDNS = dnsServers
	if len(dnsServers) == 0 {
		log.Warning("Failed to obtain upstream DNS; using default DNS [%s] as the upstream DNS", common.DefaultUpstreamDNS)
		h.upstreamDNS = common.DefaultUpstreamDNS
	} else {
		h.upstreamDNS = dnsServers[0]
	}
	return err
}

func (h *MacHandler) isRootPermission() {
	h.isRoot = os.Geteuid() == 0
}

func (h *MacHandler) GetUpstreamDNS() string {
	return h.upstreamDNS
}
