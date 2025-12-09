//go:build windows

package handler

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenNHP/StealthDNS/common"
	"github.com/OpenNHP/opennhp/nhp/log"
	"github.com/StackExchange/wmi"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

type WindowsHandler struct {
	isAdmin       bool
	dnsSetupFlag  bool
	interfaceName string
	backupDNS     []string
	upstreamDNS   string
	codePage      uint32
}

type NetworkAdapter struct {
	Index           uint32 `wmi:"Index"`
	GUID            string `wmi:"GUID"`
	NetConnectionID string `wmi:"NetConnectionID"`
}

type NetworkAdapterConfig struct {
	InterfaceIndex       uint32   `wmi:"InterfaceIndex"`
	SettingID            string   `wmi:"SettingID"`
	DNSServerSearchOrder []string `wmi:"DNSServerSearchOrder"`
	IPEnabled            bool     `wmi:"IPEnabled"`
}

type WinInterInfo struct {
	Name   string
	Metric int
}

func NewDNSHandler() *WindowsHandler {
	h := &WindowsHandler{
		dnsSetupFlag: false,
	}
	h.isAdminPermission()
	h.getOEMCodePage()
	return h
}

func (h *WindowsHandler) GetUpstreamDNS() string {
	return h.upstreamDNS
}

func (h *WindowsHandler) SetStealthDNS() (bool, error) {
	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return true, nil
	}
	info, err := h.getActiveInterfacesInfo()
	if err != nil {
		log.Error("get active interface info fail: %v", err)
		return false, err
	}
	err = h.getPrimaryInterface(info)
	if err != nil {
		log.Error("get primary interface fail: %v", err)
		return false, err
	}
	err = h.setDNS(false)
	if err != nil {
		log.Error("set stealth dns proxy fail: %v", err)
		return false, err
	}
	h.dnsSetupFlag = true
	return true, nil
}
func (h *WindowsHandler) RemoveStealthDNS() {
	if !h.dnsSetupFlag {
		log.Warning("The DNS proxy address is configured manually, and the StealthDNS service does not automatically restore the settings.")
		return
	}
	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}
	err := h.setDNS(true)
	if err != nil {
		log.Error("restore the DNS configuration fail: %v", err)
		log.Debug("please manually restore the DNS configuration.")
	}
}

func (h *WindowsHandler) isAdminPermission() {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err == nil {
		h.isAdmin = true
	} else {
		h.isAdmin = false
	}
}

func (h *WindowsHandler) getOEMCodePage() {
	cmd := exec.Command("powershell", "-Command", "[Console]::OutputEncoding.CodePage")
	output, err := cmd.Output()
	if err != nil {
		h.codePage = 437
	}
	cp, err := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32)
	if err != nil {
		h.codePage = 437
	}
	h.codePage = uint32(cp)
}

func (h *WindowsHandler) decodeFromCodePage(data []byte) (string, error) {
	if h.codePage == 65001 {
		return string(data), nil
	}

	mimeName := fmt.Sprintf("cp%d", h.codePage)
	enc, err := ianaindex.MIME.Encoding(mimeName)
	if err != nil {
		return string(data), nil
	}

	decoder := enc.NewDecoder()
	utf8Str, _, err := transform.String(decoder, string(data))
	if err != nil {
		return string(data), nil
	}
	return utf8Str, nil
}

func (h *WindowsHandler) getActiveInterfacesWithMetric() ([]WinInterInfo, error) {
	cmd := exec.Command("netsh", "interface", "ipv4", "show", "interfaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	text, err := h.decodeFromCodePage(output)
	if err != nil {
		text = string(output)
	}

	lines := strings.Split(text, "\n")
	var interfaces []WinInterInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "connected") {
			continue
		}

		fields := regexp.MustCompile(`\s{2,}`).Split(line, -1)
		if len(fields) < 4 {
			continue
		}

		if _, err := strconv.Atoi(fields[0]); err != nil {
			continue
		}

		metricStr := fields[1]
		name := strings.TrimSpace(fields[len(fields)-1])

		metric, err := strconv.Atoi(metricStr)
		if err != nil {
			continue
		}

		interfaces = append(interfaces, WinInterInfo{Name: name, Metric: metric})
	}

	return interfaces, nil
}

// Get the interface with the lowest metric (highest priority)
func (h *WindowsHandler) getPrimaryInterface(interMap map[string][]string) error {
	interfaces, err := h.getActiveInterfacesWithMetric()
	if err != nil {
		return err
	}
	if len(interfaces) == 0 {
		return fmt.Errorf("no active network interface found")
	}

	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Metric < interfaces[j].Metric
	})
	for _, interfaceInfo := range interfaces {
		if dnsList, ok := interMap[interfaceInfo.Name]; ok {
			h.interfaceName = interfaceInfo.Name
			h.backupDNS = make([]string, 0)
			for _, dnsIp := range dnsList {
				if dnsIp == common.StealthDnsIp {
					continue
				}
				h.backupDNS = append(h.backupDNS, dnsIp)
			}
			if len(h.backupDNS) == 0 {
				log.Warning("Failed to obtain upstream DNS; using default DNS [%s] as the upstream DNS", common.DefaultUpstreamDNS)
				h.upstreamDNS = common.DefaultUpstreamDNS
			} else {
				h.upstreamDNS = h.backupDNS[0]
			}
			return nil
		}
	}
	return fmt.Errorf("no active network interface found")
}

func (h *WindowsHandler) setDNS(restoreFlag bool) error {
	dnsList := make([]string, 0)
	if !restoreFlag {
		dnsList = append(dnsList, common.StealthDnsIp)
	}
	for _, dnsIp := range h.backupDNS {
		dnsList = append(dnsList, dnsIp)
	}
	cmd := exec.Command("netsh", "interface", "ip", "delete", "dns", h.interfaceName, "all")
	err := cmd.Run()
	if err != nil {
		log.Error("Failed to clear system DNS server configuration: %v", err)
		return err
	}
	cmd = exec.Command("netsh", "interface", "ipv4", "set", "dns",
		h.interfaceName, "static", dnsList[0], "primary")

	err = cmd.Run()

	if err != nil {
		log.Warning("Failed to set primary DNS：%v", err)
		return err
	}
	for i := 1; i < len(dnsList); i++ {
		cmd = exec.Command("netsh", "interface", "ipv4", "add", "dns",
			h.interfaceName, dnsList[i], "index="+strconv.Itoa(i+1))
		_ = cmd.Run() // Ignore failure to add secondary DNS.
	}

	return nil
}

func (h *WindowsHandler) getActiveInterfacesInfo() (map[string][]string, error) {
	var adapters []NetworkAdapter
	err := wmi.Query("SELECT * FROM Win32_NetworkAdapter", &adapters)
	if err != nil {
		return nil, err
	}

	adapterMap := make(map[string]string)
	for _, a := range adapters {
		if a.NetConnectionID != "" {
			adapterMap[a.GUID] = a.NetConnectionID
		}
	}

	var configs []NetworkAdapterConfig
	err = wmi.Query("SELECT * FROM Win32_NetworkAdapterConfiguration WHERE IPEnabled=True", &configs)
	if err != nil {
		return nil, err
	}
	interMap := make(map[string][]string)
	for _, cfg := range configs {
		if len(cfg.DNSServerSearchOrder) == 0 {
			continue
		}
		if name, ok := adapterMap[cfg.SettingID]; ok {
			interMap[name] = cfg.DNSServerSearchOrder
		}
	}
	return interMap, nil
}
