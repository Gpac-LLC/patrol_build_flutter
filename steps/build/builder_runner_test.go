package builder

import (
	"os"
	"testing"

	build_constants "patrol_build_flutter/steps/build/constants"
)

func TestBuilderRunner_BuildParametersFromEnv_Success(t *testing.T) {
	os.Setenv(build_constants.Platform, "android")
	os.Setenv(build_constants.BuildType, "release")
	defer func() {
		os.Unsetenv(build_constants.Platform)
		os.Unsetenv(build_constants.BuildType)
	}()

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
	os.Unsetenv(build_constants.Platform)
	os.Setenv(build_constants.BuildType, "release")
	defer os.Unsetenv(build_constants.BuildType)

	runner := &BuilderRunner{}
	cmds, err := runner.BuildParametersFromEnv()
	if err == nil {
		t.Fatal("expected error when platform is missing, got nil")
	}
	if cmds != nil {
		t.Fatalf("expected nil commands on error, got %v", cmds)
	}
}
