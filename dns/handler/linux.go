//go:build linux

package handler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenNHP/StealthDNS/common"
	"github.com/OpenNHP/opennhp/nhp/log"
)

type DNSManagement struct {
	Method        string
	ResolvConf    string
	IsSymlink     bool
	SymlinkTarget string
	IsActive      bool
}

type LinuxInterInfo struct {
	Name    string
	Metric  int
	IsUp    bool
	HasIPv4 bool
}

type LinuxHandler struct {
	isRoot        bool
	dnsSetupFlag  bool
	interfaceInfo LinuxInterInfo
	backupDNS     []string
	upstreamDNS   string
	dnsManagement *DNSManagement
}

func NewDNSHandler() *LinuxHandler {
	handler := &LinuxHandler{
		dnsSetupFlag: false,
		backupDNS:    make([]string, 0),
	}
	handler.isRootPermission()
	_ = handler.detectDNSManagement()
	return handler
}

func (h *LinuxHandler) isRootPermission() {
	h.isRoot = os.Geteuid() == 0
}

func (h *LinuxHandler) SetStealthDNS() (bool, error) {
	if !h.isRoot {
		log.Warning("The current account does not have root privileges. Please manually configure the alternate DNS.")
		return false, nil
	}
	interfaces, err := getActiveInterfaces()
	if err != nil {
		return false, err
	}
	if len(interfaces) == 0 {
		return false, fmt.Errorf("no active network interfaces")
	}
	h.interfaceInfo = interfaces[0]
	err = h.getCurrentDNSServers()
	if err != nil {
		return false, err
	}
	dnsList := append([]string{common.StealthDnsIp}, h.backupDNS...)

	err = h.setDNS(dnsList)
	if err != nil {
		return false, err
	}
	h.dnsSetupFlag = true
	return true, nil
}

func (h *LinuxHandler) RemoveStealthDNS() {
	if !h.dnsSetupFlag {
		log.Warning("The DNS proxy address is configured manually, and the StealthDNS service does not automatically restore the settings.")
		return
	}

	if !h.isRoot {
		log.Warning("The current account does not have root privileges. Please manually configure the alternate DNS.")
		return
	}

	err := h.restoreDNS()
	if err != nil {
		log.Error("dns network [%s] reset fail: %v", h.interfaceInfo.Name, err)
	}
}

func (h *LinuxHandler) GetUpstreamDNS() string {
	return h.upstreamDNS
}

func (h *LinuxHandler) detectDNSManagement() error {
	h.dnsManagement = &DNSManagement{Method: "unknown"}
	resolv := "/etc/resolv.conf"
	fi, err := os.Lstat(resolv)
	if err != nil {
		return err
	}
	h.dnsManagement.ResolvConf = resolv
	h.dnsManagement.IsSymlink = fi.Mode()&os.ModeSymlink != 0
	if h.dnsManagement.IsSymlink {
		target, _ := os.Readlink(resolv)
		h.dnsManagement.SymlinkTarget = target
	}

	// systemd-resolved
	if isServiceActive("systemd-resolved") && strings.Contains(h.dnsManagement.SymlinkTarget, "systemd/resolve") {
		h.dnsManagement.Method = "systemd-resolved"
		h.dnsManagement.IsActive = true
		return nil
	}

	// NetworkManager
	if isServiceActive("NetworkManager") && strings.Contains(h.dnsManagement.SymlinkTarget, "NetworkManager") {
		h.dnsManagement.Method = "networkmanager"
		h.dnsManagement.IsActive = true
		return nil
	}

	// dhcpcd
	if isDhcpcdRunning() {
		h.dnsManagement.Method = "dhcpcd"
		h.dnsManagement.IsActive = true
		return nil
	}

	h.dnsManagement.Method = "static"
	h.dnsManagement.IsActive = true
	return nil
}

func isServiceActive(service string) bool {
	cmd := exec.Command("systemctl", "is-active", service)
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "active"
}

func isDhcpcdRunning() bool {
	return exec.Command("pidof", "dhcpcd").Run() == nil
}

