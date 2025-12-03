//go:build linux || darwin

package agent

/*
// Linux .so
#cgo linux LDFLAGS: -L../sdk -l:nhp-agent.so -Wl,-rpath,'$ORIGIN/sdk'
// macOS  .dylib
#cgo darwin LDFLAGS: ${SRCDIR}/../sdk/nhp-agent.dylib
#include <stdlib.h>
#include <stdbool.h>

// declare functions from a dynamic library
bool nhp_agent_init(const char* workingDir, int logLevel);
void nhp_agent_close();
void nhp_free_cstring(const char* ptr);
char* nhp_agent_knock_resource(const char* aspId, const char* resId, const char* serverIp, const char* serverHostName, int serverPort);


*/
import "C"

import (
	"errors"
	"unsafe"
)

type UnixAgent struct {
}

func (a *UnixAgent) AgentInit(workingDir string, logLevel int) error {
	dir := C.CString(workingDir)
	defer C.free(unsafe.Pointer(dir))
	flag := C.nhp_agent_init(dir, C.int(logLevel))
	if bool(flag) {
		return nil
	}
	return errors.New("failed to initialize the nhp-agent")
}

func (a *UnixAgent) AgentClose() error {
	C.nhp_agent_close()
	return nil
}

func (a *UnixAgent) AgentKnockResource(aspId, resId, serverIp, serverHostname string, serverPort int) (string, error) {
	cAspId := C.CString(aspId)
	defer C.free(unsafe.Pointer(cAspId))
	cResId := C.CString(resId)
	defer C.free(unsafe.Pointer(cResId))
	cServerIp := C.CString(serverIp)
	defer C.free(unsafe.Pointer(cServerIp))
	cServerHostname := C.CString(serverHostname)
	defer C.free(unsafe.Pointer(cServerHostname))

	result := C.nhp_agent_knock_resource(cAspId, cResId, cServerIp, cServerHostname, C.int(serverPort))
	defer C.nhp_free_cstring(result)
	goResult := C.GoString(result)
	return goResult, nil
}

func newNhpAgent(libPath string) (*UnixAgent, error) {
	return &UnixAgent{}, nil
}
