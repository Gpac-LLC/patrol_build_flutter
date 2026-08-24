package install_cli_tool

import (
	"os"

	"patrol_install/commands"
	constants "patrol_install/steps/build/constants"
	"patrol_install/utils/exec"
	print "patrol_install/utils/print"
)

var patrolInstall = commands.PatrolInstall

type CommandExecutor func(cmd commands.Command) (string, error)

// InstallPatrolCLI installs the Patrol CLI, using a custom version if provided.
// The executor parameter allows for dependency injection in tests. Pass nil to use the default executor.
func InstallPatrolCLI(executor CommandExecutor) (string, error) {
	customVersion := os.Getenv(constants.CustomPatrolCLIVersion)

	if customVersion == "" {
		print.Warning("Version was not provided. Using the latest version.")
	} else {
		print.Actionf("Installing custom version provided: %s", customVersion)
	}

	installCmd := buildInstallCommand(customVersion)

	cmdExecutor := executor
	if cmdExecutor == nil {
		cmdExecutor = exec.Command
	}

	output, err := cmdExecutor(installCmd)
	if err != nil {
		return output, err
	}

	print.Success("Patrol CLI installed successfully.")
	return output, nil
}

// buildInstallCommand returns the appropriate Command struct based on the version.
func buildInstallCommand(version string) commands.Command {
	if version == "" {
		return patrolInstall
	}
	return patrolInstall.CopyWith(nil, append(patrolInstall.Args, version))
}
