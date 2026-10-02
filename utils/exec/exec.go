package exec

import (
	"bytes"
	"fmt"
	"os/exec"

	"patrol_build_flutter/commands"
)

// / Executes a command and returns its output as a string.
func Command(cmd commands.Command) (string, error) {
	command := exec.Command(cmd.Name, cmd.Args...)
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()

	if err != nil {
		return "", fmt.Errorf("failed to run %s %s: %w", command, cmd.Args, err)
	}

	return out.String(), nil
}
