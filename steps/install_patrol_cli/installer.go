package install_patrol_cli

import (
	"fmt"

	v "github.com/Masterminds/semver/v3"

	"patrol_build_flutter/utils/print"
)

type Installer interface {
	GetPatrolCLIVersion() (*v.Version, error)
	InstallPatrolCLI() error
}

func Run(installer Installer) (*v.Version, error) {
	if installer == nil {
		return nil, fmt.Errorf("installer is required")
	}

	print.StepInitiated("--- Checking if Patrol CLI is already installed ---")

	version, err := installer.GetPatrolCLIVersion()
	if err != nil {
		print.Warning("CLI is not installed, attempting installation...")
		if err := installer.InstallPatrolCLI(); err != nil {
			print.Errorf("❌ Installation failed: %s", err.Error())
			return nil, err
		}

		version, err = installer.GetPatrolCLIVersion()
		if err != nil {
			print.Errorf("❌ Failed to verify version after install: %s", err.Error())
			return nil, err
		}

		print.StepCompletedf("✅ PATROL CLI installed successfully. Version: %s\n", version.String())
		return version, nil
	}

	print.StepCompletedf("✅ Tool already installed. Version: %s\n", version.String())
	return version, nil
}
