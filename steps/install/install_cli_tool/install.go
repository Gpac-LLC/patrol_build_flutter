package install_cli_tool

import (
	"os"

	"patrol_install/commands"
	constants "patrol_install/steps/build/constants"
	"patrol_install/utils/exec"
	print "patrol_install/utils/print"
)

var patrolInstall = commands.PatrolInstall

func Install() (string, error) {

	customVersion := os.Getenv(constants.CustomPatrolCLIVersion)

	if customVersion == "" {
		print.Warning("Version was not provided. Using the latest version.")
	}

	command := command(customVersion)

	output, err := exec.Command(command)

	if err != nil {
		return output, err
	}

	return "", nil
}

func command(version string) commands.Command {
	// Ensure version is not empty before appending
	if version == "" {
		return patrolInstall
	}

	// Create a modified copy of patrolInstall with the version appended to Args
	cmd := patrolInstall.CopyWith(nil, append(patrolInstall.Args, version))
	return cmd
}
