package builder

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"patrol_install/utils/print"
)

type Builder interface {
	BuildParametersFromEnv() ([]string, error)
}

func Run(installer Builder) error {
	print.StepInitiated("--- Starting Build Process ---")

	commands, err := installer.BuildParametersFromEnv()

	if err != nil {
		print.Error(fmt.Sprintf("❌ Failed to retrieve build commands: %s", err))
		return err
	}

	for _, cmd := range commands {
		print.Action(fmt.Sprintf("Executing build command: %s", cmd))

		if err := executeCommand(cmd); err != nil {
			print.Error(fmt.Sprintf("❌ Command failed: %s\n", err))
			return fmt.Errorf("build aborted: failed to execute '%s': %w", cmd, err)
		}

		print.Success(fmt.Sprintf("✅ Command '%s' executed successfully.\n", cmd))
	}

	print.StepCompleted("✅ All build commands executed successfully.")
	return nil
}

func executeCommand(command string) error {
	// Use 'sh -c' to allow complex shell expressions
	cmd := exec.Command("sh", "-c", command)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Stream output in real time
	go streamOutput(stdoutPipe)
	go streamOutput(stderrPipe)

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func streamOutput(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
