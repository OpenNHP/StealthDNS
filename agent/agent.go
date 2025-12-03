package agent

import (
	"fmt"
	"syscall"
	"unsafe"
)

type NhpAgent interface {
	AgentInit(workingDir string, logLevel int) error
	AgentClose() error
	AgentKnockResource(aspId, resId, serverIp, serverHostname string, serverPort int) (string, error)
}

func NewNhpAgent(libPath string) (NhpAgent, error) {
	return newNhpAgent(libPath)
}

func StringToPtr(s string) (uintptr, error) {
	if s == "" {
		return 0, nil
	}

	ptr, err := syscall.BytePtrFromString(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to pointer: %v", err)
	}
	return uintptr(unsafe.Pointer(ptr)), nil
}

func PtrToString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}

	length := 0
	tempPtr := ptr
	for {
		b := *(*byte)(unsafe.Pointer(tempPtr))
		if b == 0 {
			break
		}
		length++
		tempPtr++
	}

	if length == 0 {
		return ""
	}

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
	}

	return string(bytes)
}
