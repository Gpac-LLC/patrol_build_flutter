package builder

import (
	"os"
	"testing"

	buildconstants "patrol_build_flutter/steps/build/constants"
)

func TestBuilderRunner_BuildParametersFromEnv_Success(t *testing.T) {
	t.Setenv(buildconstants.Platform, "android")
	t.Setenv(buildconstants.BuildType, "release")

	runner := &BuilderRunner{}
	cmds, err := runner.BuildParametersFromEnv()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cmds) == 0 {
		t.Fatal("expected at least one command, got empty slice")
	}
	if cmds[0] != "patrol build android --release" {
		t.Errorf("expected 'patrol build android --release', got %q", cmds[0])
	}
}

func TestBuilderRunner_BuildParametersFromEnv_MissingPlatform(t *testing.T) {
	t.Setenv(buildconstants.Platform, "")
	if err := os.Unsetenv(buildconstants.Platform); err != nil {
		t.Fatalf("unset platform: %v", err)
	}
	t.Setenv(buildconstants.BuildType, "release")

	runner := &BuilderRunner{}
	cmds, err := runner.BuildParametersFromEnv()
	if err == nil {
		t.Fatal("expected error when platform is missing, got nil")
	}
	if cmds != nil {
		t.Fatalf("expected nil commands on error, got %v", cmds)
	}
}