func getActiveInterfaces() ([]LinuxInterInfo, error) {
	cmd := exec.Command("ip", "route", "show", "default")
	output, err := cmd.Output()
	if err != nil {
		return getAllInterfacesWithMetric()
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "default") {
			re := regexp.MustCompile(`dev\s+(\S+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				dev := matches[1]
				metric := 0
				if m := regexp.MustCompile(`metric\s+(\d+)`).FindStringSubmatch(line); len(m) > 1 {
					if v, _ := strconv.Atoi(m[1]); v > 0 {
						metric = v
					}
				}
				return []LinuxInterInfo{{Name: dev, Metric: metric, IsUp: true, HasIPv4: true}}, nil
			}
		}
	}

	return getAllInterfacesWithMetric()
}

func getAllInterfacesWithMetric() ([]LinuxInterInfo, error) {
	var interfaces []LinuxInterInfo

	cmd := exec.Command("ip", "link", "show", "up")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	linkRe := regexp.MustCompile(`^\d+:\s+(\S+?):`)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := linkRe.FindStringSubmatch(line); len(matches) > 1 {
			name := strings.TrimSuffix(matches[1], ":")
			// skip lo
			if name == "lo" {
				continue
			}
			hasIPv4 := hasIPv4Address(name)
			if !hasIPv4 {
				continue
			}
			// Sort by Metric in ascending order (lower values have higher priority).
			metric := getInterfaceMetric(name)
			interfaces = append(interfaces, LinuxInterInfo{
				Name: name, Metric: metric, IsUp: true, HasIPv4: true,
			})
		}
	}

	// Sort by Metric in ascending order (lower values have higher priority).
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Metric < interfaces[j].Metric
	})

	return interfaces, nil
}

func hasIPv4Address(iface string) bool {
	cmd := exec.Command("ip", "addr", "show", iface)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "inet ")
}

func getInterfaceMetric(iface string) int {
	cmd := exec.Command("ip", "route", "show", "dev", iface)
	output, err := cmd.Output()
	if err != nil {
		return 100 // default high metric
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "default") || strings.Contains(line, "via") {
			if m := regexp.MustCompile(`metric\s+(\d+)`).FindStringSubmatch(line); len(m) > 1 {
				if v, _ := strconv.Atoi(m[1]); v > 0 {
					return v
				}
			}
			return 100
		}
	}
	return 100
}

func (h *LinuxHandler) getCurrentDNSServers() error {
	var dnsList []string
	switch h.dnsManagement.Method {
	case "systemd-resolved":
		cmd := exec.Command("resolvectl", "dns", h.interfaceInfo.Name)
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		dnsList = parseResolvectlDNS(string(output))
	case "networkmanager":
		cmd := exec.Command("nmcli", "-g", "ipv4.dns", "con", "show", getNMConnectionForInterface(h.interfaceInfo.Name))
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		dnsStr := strings.TrimSpace(string(output))
		if dnsStr == "" {
			return err
		}
		dnsList = strings.Split(dnsStr, ",")
	case "dhcpcd", "static":
		dnsList = readDNSServersFromResolvConf()
	default:
		dnsList = readDNSServersFromResolvConf()
	}
	for _, dnsIp := range dnsList {
		if len(dnsIp) > 0 && !strings.EqualFold(common.StealthDnsIp, dnsIp) {
			h.backupDNS = append(h.backupDNS, dnsIp)
		}
	}
	if len(h.backupDNS) == 0 {
		log.Warning("Failed to obtain upstream DNS; using default DNS [%s] as the upstream DNS", common.DefaultUpstreamDNS)
		h.upstreamDNS = common.DefaultUpstreamDNS
	} else {
		h.upstreamDNS = h.backupDNS[0]
	}
	return nil
}

func parseResolvectlDNS(output string) []string {
	var servers []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "DNS Servers:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				serversStr := strings.TrimSpace(parts[1])
				if serversStr != "" {
					servers = strings.Split(serversStr, " ")
				}
			}
		}
	}
	return servers
}

func getNMConnectionForInterface(iface string) string {
	cmd := exec.Command("nmcli", "-t", "-f", "NAME,DEVICE", "con", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && parts[1] == iface {
			return parts[0]
		}
	}
	return ""
}

func readDNSServersFromResolvConf() []string {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				servers = append(servers, parts[1])
			}
		}
	}
	return servers
}

func (h *LinuxHandler) setDNS(newDNS []string) error {
	switch h.dnsManagement.Method {
	case "systemd-resolved":
		args := []string{"dns", h.interfaceInfo.Name}
		args = append(args, newDNS...)
		cmd := exec.Command("resolvectl", args...)
		return cmd.Run()
	case "networkmanager":
		connName := getNMConnectionForInterface(h.interfaceInfo.Name)
		if connName == "" {
			return fmt.Errorf("no NetworkManager connection found for interface %s", h.interfaceInfo.Name)
		}
		// Set the primary DNS first.
		cmd := exec.Command("nmcli", "con", "mod", connName, "ipv4.dns", strings.Join(newDNS, " "))
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("nmcli", "con", "up", connName)
		return cmd.Run()
	case "dhcpcd":
		return updateDhcpcdConf(newDNS)
	case "static":
		return writeResolvConf(newDNS)
	default:
		return fmt.Errorf("unsupported DNS management method: %s", h.dnsManagement.Method)
	}
}

func updateDhcpcdConf(dnsList []string) error {
	confPath := "/etc/dhcpcd.conf"
	content, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0)

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "static domain_name_servers=") ||
			strings.HasPrefix(strings.TrimSpace(line), "#static domain_name_servers=") {
			continue
		}
		newLines = append(newLines, line)
	}

	if len(dnsList) > 0 {
		newLines = append(newLines, fmt.Sprintf("static domain_name_servers=%s", strings.Join(dnsList, " ")))
	}

	return os.WriteFile(confPath, []byte(strings.Join(newLines, "\n")+"\n"), 0644)
}

func writeResolvConf(dnsList []string) error {
	content := "# Managed by go-dns-tool\n"
	for _, dns := range dnsList {
		content += fmt.Sprintf("nameserver %s\n", dns)
	}
	return os.WriteFile("/etc/resolv.conf", []byte(content), 0644)
}

func (h *LinuxHandler) restoreDNS() error {
	if len(h.backupDNS) == 0 {
		switch h.dnsManagement.Method {
		case "systemd-resolved":
			return exec.Command("resolvectl", "revert", h.interfaceInfo.Name).Run()
		case "networkmanager":
			connName := getNMConnectionForInterface(h.interfaceInfo.Name)
			if connName == "" {
				return nil
			}
			err := exec.Command("nmcli", "con", "mod", connName, "ipv4.dns", "").Run()
			if err != nil {
				return err
			}
			err = exec.Command("nmcli", "con", "mod", connName, "ipv4.ignore-auto-dns", "no").Run()
			if err != nil {
				return err
			}
			return exec.Command("nmcli", "con", "up", connName).Run()
		case "dhcpcd":
			err := updateDhcpcdConf(nil) // 移除 static 行
			if err != nil {
				return err
			}
			return exec.Command("systemctl", "restart", "dhcpcd").Run()
		case "static":
			return os.WriteFile("/etc/resolv.conf", []byte("# Restored by go-dns-tool\n"), 0644)
		}
	} else {
		return h.setDNS(h.backupDNS)
	}
	return nil
}
