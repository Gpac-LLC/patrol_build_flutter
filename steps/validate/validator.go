package validate

import (
	"errors"
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
		return err
	}

	print.StepCompleted("✅ Flutter Version: " + flutterVersion.String() + "\n")

	print.StepInitiated("--- Getting Patrol Version ---")
	patrolVersion, patrolErr := runner.GetPatrolVersion()

	if patrolErr != nil {
		print.Warning("❌ Failed to get Patrol version")
		print.Error(patrolErr.Error())
		return patrolErr
	}

	print.StepCompleted("✅ Patrol Version: " + patrolVersion.String() + "\n")

	validatorParams := versions.ValidateRunParams{
		FlutterVersion: flutterVersion,
		CliVersion:     params.CliVersion,
		PatrolVersion:  patrolVersion,
	}

	print.StepInitiated("--- Checking Compatibility ---")
	isCompatible := versions.CheckCompatibility(validatorParams)

	if isCompatible {
		message := fmt.Sprintf("✅ Flutter %s, Patrol CLI %s and Patrol %s are compatible",
			flutterVersion.String(), params.CliVersion.String(), patrolVersion.String())
		print.StepCompleted(message)
		return nil
	}
	errorMessage := fmt.Sprintf("❌ Flutter %s, Patrol CLI %s and Patrol %s are not compatible",
		flutterVersion.String(), params.CliVersion.String(), patrolVersion.String())
	print.Error(errorMessage)
	return errors.New(errorMessage)

}
