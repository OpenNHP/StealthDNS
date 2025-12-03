package handler

import (
	"bufio"
	"github.com/OpenNHP/opennhp/nhp/log"
	"os"
	"strings"
)

const (
	dnsConfigFile  = "/etc/resolv.conf"
	dnsBackupFile  = "/etc/resolv_stealthdns.conf"
	stealthDNSPath = "nameserver 127.0.0.1"
)

type LinuxHandler struct {
	isAdmin            bool
	dnsSetupFlag       bool
	removeLocalDNSFlag bool
}

func (h *LinuxHandler) SetStealthDNS() (bool, error) {
	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return false, nil
	}
	err := h.backupDNSFile()
	if err != nil {
		log.Error("backup dns configuration file failed: %v", dnsConfigFile, err)
		return false, err
	}
	err = h.createStealthDNSFile()
	if err != nil {
		log.Error("set stealth dns config file fail: %v", err)
		return false, err
	}
	h.dnsSetupFlag = true
	return true, nil
}

func (h *LinuxHandler) RemoveStealthDNS() {
	if !h.removeLocalDNSFlag && !h.dnsSetupFlag {
		log.Debug("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}

	if !h.isAdmin {
		log.Warning("The current account does not have administrator privileges. Please manually configure the alternate DNS.")
		return
	}

	//	delete dnsConfigFile
	err := os.Remove(dnsConfigFile)
	if err != nil {
		log.Error("remove stealth dns configuration file /etc/resolv.conf fail: %v", err)
		return
	}
	err = os.Rename(dnsBackupFile, dnsConfigFile)
	if err != nil {
		log.Error("reset system dns configuration file /etc/resolv.conf fail: %v", err)
		return
	}
	err = os.Remove(dnsBackupFile)
	if err != nil {
		log.Error("remove stealth dns backup file /etc/resolv_stealthdns.conf fail: %v", err)
		return
	}
	_ = os.Chmod(dnsConfigFile, 0644)
}

func (h *LinuxHandler) backupDNSFile() error {
	_, err := os.Stat(dnsConfigFile)
	if err != nil {
		return err
	}
	if _, err = os.Stat(dnsBackupFile); err == nil {
		if err = os.Remove(dnsBackupFile); err != nil {
			return err
		}
	} else {
		if !os.IsNotExist(err) {
			return err
		}
	}

	err = os.Rename(dnsConfigFile, dnsBackupFile)

	return err
}

func (h *LinuxHandler) createStealthDNSFile() error {
	src, err := os.Open(dnsBackupFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(dnsConfigFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	scanner := bufio.NewScanner(src)
	writer := bufio.NewWriter(dst)
	defer writer.Flush()

	setStealthDNS := false

	for scanner.Scan() {
		originalLine := scanner.Text()
		if strings.HasSuffix(originalLine, stealthDNSPath) {
			continue
		}

		if strings.HasPrefix(originalLine, "nameserver") && !setStealthDNS {
			_, err = writer.WriteString(stealthDNSPath + "\n")
			if err != nil {
				return err
			}
			setStealthDNS = true
		}
		_, err = writer.WriteString(originalLine + "\n")
		if err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	_ = os.Chmod(dnsConfigFile, 0644)
	return nil
}

func (h *LinuxHandler) isAdminPermission() {
	h.isAdmin = os.Geteuid() == 0
}

func NewLinuxHandler() *LinuxHandler {
	handler := &LinuxHandler{
		dnsSetupFlag: false,
	}
	handler.isAdminPermission()

	return handler
}
