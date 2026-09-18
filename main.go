package main

import (
	"os"

	build "patrol_build_flutter/steps/build"
	"patrol_build_flutter/steps/export_artifacts"
	"patrol_build_flutter/steps/install_patrol_cli"
	"patrol_build_flutter/steps/validate"
	"patrol_build_flutter/utils/print"
)

func exitOnError(step string, err error) {
	if err != nil {
		print.Errorf("❌ %s failed: %s", step, err)
		os.Exit(1)
	}
}

func main() {
	cliVersion, err := install_patrol_cli.Run(&install_patrol_cli.InstallerRunner{})
	exitOnError("Install Patrol CLI", err)
	print.Successf("✅ Patrol CLI installed: %s", cliVersion.String())

	err = validate.Run(validate.ValidatorRunParams{
		Runner:     &validate.ValidatorRunner{},
		CliVersion: cliVersion,
	})
	exitOnError("Version validation", err)

	err = build.Run(&build.BuilderRunner{})
	exitOnError("Build", err)

	err = export_artifacts.Run(&export_artifacts.ExporterRunner{})
	exitOnError("Export artifacts", err)
}
