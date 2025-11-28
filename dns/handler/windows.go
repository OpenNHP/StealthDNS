package handler

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/OpenNHP/StealthDNS/common"
	"github.com/OpenNHP/opennhp/nhp/log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

type WindowsHandler struct {
	codePage           int
	isAdmin            bool
	dnsSetupFlag       bool
	removeLocalDNSFlag bool
	interfaces         []string
}

type CodePage struct {
	CodeSet     string
	OSLanguage  uint32
	CountryCode string
}

func NewWindowsHandler(removeLocalDNS bool) *WindowsHandler {
	h := &WindowsHandler{
		dnsSetupFlag:       false,
		removeLocalDNSFlag: removeLocalDNS,
		interfaces:         make([]string, 0),
	}
	h.isAdminPermission()
	h.getCodePage()
	return h
}

func (h *WindowsHandler) SetStealthDNS() (bool, error) {
	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return false, nil
	}
	interfaces, err := h.getNetworkInterfaces()
	if err != nil {
		log.Error("get network interfaces fail, %v", err)
		return false, err
	}
	index := 0
	for _, s := range interfaces {
		r, e := h.setLocalDNS(s)
		if e != nil {
			log.Error("set DNS [%s] fail: %v", err)
			continue
		}
		if r {
			h.interfaces = append(h.interfaces, s)
			h.dnsSetupFlag = true
			index++
		}
	}
	if index == 0 {
		return false, errors.New("no network interface dns setup success")
	}
	if index < len(interfaces) {
		log.Warning("Not all DNS settings were applied successfully. Please check your DNS configuration.")
	}
	return true, nil
}

func (h *WindowsHandler) RemoveStealthDNS() {
	if !h.removeLocalDNSFlag && !h.dnsSetupFlag {
		log.Debug("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}

	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}

	if len(h.interfaces) == 0 {
		log.Warning("No interfaces found, DNS reset is not required.")
		return
	}

	for _, s := range h.interfaces {
		r, _ := h.removeLocalDNS(s)
		if r {
			h.dnsSetupFlag = false
		}
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

func (h *WindowsHandler) getCodePage() {
	cp := getCodePageByWmi()
	if cp == 0 {
		cp = getCodePageByPowerShell()
	}
	h.codePage = cp
}

func getCodePageByWmi() int {
	cmd := exec.Command("wmic", "os", "get", "CodeSet", "/value")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(string(output), "\r\n")
	for _, line := range lines {
		if strings.Contains(line, "CodeSet") {
			codePage := line[8:]
			if cpNum, err := strconv.Atoi(codePage); err == nil {
				return cpNum
			}
		}
	}
	return 0
}

func getCodePageByPowerShell() int {
	cmd := exec.Command("powershell", "-Command", "[Console]::OutputEncoding.CodePage")
	output, err := cmd.Output()
	if err == nil {
		cpStr := strings.TrimSpace(string(output))
		if cp, err := strconv.Atoi(cpStr); err == nil && cp > 0 {
			return cp
		}
	}
	return 0
}

func (h *WindowsHandler) getNetworkInterfaces() ([]string, error) {
	cmd := exec.Command("powershell", "-Command", `
        $defaultRoute = Get-NetRoute -DestinationPrefix "0.0.0.0/0" | Sort-Object -Property InterfaceMetric | Select-Object -First 1
        if ($defaultRoute) {
            $adapter = Get-NetAdapter -InterfaceIndex $defaultRoute.InterfaceIndex
            $adapter.Name
        } else {
            # If default route is missing, pick the first connected network adapter.
            $connectedAdapters = Get-NetAdapter | Where-Object { $_.Status -eq 'Up' -and $_.InterfaceOperationalStatus -eq 'Connected' }
            if ($connectedAdapters) {
                $connectedAdapters | Sort-Object -Property InterfaceMetric | Select-Object -First 1 -ExpandProperty Name
            }
        }
    `)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	interfaces := make([]string, 0)
	str, err := h.BytesToString(output)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(str, "\r\n")
	for _, line := range lines {
		if len(line) > 0 {
			interfaces = append(interfaces, line)
		}
	}
	return interfaces, nil
}

func (h *WindowsHandler) setLocalDNS(interfaceName string) (bool, error) {
	cmd := exec.Command("netsh", "interface", "ip", "show", "dns", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	str, err := h.BytesToString(output)
	if err != nil {
		return false, err
	}
	lines := strings.Split(str, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, common.StealthDnsIp) {
			return true, nil
		}
	}
	//Set 127.0.0.1 as the primary DNS because when the primary DNS returns DNS_PROBE_FINISHED_NXDOMAIN, the server will not query the secondary DNS.
	//This would cause 127.0.0.1 to fail to take effect if it were set as the secondary DNS. Therefore, 127.0.0.1 must be set as the primary DNS to ensure
	//all requests are processed by StealthDNS. For domains it doesn't recognize, StealthDNS will reply with NOANSWER or forward the request to the configured
	//upstream DNS, ensuring that non-nhp related domains can be resolved normally.
	cmd = exec.Command("netsh", "interface", "ip", "add", "dns", interfaceName, common.StealthDnsIp, "index=1")
	//cmd = exec.Command("netsh", "interface", "ip", "set", "dns", interfaceName, "static", common.StealthDnsIp)
	err = cmd.Run()
	if err == nil {
		return true, nil
	}
	return false, err
}

func (h *WindowsHandler) removeLocalDNS(interfaceName string) (bool, error) {
	cmd := exec.Command("netsh", "interface", "ip", "show", "dns", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	str, err := h.BytesToString(output)
	if err != nil {
		return false, err
	}
	lines := strings.Split(str, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, common.StealthDnsIp) {
			cmd = exec.Command("netsh", "interface", "ip", "delete", "dns", interfaceName, common.StealthDnsIp)
			err = cmd.Run()
			if err == nil {
				return true, nil
			}
			return false, err
		}
	}
	return false, err
}

func (h *WindowsHandler) BytesToString(output []byte) (string, error) {
	switch h.codePage {
	case 936:
		// GBK
		return simplifiedchinese.GBK.NewDecoder().String(string(output))
	case 950:
		return traditionalchinese.Big5.NewDecoder().String(string(output))
	default:
		// UTF-8
		return string(output), nil
	}
}
