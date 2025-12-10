package validate_versions

import (
	v "github.com/Masterminds/semver/v3"
)

type ValidateRunParams struct {
	FlutterVersion *v.Version
	CliVersion     *v.Version
	PatrolVersion  *v.Version
}

func CheckCompatibility(params ValidateRunParams) bool {
	flutterV := params.FlutterVersion
	patrolCLIV := params.CliVersion
	patrolV := params.PatrolVersion

	if flutterV == nil {
		panic("FlutterVersion cannot be nil in CheckCompatibility")
	}
	if patrolCLIV == nil {
		panic("panic: CliVersion cannot be nil in CheckCompatibility")
	}
	if patrolV == nil {
		panic("PatrolVersion cannot be nil in CheckCompatibility")
	}

	for _, entry := range CompatibilityTable {
		if isVersionInRange(patrolCLIV, entry.PatrolCLIRange) &&
			isVersionInRange(patrolV, entry.PatrolRange) &&
			flutterV.Equal(entry.FlutterVersion) {
			return true
		}
	}
	return false
}

func isVersionInRange(v *v.Version, r VersionRange) bool {
	return (v.Equal(r.Min) || v.GreaterThan(r.Min)) &&
		(v.Equal(r.Max) || v.LessThan(r.Max))
}
