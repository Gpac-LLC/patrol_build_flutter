package install_patrol_cli

import (
	"errors"
	"testing"

	v "github.com/Masterminds/semver/v3"
)

type mockInstaller struct {
	getPatrolCLIVersionFunc func() (*v.Version, error)
	installPatrolCLIFunc    func() error
}

func (m *mockInstaller) GetPatrolCLIVersion() (*v.Version, error) {
	if m.getPatrolCLIVersionFunc != nil {
		return m.getPatrolCLIVersionFunc()
	}
	return v.MustParse("3.11.0"), nil
}

func (m *mockInstaller) InstallPatrolCLI() error {
	if m.installPatrolCLIFunc != nil {
		return m.installPatrolCLIFunc()
	}
	return nil
}

func TestInstaller_NilInstaller(t *testing.T) {
	version, err := Run(nil)
	if err == nil {
		t.Fatal("expected error when installer is nil, got nil")
	}
	if version != nil {
		t.Fatalf("expected version to be nil, got %v", version)
	}
}

func TestInstaller_AlreadyInstalled(t *testing.T) {
	installed := false
	mock := &mockInstaller{
		getPatrolCLIVersionFunc: func() (*v.Version, error) {
			return v.MustParse("3.11.0"), nil
		},
		installPatrolCLIFunc: func() error {
			installed = true
			return nil
		},
	}

	version, err := Run(mock)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if version == nil || version.String() != "3.11.0" {
		t.Fatalf("expected version 3.11.0, got %v", version)
	}
	if installed {
		t.Fatal("install should not be called when CLI is already installed")
	}
}

func TestInstaller_InstallSuccess(t *testing.T) {
	calls := 0
	mock := &mockInstaller{
		getPatrolCLIVersionFunc: func() (*v.Version, error) {
			calls++
			if calls == 1 {
				return nil, errors.New("command not found")
			}
			return v.MustParse("3.11.0"), nil
		},
		installPatrolCLIFunc: func() error {
			return nil
		},
	}

	version, err := Run(mock)
	if err != nil {
		t.Fatalf("expected successful install, got error: %v", err)
	}
	if version == nil || version.String() != "3.11.0" {
		t.Fatalf("expected version 3.11.0, got %v", version)
	}
	if calls != 2 {
		t.Fatalf("expected GetPatrolCLIVersion to be called twice, got %d", calls)
	}
}

func TestInstaller_InstallFailure(t *testing.T) {
	mock := &mockInstaller{
		getPatrolCLIVersionFunc: func() (*v.Version, error) {
			return nil, errors.New("command not found")
		},
		installPatrolCLIFunc: func() error {
			return errors.New("failed to activate package")
		},
	}

	version, err := Run(mock)
	if err == nil {
		t.Fatal("expected error on install failure, got nil")
	}
	if version != nil {
		t.Fatalf("expected version to be nil, got %v", version)
	}
}

func TestInstaller_VerifyAfterInstallFailure(t *testing.T) {
	calls := 0
	mock := &mockInstaller{
		getPatrolCLIVersionFunc: func() (*v.Version, error) {
			calls++
			return nil, errors.New("failed to get version")
		},
		installPatrolCLIFunc: func() error {
			return nil
		},
	}

	version, err := Run(mock)
	if err == nil {
		t.Fatal("expected error on verify failure, got nil")
	}
	if version != nil {
		t.Fatalf("expected version to be nil, got %v", version)
	}
	if calls != 2 {
		t.Fatalf("expected 2 version check calls, got %d", calls)
	}
}
