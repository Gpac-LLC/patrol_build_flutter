package validate

import (
	"errors"
	"testing"

	v "github.com/Masterminds/semver/v3"
)

type mockValidatorRunner struct {
	getFlutterVersionFunc func() (*v.Version, error)
	getPatrolVersionFunc  func() (*v.Version, error)
}

func (m *mockValidatorRunner) GetFlutterVersion() (*v.Version, error) {
	if m.getFlutterVersionFunc != nil {
		return m.getFlutterVersionFunc()
	}
	return v.MustParse("3.32.0"), nil
}

func (m *mockValidatorRunner) GetPatrolVersion() (*v.Version, error) {
	if m.getPatrolVersionFunc != nil {
		return m.getPatrolVersionFunc()
	}
	return v.MustParse("3.20.0"), nil
}

func TestValidator_NilRunner(t *testing.T) {
	err := Run(ValidatorRunParams{
		Runner:     nil,
		CliVersion: v.MustParse("3.11.0"),
	})
	if err == nil {
		t.Fatal("expected error when runner is nil, got nil")
	}
}

func TestValidator_NilCliVersion(t *testing.T) {
	mock := &mockValidatorRunner{}
	err := Run(ValidatorRunParams{
		Runner:     mock,
		CliVersion: nil,
	})
	if err == nil {
		t.Fatal("expected error when CliVersion is nil, got nil")
	}
}

func TestValidator_GetFlutterVersionError(t *testing.T) {
	mock := &mockValidatorRunner{
		getFlutterVersionFunc: func() (*v.Version, error) {
			return nil, errors.New("failed to get flutter version")
		},
	}
	err := Run(ValidatorRunParams{
		Runner:     mock,
		CliVersion: v.MustParse("3.11.0"),
	})
	if err == nil {
		t.Fatal("expected error when flutter version check fails, got nil")
	}
}

func TestValidator_GetPatrolVersionError(t *testing.T) {
	mock := &mockValidatorRunner{
		getPatrolVersionFunc: func() (*v.Version, error) {
			return nil, errors.New("failed to get patrol version")
		},
	}
	err := Run(ValidatorRunParams{
		Runner:     mock,
		CliVersion: v.MustParse("3.11.0"),
	})
	if err == nil {
		t.Fatal("expected error when patrol version check fails, got nil")
	}
}

func TestValidator_CompatibleVersions(t *testing.T) {
	mock := &mockValidatorRunner{
		getFlutterVersionFunc: func() (*v.Version, error) {
			return v.MustParse("3.32.0"), nil
		},
		getPatrolVersionFunc: func() (*v.Version, error) {
			return v.MustParse("3.20.0"), nil
		},
	}
	err := Run(ValidatorRunParams{
		Runner:     mock,
		CliVersion: v.MustParse("3.11.0"),
	})
	if err != nil {
		t.Fatalf("expected compatible versions to succeed, got error: %v", err)
	}
}

func TestValidator_IncompatibleVersions(t *testing.T) {
	mock := &mockValidatorRunner{
		getFlutterVersionFunc: func() (*v.Version, error) {
			return v.MustParse("3.32.0"), nil
		},
		getPatrolVersionFunc: func() (*v.Version, error) {
			return v.MustParse("1.0.0"), nil
		},
	}
	err := Run(ValidatorRunParams{
		Runner:     mock,
		CliVersion: v.MustParse("1.0.0"),
	})
	if err == nil {
		t.Fatal("expected incompatible versions to return error, got nil")
	}
}
