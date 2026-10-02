package export_artifacts_utils

import (
	"testing"

	"patrol_build_flutter/commands"
)

type commandExecutorStub struct {
	called bool
	cmd    commands.Command
	err    error
}

func (s *commandExecutorStub) Run(cmd commands.Command) (string, error) {
	s.called = true
	s.cmd = cmd
	return "", s.err
}

func TestZipFiles(t *testing.T) {
	// GIVEN a zip path and input paths
	zipPath := "/tmp/ios_tests.zip"
	inputs := []string{"/tmp/BuildDir", "/tmp/Runner.xctestrun"}
	stub := &commandExecutorStub{}

	// WHEN building and executing the zip command
	result, err := ZipFiles(zipPath, inputs, stub.Run)

	// THEN the command is executed with correct args
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !stub.called {
		t.Fatalf("expected executor to be called")
	}
	if result != zipPath {
		t.Fatalf("expected zip path %s, got %s", zipPath, result)
	}
	expectedArgs := []string{"-r", zipPath, inputs[0], inputs[1]}
	if len(stub.cmd.Args) != len(expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, stub.cmd.Args)
	}
	for i, arg := range expectedArgs {
		if stub.cmd.Args[i] != arg {
			t.Fatalf("expected arg %s at %d, got %s", arg, i, stub.cmd.Args[i])
		}
	}
}

func TestZipFiles_EmptyZipPath(t *testing.T) {
	// GIVEN an empty zip path
	stub := &commandExecutorStub{}

	// WHEN building the command
	_, err := ZipFiles("", []string{"/tmp/BuildDir"}, stub.Run)

	// THEN it fails before executing
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if stub.called {
		t.Fatalf("expected executor not to be called")
	}
}

func TestZipFiles_NoInputs(t *testing.T) {
	// GIVEN no input paths
	stub := &commandExecutorStub{}

	// WHEN building the command
	_, err := ZipFiles("/tmp/ios_tests.zip", nil, stub.Run)

	// THEN it fails before executing
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if stub.called {
		t.Fatalf("expected executor not to be called")
	}
}
