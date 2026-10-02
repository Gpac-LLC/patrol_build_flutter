package validate

import (
	v "github.com/Masterminds/semver/v3"

	flutter "patrol_build_flutter/steps/validate/get_flutter_version"
	patrol "patrol_build_flutter/steps/validate/get_patrol_version"
)

type ValidatorRunner struct{}

func (p *ValidatorRunner) GetFlutterVersion() (*v.Version, error) {
	return flutter.GetFlutterVersion(flutter.FlutterVersionCmd)
}

func (p *ValidatorRunner) GetPatrolVersion() (*v.Version, error) {
	return patrol.GetPatrolVersion(patrol.FlutterPubDepsCmd)
}
