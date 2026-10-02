package validate

import (
	"fmt"

	v "github.com/Masterminds/semver/v3"

	versions "patrol_build_flutter/steps/validate/validate_versions"
	"patrol_build_flutter/utils/print"
)

type Validator interface {
	GetFlutterVersion() (*v.Version, error)
	GetPatrolVersion() (*v.Version, error)
}

type ValidatorRunParams struct {
	Runner     Validator
	CliVersion *v.Version
}

func Run(params ValidatorRunParams) error {
	if params.Runner == nil {
		return fmt.Errorf("runner is required for validation")
	}
	if params.CliVersion == nil {
		return fmt.Errorf("CLI version is required for validation — was the installation step successful?")
	}

	runner := params.Runner

	print.StepInitiated("--- Getting Flutter Version ---")

	flutterVersion, err := runner.GetFlutterVersion()
	if err != nil {
		print.Warning("❌ Failed to get Flutter version")
		print.Error(err.Error())
		return err
	}

	print.StepCompletedf("✅ Flutter Version: %s\n", flutterVersion.String())

	print.StepInitiated("--- Getting Patrol Version ---")
	patrolVersion, patrolErr := runner.GetPatrolVersion()

	if patrolErr != nil {
		print.Warning("❌ Failed to get Patrol version")
		print.Error(patrolErr.Error())
		return patrolErr
	}

	print.StepCompletedf("✅ Patrol Version: %s\n", patrolVersion.String())

	validatorParams := versions.ValidateRunParams{
		FlutterVersion: flutterVersion,
		CliVersion:     params.CliVersion,
		PatrolVersion:  patrolVersion,
	}

	print.StepInitiated("--- Checking Compatibility ---")
	isCompatible := versions.CheckCompatibility(validatorParams)

	if isCompatible {
		print.StepCompletedf("✅ Flutter %s, Patrol CLI %s and Patrol %s are compatible",
			flutterVersion.String(), params.CliVersion.String(), patrolVersion.String())
		return nil
	}

	compatErr := fmt.Errorf("❌ Flutter %s, Patrol CLI %s and Patrol %s are not compatible",
		flutterVersion.String(), params.CliVersion.String(), patrolVersion.String())
	print.Error(compatErr.Error())
	return compatErr
}
