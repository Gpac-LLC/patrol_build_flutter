package builder

import (
	"fmt"

	createparams "patrol_build_flutter/steps/build/steps/create_parameters"
	"patrol_build_flutter/utils/print"
)

type BuilderRunner struct{}

func (p *BuilderRunner) BuildParametersFromEnv() ([]string, error) {
	command, err := createparams.BuildParametersFromEnv()
	if err != nil {
		print.Errorf("Build failed: %s", err)
		return nil, err
	}

	finalCommand := command.Command()
	if len(finalCommand) == 0 {
		buildErr := fmt.Errorf("no build commands generated: Command() returned empty")
		print.Error(buildErr.Error())
		return nil, buildErr
	}

	return finalCommand, nil
}
