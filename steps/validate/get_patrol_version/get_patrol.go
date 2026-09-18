package get_patrol_version

import (
	"bufio"
	"fmt"
	"strings"

	v "github.com/Masterminds/semver/v3"

	"patrol_build_flutter/commands"
	commands_utils "patrol_build_flutter/commands/utils"
	"patrol_build_flutter/utils/exec"
)

var FlutterPubDepsCmd = commands.FlutterPubDependencies

func GetPatrolVersion(cmd commands.Command) (*v.Version, error) {

	if !commands_utils.IsSameCommand(cmd, FlutterPubDepsCmd) {
		return nil, fmt.Errorf("should use FlutterPubDependencies command")
	}

	// Todo: create execture function that returns the log
	output, err := exec.Command(cmd)
	if err != nil {
		return nil, err
	}

	version, err := GetPatrolVersionFromLog(output)
	if err != nil {
		return nil, fmt.Errorf("could not find version in output")
	}

	return version, nil
}

func GetPatrolVersionFromLog(log string) (*v.Version, error) {
	scanner := bufio.NewScanner(strings.NewReader(log))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "- patrol ") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				rawVersion := fields[2]
				version, err := v.NewVersion(rawVersion)
				if err != nil {
					return nil, fmt.Errorf("invalid version format: %v", err)
				}
				return version, nil
			}
		}
	}
	return nil, fmt.Errorf("patrol package not found")
}
