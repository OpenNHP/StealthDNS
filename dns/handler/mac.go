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
	isAdmin            bool
	dnsSetupFlag       bool
	removeLocalDNSFlag bool
	dnsServersCache    map[string][]string
	networkServices    []string
}

func NewMacHandler(removeLocalDNS bool) *MacHandler {
	h := &MacHandler{
		dnsSetupFlag:       false,
		removeLocalDNSFlag: removeLocalDNS,
		dnsServersCache:    make(map[string][]string),
		networkServices:    make([]string, 0),
	}
	h.isAdminPermission()
	return h
}

func (h *MacHandler) SetStealthDNS() (bool, error) {
	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return false, nil
	}
	err := h.getNetworkServices()
	if err != nil {
		log.Error("Failed to retrieve network settings: %v", err)
		return false, err
	}
	index := 0
	for _, service := range h.networkServices {
		err = h.backupAndSet(service)
		if err != nil {
			log.Error("set DNS [%s] fail: %v", err)
			continue
		}
		index++
	}
	if index == 0 {
		return false, errors.New("no network interface dns setup success")
	}
	if index < len(h.networkServices) {
		log.Warning("Not all DNS settings were applied successfully. Please check your DNS configuration.")
	}
	return true, nil
}

func (h *MacHandler) RemoveStealthDNS() {
	if !h.removeLocalDNSFlag && !h.dnsSetupFlag {
		log.Debug("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}

	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}

	if len(h.networkServices) == 0 {
		log.Warning("No interfaces found, DNS reset is not required.")
		return
	}
	for _, service := range h.networkServices {
		err := h.setDNS(service, h.dnsServersCache[service])
		if err != nil {
			log.Error("dns network [%s] reset fail: %v", err)
		}
	}
}

func (h *MacHandler) getNetworkServices() error {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var services []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Skip the first header row.
	if scanner.Scan() {
		for scanner.Scan() {
			service := strings.TrimSpace(scanner.Text())
			if service != "" && !strings.HasPrefix(service, "*") {
				services = append(services, service)
			}
		}
	}
	if len(services) == 0 {
		return errors.New("no network service found")
	}
	h.networkServices = services
	return nil
}

func (h *MacHandler) backupAndSet(service string) error {
	err := h.backupDNS(service)
	if err != nil {
		return err
	}
	dnsServices := make([]string, 0)
	dnsServices = append(dnsServices, common.StealthDnsIp)
	dnsServices = append(dnsServices, h.dnsServersCache[service]...)
	err = h.setDNS(service, dnsServices)
	return err
}

func (h *MacHandler) setDNS(service string, dnsServers []string) error {
	args := []string{"-setdnsservers", service}
	args = append(args, dnsServers...)

	cmd := exec.Command("networksetup", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("set DNS fail: %v, output: %s", err, string(output))
	}
	return nil
}

func (h *MacHandler) backupDNS(service string) error {
	cmd := exec.Command("networksetup", "-getdnsservers", service)
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
		if dns != "" {
			dnsServers = append(dnsServers, dns)
		}
	}
	h.dnsServersCache[service] = dnsServers
	return err
}

func (h *MacHandler) isAdminPermission() {
	h.isAdmin = os.Geteuid() == 0
}
