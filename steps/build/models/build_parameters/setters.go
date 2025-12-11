package build_parameters

import (
	"errors"
	"fmt"
	"strings"
)

// SetPlatform sets the build platform. Accepted: "android", "iOS".
func SetPlatform(bp *BuildParameters, value string) error {
	var platform = strings.ToLower(value)
	switch platform {
	case "android", "ios":
		bp.Platform = platform
		return nil
	case "both":
		bp.Platform = value
		return nil
	default:
		return errors.New("invalid platform: expected 'android', 'ios' or 'both'")
	}
}

// SetTarget sets the target. Required and must not be empty.
func SetTarget(bp *BuildParameters, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("target cannot be empty")
	}
	bp.Target = value
	return nil
}

func SetBuildType(bp *BuildParameters, value string) error {
	switch value {
	case "release", "simulator":
		bp.BuildType = value
		return nil
	default:
		return errors.New("invalid build type: expected 'release' or 'simulator'")
	}
}

func SetTags(bp *BuildParameters, value string) error {
	bp.Tags = formatTags(value)
	return nil
}

func SetExcludedTags(bp *BuildParameters, value string) error {
	bp.ExcludedTags = formatTags(value)
	return nil
}

func SetVerbose(bp *BuildParameters, value string) error {
	return setFlag(value, "--verbose", &bp.IsVerbose, "verbose")
}

// formatTags converts comma-separated values to '( tag1 && tag2 )' format.
func formatTags(input string) string {
	tags := strings.Split(input, ",")
	var trimmed []string
	for _, t := range tags {
		if tag := strings.TrimSpace(t); tag != "" {
			trimmed = append(trimmed, tag)
		}
	}
	if len(trimmed) == 0 {
		return ""
	}
	return "'( " + strings.Join(trimmed, " && ") + " )'"
}

func setFlag(value, flag string, target *string, name string) error {
	switch strings.ToLower(value) {
	case "true":
		*target = flag
	case "false":
		*target = ""
	default:
		return fmt.Errorf("invalid value for %s: expected 'true' or 'false'", name)
	}
	return nil
}
