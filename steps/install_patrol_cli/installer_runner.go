package install_patrol_cli

import (
	get_cli_version "patrol_build_flutter/steps/install_patrol_cli/get_cli_version"
	install_cli_tool "patrol_build_flutter/steps/install_patrol_cli/install_cli_tool"

	v "github.com/Masterminds/semver/v3"
)

type InstallerRunner struct{}

func (p *InstallerRunner) GetPatrolCLIVersion() (*v.Version, error) {
	return get_cli_version.GetPatrolCLIVersion()
}

func (p *InstallerRunner) InstallPatrolCLI() error {
	_, err := install_cli_tool.InstallPatrolCLI(nil)
	return err
}
