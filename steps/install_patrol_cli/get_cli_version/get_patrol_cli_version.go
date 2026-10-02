package get_cli_version

import (
	"fmt"
	"strings"

	v "github.com/Masterminds/semver/v3"

	"patrol_build_flutter/commands"
	regex "patrol_build_flutter/constants"
	"patrol_build_flutter/utils/exec"
)

var patrolDoctor = commands.PatrolDoctor

func GetPatrolCLIVersion() (*v.Version, error) {
	output, err := exec.Command(patrolDoctor)
	if err != nil {
		return nil, err
	}

	// Use regex to extract the version
	re := regex.Version("Patrol CLI Version")
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		parsedVersion, err := v.NewVersion(cleanVersion(match[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid semantic version: %w", err)
		}
		return parsedVersion, nil
	}

	return nil, fmt.Errorf("could not find version in output")
}

func cleanVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	return version
}
