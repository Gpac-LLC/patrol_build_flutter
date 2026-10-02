package export_artifacts

import (
	"errors"
	"testing"

	build_constants "patrol_build_flutter/steps/build/constants"
)

type exportCallState struct {
	androidCalled bool
	iosCalled     bool
}

func stubExports(t *testing.T, androidErr, iosErr error) *exportCallState {
	state := &exportCallState{}
	originalAndroid := exportAndroid
	originalIOS := exportIOS

	exportAndroid = func() error {
		state.androidCalled = true
		return androidErr
	}
	exportIOS = func() error {
		state.iosCalled = true
		return iosErr
	}

	t.Cleanup(func() {
		exportAndroid = originalAndroid
		exportIOS = originalIOS
	})
	return state
}

func TestFindAndExport_AndroidOnly(t *testing.T) {
	// GIVEN Android selected
	t.Setenv(build_constants.Platform, build_constants.PlatformAndroid)
	state := stubExports(t, nil, nil)
	runner := &ExporterRunner{}

	// WHEN running exports
	err := runner.FindAndExport()

	// THEN only Android export runs
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !state.androidCalled || state.iosCalled {
		t.Fatalf("expected only android export, got android=%v ios=%v", state.androidCalled, state.iosCalled)
	}
}

func TestFindAndExport_IOSOnly(t *testing.T) {
	// GIVEN iOS selected
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	state := stubExports(t, nil, nil)
	runner := &ExporterRunner{}

	// WHEN running exports
	err := runner.FindAndExport()

	// THEN only iOS export runs
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !state.iosCalled || state.androidCalled {
		t.Fatalf("expected only ios export, got android=%v ios=%v", state.androidCalled, state.iosCalled)
	}
}

func TestFindAndExport_BothPlatforms(t *testing.T) {
	// GIVEN both selected
	t.Setenv(build_constants.Platform, build_constants.PlatformBoth)
	state := stubExports(t, nil, nil)
	runner := &ExporterRunner{}

	// WHEN running exports
	err := runner.FindAndExport()

	// THEN both exports run
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !state.iosCalled || !state.androidCalled {
		t.Fatalf("expected both exports, got android=%v ios=%v", state.androidCalled, state.iosCalled)
	}
}

func TestFindAndExport_AndroidError(t *testing.T) {
	// GIVEN Android export fails
	t.Setenv(build_constants.Platform, build_constants.PlatformAndroid)
	state := stubExports(t, errors.New("android failed"), nil)
	runner := &ExporterRunner{}

	// WHEN running exports
	err := runner.FindAndExport()

	// THEN error is returned and iOS is not called
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if state.iosCalled {
		t.Fatalf("expected ios export not to run after android failure")
	}
}

func TestFindAndExport_IOSError(t *testing.T) {
	// GIVEN iOS export fails
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	state := stubExports(t, nil, errors.New("ios failed"))
	runner := &ExporterRunner{}

	// WHEN running exports
	err := runner.FindAndExport()

	// THEN error is returned
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !state.iosCalled {
		t.Fatalf("expected ios export to run")
	}
}
