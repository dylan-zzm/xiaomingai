package gateway

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// RunHelper executes the chatlog helper CLI tool with the given mode and arguments.
// It automatically elevates privileges: sudo on macOS, runas on Windows.
func RunHelper(mode string, kv map[string]string) error {
	helperPath := "./helper"

	args := []string{"--mode", mode}
	for k, v := range kv {
		if v != "" {
			args = append(args, fmt.Sprintf("--%s", k), v)
		}
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		sudoArgs := append([]string{helperPath}, args...)
		cmd = exec.Command("sudo", sudoArgs...)
	} else if runtime.GOOS == "windows" {
		// Note: /savecred might cache credentials, which can be a security consideration.
		argsLine := strings.Join(append([]string{helperPath}, args...), " ")
		cmd = exec.Command("runas", "/user:Administrator", "/savecred", argsLine)
	} else {
		cmd = exec.Command(helperPath, args...)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}