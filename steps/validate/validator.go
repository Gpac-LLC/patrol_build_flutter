package validate

import (
	"fmt"

	v "github.com/Masterminds/semver/v3"

	versions "patrol_install/steps/validate/validate_versions"
	"patrol_install/utils/print"
)

type Validator interface {
	GetVersion() (*v.Version, error)
	GetPatrolVersion() (*v.Version, error)
}

type ValidatorRunParams struct {
	Runner     Validator
	CliVersion *v.Version
}

func Run(params ValidatorRunParams) error {
	runner := params.Runner

	print.StepInitiated("--- Getting Flutter Version ---")

	flutterVersion, err := runner.GetVersion()
	if err != nil {
		print.Warning("❌ Failed to get Flutter version")
		print.Error(err.Error())
	}

	print.StepCompleted("✅ Flutter Version: " + flutterVersion.String() + "\n")

	print.StepInitiated("--- Getting Patrol Version ---")
	patrolVersion, patrolErr := runner.GetPatrolVersion()

	if patrolErr != nil {
		print.Warning("❌ Failed to get Patrol version")
		print.Error(patrolErr.Error())
	}

	print.StepCompleted("✅ Patrol Version: " + patrolVersion.String() + "\n")

	validatorParams := versions.ValidateRunParams{
		FlutterVersion: flutterVersion,
		CliVersion:     params.CliVersion,
		PatrolVersion:  patrolVersion,
	}

	print.StepInitiated("--- Checking Compatibility ---")
	validationError := versions.CheckCompatibility(validatorParams)
	if validationError != nil {
		print.Warning("❌ Failed to check compatibility")
		print.Error(validationError.Error())
	}

	message := fmt.Sprintf("✅ Flutter %s, Patrol CLI %s and Patrol %s are compatible", flutterVersion.String(), params.CliVersion.String(), patrolVersion.String())
	print.StepCompleted(message)

	return nil
}
