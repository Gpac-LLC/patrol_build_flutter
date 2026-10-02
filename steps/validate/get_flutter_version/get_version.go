package get_flutter_version

import (
	"fmt"
	"strings"

	v "github.com/Masterminds/semver/v3"

	"patrol_build_flutter/commands"
	regex "patrol_build_flutter/constants"
	"patrol_build_flutter/utils/exec"
)

var FlutterVersionCmd = commands.FlutterVersion

// CleanVersion receives the command output string and extracts the version substring.
func CleanVersion(output string) (string, error) {
	re := regex.Version("Flutter")
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		return cleanVersion(match[1]), nil
	}
	return "", fmt.Errorf("could not find version in output")
}

// ParseVersion receives the cleaned version string and parses it into a semver.Version.
func ParseVersion(versionStr string) (*v.Version, error) {
	parsedVersion, err := v.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid semantic version: %w", err)
	}
	return parsedVersion, nil
}

func GetFlutterVersion(cmd commands.Command) (*v.Version, error) {
	output, err := exec.Command(cmd)
	if err != nil {
		return nil, err
	}

	cleaned, err := CleanVersion(output)
	if err != nil {
		return nil, err
	}

	parsedVersion, err := ParseVersion(cleaned)
	if err != nil {
		return nil, err
	}
	return parsedVersion, nil
}

func cleanVersion(version string) string {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	return version
}
