package builder

import (
	"fmt"

	createparams "patrol_install/steps/build/steps/create_parameters"
	"patrol_install/utils/print"
)

type BuilderRunner struct{}

func (p *BuilderRunner) BuildParametersFromEnv() ([]string, error) {
	command, err := createparams.BuildParametersFromEnv()
	if err != nil {
		print.Error(fmt.Sprintf("Build failed: %s", err))
		return nil, err
	}

	finalCommand := command.Command()
	if finalCommand == nil {
		buildErr := fmt.Errorf("no build commands generated: Command() returned nil")
		print.Error(buildErr.Error())
		return nil, buildErr
	}

	return finalCommand, nil
}
