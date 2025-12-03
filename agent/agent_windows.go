//go:build windows

package agent

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/OpenNHP/StealthDNS/common"
	"github.com/OpenNHP/opennhp/nhp/log"
	"golang.org/x/sys/windows"
)

type WindowsAgent struct {
	handle                windows.Handle
	nhpAgentInit          uintptr
	nhpAgentClose         uintptr
	nhpAgentKnockResource uintptr
	nhpFreeCString        uintptr
}

func (a *WindowsAgent) AgentInit(workingDir string, logLevel int) error {
	ptr, err := StringToPtr(workingDir)
	if err != nil {
		return err
	}
	_, _, errno := syscall.SyscallN(a.nhpAgentInit, ptr, uintptr(logLevel))
	if errno == 0 {
		return nil
	}
	return errno
}

func (a *WindowsAgent) AgentClose() error {
	_, _, _ = syscall.SyscallN(a.nhpAgentClose)
	return windows.FreeLibrary(a.handle)
}

func (a *WindowsAgent) AgentKnockResource(aspId, resId, serverIp, serverHostname string, serverPort int) (string, error) {
	aspIdPtr, err := StringToPtr(aspId)
	if err != nil {
		return "", err
	}
	resIdPtr, err := StringToPtr(resId)
	if err != nil {
		return "", err
	}
	serverIpPtr, err := StringToPtr(serverIp)
	if err != nil {
		return "", err
	}
	hostNamePtr, err := StringToPtr(serverHostname)
	if err != nil {
		return "", err
	}
	ret, _, errno := syscall.SyscallN(a.nhpAgentKnockResource, aspIdPtr, resIdPtr, serverIpPtr, hostNamePtr, uintptr(serverPort))
	defer func() {
		_, _, err := syscall.SyscallN(a.nhpFreeCString, ret)
		if err != 0 {
			log.Warning("Failed to free memory")
		}
	}()
	if errno != 0 {
		return "", errno
	}
	result := PtrToString(ret)

	return result, err
}

func newNhpAgent(libPath string) (*WindowsAgent, error) {
	fileName := filepath.Join(libPath, "sdk", "nhp-agent.dll")
	_, err := os.Stat(fileName)
	if err != nil {
		log.Error("find nhp-agent.dll fail: %v", err)
		return nil, err
	}
	handle, err := windows.LoadLibrary(fileName)
	if err != nil {
		log.Error("load nhp-agent.dll fail: %v", err)
		return nil, err
	}
	a := &WindowsAgent{handle: handle}
	a.nhpAgentInit, err = windows.GetProcAddress(a.handle, common.AgentInit)
	if err != nil {
		log.Error("load nhp_agent_init func fail: %v", err)
		return nil, err
	}

	a.nhpAgentClose, err = windows.GetProcAddress(a.handle, common.AgentClose)
	if err != nil {
		log.Error("load nhp_agent_close func fail: %v", err)
		return nil, err
	}

	a.nhpAgentKnockResource, err = windows.GetProcAddress(a.handle, common.AgentKnockResource)
	if err != nil {
		log.Error("load nhp_agent_knock_resource func fail: %v", err)
		return nil, err
	}

	a.nhpFreeCString, err = windows.GetProcAddress(a.handle, common.AgentFreeCString)
	if err != nil {
		log.Error("load nhp_agent_init func fail: %v", err)
		return nil, err
	}

	return a, nil
}
