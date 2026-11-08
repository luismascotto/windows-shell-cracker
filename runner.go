package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// launchExplorerShell opens the Windows Explorer view for a shell: URI without blocking the TUI.
func launchExplorerShell(shellCommand string) error {
	if shellCommand == "" {
		return errNoSelection
	}
	switch runtime.GOOS {
	case "windows":
		// start "" explorer <uri> — detached child; TUI keeps running.
		cmd := exec.Command("cmd", "/c", "start", "", "explorer", shellCommand)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Start()
	default:
		// Dev / non-Windows: log intent so the TUI can be exercised elsewhere.
		fmt.Fprintf(os.Stderr, "[windows-shell-cracker] would run: start explorer %s\n", shellCommand)
		return nil
	}
}
