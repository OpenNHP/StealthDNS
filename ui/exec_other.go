//go:build !windows

package main

import (
	"os/exec"
)

// hideWindow empty implementation for non-Windows platforms
func hideWindow(cmd *exec.Cmd) {
	// No need to hide window on non-Windows platforms
}

