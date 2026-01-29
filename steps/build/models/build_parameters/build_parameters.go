package build_parameters

import (
	"fmt"
	"os"
	"strings"

	build_constants "patrol_install/steps/build/constants"
)

// BuildParameters holds validated and formatted build configuration.
type BuildParameters struct {
	Target       string
	Platform     string
	BuildType    string
	Tags         string
	ExcludedTags string
	IsVerbose    string
}

// NewBuildParameters builds a BuildParameters struct from a map of environment variables.
func NewBuildParameters(envMap map[string]string) (*BuildParameters, error) {
	bp := &BuildParameters{}

	requiredFields := map[string]func(*BuildParameters, string) error{
		"platform":  SetPlatform,
		"buildType": SetBuildType,
	}

	optionalFields := map[string]func(*BuildParameters, string) error{
		"tags":         SetTags,
		"target":       SetTarget,
		"excludedTags": SetExcludedTags,
		"verbose":      SetVerbose,
	}

	// Apply required setters
	for key, setter := range requiredFields {
		val, ok := envMap[key]
		if !ok || strings.TrimSpace(val) == "" {
			return nil, fmt.Errorf("missing required field: %s", key)
		}
		if err := setter(bp, val); err != nil {
			return nil, err
		}
	}

	// Apply optional setters
	for key, setter := range optionalFields {
		if val, ok := envMap[key]; ok && strings.TrimSpace(val) != "" {
			if err := setter(bp, val); err != nil {
				return nil, err
			}
		}
	}

	return bp, nil
}

// Command constructs the final CLI command string based on the populated BuildParameters fields.
func (bp *BuildParameters) Command() []string {
	platform := os.Getenv(build_constants.Platform)
	buildType := os.Getenv(build_constants.BuildType)
	isiOS := platform != "android"
	isDebug := buildType == "debug"
	isiOSSimulator := isiOS && isDebug

	args := []string{}
	if bp.Target != "" {
		args = append(args, "--target", bp.Target)
	}
	if bp.Tags != "" {
		args = append(args, "--tags", bp.Tags)
	}
	if bp.ExcludedTags != "" {
		args = append(args, "--excludedTags", bp.ExcludedTags)
	}
	if bp.IsVerbose != "" {
		args = append(args, bp.IsVerbose)
	}

	buildTypeArg := "--" + bp.BuildType
	if isiOSSimulator {
		buildTypeArg += " --simulator"
	}

	buildCmd := func(platform, buildTypeArg string) string {
		return "patrol build " + platform + " " + buildTypeArg + " " + strings.Join(args, " ")
	}

	if bp.Platform == "both" {
		return []string{
			buildCmd("android", "--"+bp.BuildType),
			buildCmd("ios", buildTypeArg),
		}
	}

	cmd := buildCmd(bp.Platform, buildTypeArg)
	return []string{cmd}
}
