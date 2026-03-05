//go:build !windows

package dev

import (
	"os"
	"os/exec"
	"syscall"
)

// Platform-specific signals to listen for
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// setSysProcAttr sets Unix-specific process attributes (process group).
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup kills the process and its children using process groups.
func killProcessGroup(cmd *exec.Cmd, force bool) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	sig := syscall.SIGTERM
	if force {
		sig = syscall.SIGKILL
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		// Fallback to killing just the process
		return cmd.Process.Kill()
	}

	return syscall.Kill(-pgid, sig)
}
