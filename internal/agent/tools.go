package agent

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// KubectlTool executes a kubectl command safely by enforcing allowed subcommands.
func KubectlTool(ctx context.Context, subcommand string, args ...string) (string, error) {
	allowedSubcommands := map[string]bool{
		"get":      true,
		"describe": true,
		"logs":     true,
		"events":   true, // though get events is common
	}

	if !allowedSubcommands[subcommand] {
		return "", fmt.Errorf("subcommand %s is not allowed for security reasons", subcommand)
	}

	cmdArgs := append([]string{subcommand}, args...)
	cmd := exec.CommandContext(ctx, "kubectl", cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("kubectl error: %v\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
