package builder

import (
	"errors"
	"testing"
)

type mockBuilder struct {
	buildParametersFromEnvFunc func() ([]string, error)
}

func (m *mockBuilder) BuildParametersFromEnv() ([]string, error) {
	if m.buildParametersFromEnvFunc != nil {
		return m.buildParametersFromEnvFunc()
	}
	return []string{"echo hello"}, nil
}

func TestBuilder_NilBuilder(t *testing.T) {
	err := Run(nil)
	if err == nil {
		t.Fatal("expected error when builder is nil, got nil")
	}
}

func TestBuilder_BuildParametersFromEnvError(t *testing.T) {
	mock := &mockBuilder{
		buildParametersFromEnvFunc: func() ([]string, error) {
			return nil, errors.New("failed to parse build parameters")
		},
	}
	err := Run(mock)
	if err == nil {
		t.Fatal("expected error when BuildParametersFromEnv fails, got nil")
	}
}

func TestBuilder_EmptyCommandsError(t *testing.T) {
	mock := &mockBuilder{
		buildParametersFromEnvFunc: func() ([]string, error) {
			return []string{}, nil
		},
	}
	err := Run(mock)
	if err == nil {
		t.Fatal("expected error when commands list is empty, got nil")
	}
}

func TestBuilder_SuccessfulCommand(t *testing.T) {
	mock := &mockBuilder{
		buildParametersFromEnvFunc: func() ([]string, error) {
			return []string{"echo hello"}, nil
		},
	}
	err := Run(mock)
	if err != nil {
		t.Fatalf("expected successful command execution, got error: %v", err)
	}
}

func TestBuilder_FailedCommand(t *testing.T) {
	mock := &mockBuilder{
		buildParametersFromEnvFunc: func() ([]string, error) {
			return []string{"exit 1"}, nil
		},
	}
	err := Run(mock)
	if err == nil {
		t.Fatal("expected error when command fails, got nil")
	}
}
