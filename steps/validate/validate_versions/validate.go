package validate_versions

import (
	"fmt"

	v "github.com/Masterminds/semver/v3"
)

type ValidateRunParams struct {
	FlutterVersion *v.Version
	CliVersion     *v.Version
	PatrolVersion  *v.Version
}

func CheckCompatibility(params ValidateRunParams) error {
	flutterV := params.FlutterVersion
	patrolCLIV := params.CliVersion
	patrolV := params.PatrolVersion

	for _, entry := range CompatibilityTable {
		if isVersionInRange(patrolCLIV, entry.PatrolCLIRange) &&
			isVersionInRange(patrolV, entry.PatrolRange) &&
			flutterV.GreaterThan(entry.MinFlutterVersion) || flutterV.Equal(entry.MinFlutterVersion) {
			return nil
		}
	}
	return fmt.Errorf("❌ Flutter %s, Patrol CLI %s and Patrol %s are not compatible", flutterV.String(), patrolCLIV.String(), patrolV.String())
}

func isVersionInRange(v *v.Version, r VersionRange) bool {
	return (v.Equal(r.Min) || v.GreaterThan(r.Min)) &&
		(v.Equal(r.Max) || v.LessThan(r.Max))
}
