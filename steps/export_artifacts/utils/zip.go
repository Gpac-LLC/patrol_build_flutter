package export_artifacts_utils

import (
	"fmt"

	"patrol_build_flutter/commands"
	"patrol_build_flutter/utils/exec"
)

var compressIOSFiles = commands.CompressIOSFiles

type CommandExecutor func(cmd commands.Command) (string, error)

// ZipFiles builds and executes a zip command for the given input paths.
func ZipFiles(zipPath string, inputPaths []string, executor CommandExecutor) (string, error) {
	if zipPath == "" {
		return "", fmt.Errorf("zip path is empty")
	}
	if len(inputPaths) == 0 {
		return "", fmt.Errorf("no input paths to zip")
	}

	cmdArgs := append([]string{"-r", zipPath}, inputPaths...)
	cmd := compressIOSFiles.CopyWith(nil, cmdArgs)

	run := executor
	if run == nil {
		run = exec.Command
	}
	if _, err := run(cmd); err != nil {
		return "", err
	}
	return zipPath, nil
}
