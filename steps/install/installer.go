package install

import (
	v "github.com/Masterminds/semver/v3"

	"patrol_install/utils/print"
)

type Installer interface {
	GetVersion() (*v.Version, error)
	Install() error
}

func Run(installer Installer) (*v.Version, error) {
	print.StepInitiated("--- Checking if Patrol CLI is already installed ---")

	version, err := installer.GetVersion()
	if err != nil {
		print.Warning("CLI is not installed, attempting installation...")
		if err := installer.Install(); err != nil {
			print.Error("❌ Installation failed: " + err.Error())
			return nil, err
		}

		version, err = installer.GetVersion()
		if err != nil {
			print.Error("❌ Failed to verify version after install: " + err.Error())
			return nil, err
		}

		print.StepCompleted("✅ PATROL CLI installed successfully. Version: " + version.String() + "\n")
		return version, nil
	}

	print.StepCompleted("✅ Tool already installed. Version: " + version.String() + "\n")
	return version, nil
}
