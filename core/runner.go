// Package core provides shared services used by all modules.
// The most important one is Runner — every module calls RunCommand()
// instead of calling the OS directly. This keeps things clean and testable.
package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommandResult holds the output of any shell command we run.
type CommandResult struct {
	Output string // stdout combined with stderr
	Err    error  // non-nil if the command failed
}

// RunCommand runs a shell command and returns its output.
// We use /bin/sh -c so we can pass pipes and multi-part commands.
// Example: RunCommand("df -h /") or RunCommand("pm2 list")
func RunCommand(command string) CommandResult {
	// Run via shell so we can use pipes, redirects, etc.
	cmd := exec.Command("/bin/sh", "-c", command)

	// Combine stdout and stderr into one buffer so we see everything
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return CommandResult{
		Output: strings.TrimSpace(out.String()),
		Err:    err,
	}
}

// RunCommandLines is like RunCommand but returns each output line separately.
// Useful when you want to loop over results (like a list of processes).
func RunCommandLines(command string) []string {
	result := RunCommand(command)
	if result.Err != nil {
		return []string{}
	}
	if result.Output == "" {
		return []string{}
	}
	return strings.Split(result.Output, "\n")
}

// CommandExists checks if a program is installed on the system.
// Example: CommandExists("nginx") returns true if nginx is installed.
func CommandExists(name string) bool {
	result := RunCommand(fmt.Sprintf("which %s", name))
	return result.Err == nil && result.Output != ""
}
