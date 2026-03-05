//go:build windows

package dev

import (
	"os"
	"os/exec"
)

// Platform-specific signals to listen for (Windows only has Interrupt)
var shutdownSignals = []os.Signal{os.Interrupt}

// setSysProcAttr sets Windows-specific process attributes.
// On Windows, we don't need Setpgid as it doesn't exist.
func setSysProcAttr(cmd *exec.Cmd) {
	// No special attributes needed on Windows
}

// killProcessGroup kills the process on Windows.
// Windows doesn't have process groups like Unix, so we just kill the process directly.
func killProcessGroup(cmd *exec.Cmd, force bool) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}
