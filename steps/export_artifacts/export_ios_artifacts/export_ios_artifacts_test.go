package export_ios_artifacts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	build_constants "patrol_build_flutter/steps/build/constants"
	export_artifacts_utils "patrol_build_flutter/steps/export_artifacts/utils"
)

type stubEnvExporter struct {
	t        *testing.T
	exported map[string]string
}

func (s *stubEnvExporter) Export(key, value string) error {
	if s.exported == nil {
		s.exported = make(map[string]string)
	}
	s.exported[key] = value
	s.t.Setenv(key, value)
	return nil
}

func setupEnvExporterStub(t *testing.T) *stubEnvExporter {
	stub := &stubEnvExporter{
		t:        t,
		exported: make(map[string]string),
	}
	export_artifacts_utils.SetEnvExporter(stub)
	t.Cleanup(func() {
		export_artifacts_utils.SetEnvExporter(nil)
	})
	return stub
}

type zipRunnerStub struct {
	called     bool
	zipPath    string
	inputPaths []string
	err        error
}

func (s *zipRunnerStub) Run(zipPath string, inputPaths []string, _ export_artifacts_utils.CommandExecutor) (string, error) {
	s.called = true
	s.zipPath = zipPath
	s.inputPaths = append([]string(nil), inputPaths...)
	if s.err != nil {
		return "", s.err
	}
	if err := os.WriteFile(zipPath, []byte("zip"), 0644); err != nil {
		return "", err
	}
	return zipPath, nil
}

func setupZipRunnerStub(t *testing.T, err error) *zipRunnerStub {
	stub := &zipRunnerStub{err: err}
	setZipFiles(stub.Run)
	t.Cleanup(func() {
		setZipFiles(nil)
	})
	return stub
}

func setupWorkingDir(t *testing.T) string {
	workDir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	return workDir
}

func createBuildProducts(t *testing.T, workDir, buildDirName string) (string, string) {
	buildProductsPath := filepath.Join(workDir, IOSBuildProductsPath)
	buildDir := filepath.Join(buildProductsPath, buildDirName)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		t.Fatalf("mkdir build dir: %v", err)
	}
	return buildProductsPath, buildDir
}

func createAppBundle(t *testing.T, buildDir, appName string) string {
	appDir := filepath.Join(buildDir, appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "dummy.txt"), []byte("data"), 0644); err != nil {
		t.Fatalf("write dummy file: %v", err)
	}
	return appDir
}

func createXCTestRun(t *testing.T, buildProductsPath, name string) string {
	path := filepath.Join(buildProductsPath, name)
	if err := os.WriteFile(path, []byte("run"), 0644); err != nil {
		t.Fatalf("write xctestrun: %v", err)
	}
	return path
}

func assertExportedPath(t *testing.T, exported map[string]string, key, expectedPath string) {
	t.Helper()
	got, ok := exported[key]
	if !ok {
		t.Fatalf("expected export key %s", key)
	}
	if got != expectedPath {
		t.Fatalf("expected %s=%s, got %s", key, expectedPath, got)
	}
	if _, err := os.Stat(got); err != nil {
		t.Fatalf("expected exported path to exist: %v", err)
	}
}

func TestCopyIOSArtifacts_ReleaseSuccess(t *testing.T) {
	// GIVEN a release build with required artifacts
	workDir := setupWorkingDir(t)
	buildProductsPath, buildDir := createBuildProducts(t, workDir, IOSReleaseBuildDirName)
	createAppBundle(t, buildDir, IOSAppUnderTestName)
	createAppBundle(t, buildDir, IOSTestInstrumentation)
	xctestrun := createXCTestRun(t, buildProductsPath, "Runner_1.xctestrun")
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "release")
	envStub := setupEnvExporterStub(t)
	zipStub := setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN artifacts and zip are copied and exported
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !zipStub.called {
		t.Fatalf("expected zip runner to be called")
	}
	expectedZipPath := filepath.Join(IOSBuildProductsPath, IOSExportsZipName)
	if zipStub.zipPath != expectedZipPath {
		t.Fatalf("expected zip path %s, got %s", expectedZipPath, zipStub.zipPath)
	}
	expectedInputPaths := []string{
		filepath.Join(IOSBuildProductsPath, IOSReleaseBuildDirName),
		filepath.Join(IOSBuildProductsPath, filepath.Base(xctestrun)),
	}
	if len(zipStub.inputPaths) != len(expectedInputPaths) {
		t.Fatalf("expected zip inputs %v, got %v", expectedInputPaths, zipStub.inputPaths)
	}
	for i, path := range expectedInputPaths {
		if zipStub.inputPaths[i] != path {
			t.Fatalf("expected zip input %s at %d, got %s", path, i, zipStub.inputPaths[i])
		}
	}
	expectedAppPath := filepath.Join(artifactsPath, IOSAppUnderTestName)
	assertExportedPath(t, envStub.exported, IOSAppUnderTestPathEnvKey, expectedAppPath)
	expectedTestAppPath := filepath.Join(artifactsPath, IOSTestInstrumentation)
	assertExportedPath(t, envStub.exported, IOSTestInstrumentationEnvKey, expectedTestAppPath)
	expectedRunnerPath := filepath.Join(artifactsPath, filepath.Base(xctestrun))
	assertExportedPath(t, envStub.exported, IOSRunnerFilePathEnvKey, expectedRunnerPath)
	expectedExportZipPath := filepath.Join(artifactsPath, IOSExportsZipName)
	assertExportedPath(t, envStub.exported, IOSBuildExportsZipPathEnvKey, expectedExportZipPath)
}

func TestCopyIOSArtifacts_DebugSimulatorSuccess(t *testing.T) {
	// GIVEN a debug simulator build with required artifacts
	workDir := setupWorkingDir(t)
	buildProductsPath, buildDir := createBuildProducts(t, workDir, IOSDebugBuildDirName)
	createAppBundle(t, buildDir, IOSAppUnderTestName)
	createAppBundle(t, buildDir, IOSTestInstrumentation)
	createXCTestRun(t, buildProductsPath, "Runner_2.xctestrun")
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "debug")
	envStub := setupEnvExporterStub(t)
	zipStub := setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN artifacts and zip are copied and exported
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !zipStub.called {
		t.Fatalf("expected zip runner to be called")
	}
	expectedAppPath := filepath.Join(artifactsPath, IOSAppUnderTestName)
	assertExportedPath(t, envStub.exported, IOSAppUnderTestPathEnvKey, expectedAppPath)
	expectedTestAppPath := filepath.Join(artifactsPath, IOSTestInstrumentation)
	assertExportedPath(t, envStub.exported, IOSTestInstrumentationEnvKey, expectedTestAppPath)
	expectedExportZipPath := filepath.Join(artifactsPath, IOSExportsZipName)
	assertExportedPath(t, envStub.exported, IOSBuildExportsZipPathEnvKey, expectedExportZipPath)
}

func TestCopyIOSArtifacts_MissingArtifacts(t *testing.T) {
	// GIVEN a build directory missing the RunnerUITests app
	workDir := setupWorkingDir(t)
	buildProductsPath, buildDir := createBuildProducts(t, workDir, IOSReleaseBuildDirName)
	createAppBundle(t, buildDir, IOSAppUnderTestName)
	createXCTestRun(t, buildProductsPath, "Runner_1.xctestrun")
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "release")
	envStub := setupEnvExporterStub(t)
	setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN it fails and does not export paths
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if len(envStub.exported) != 0 {
		t.Fatalf("expected no exports, got %v", envStub.exported)
	}
}

func TestCopyIOSArtifacts_MissingXCTestRun(t *testing.T) {
	// GIVEN a build directory without xctestrun files
	workDir := setupWorkingDir(t)
	_, buildDir := createBuildProducts(t, workDir, IOSReleaseBuildDirName)
	createAppBundle(t, buildDir, IOSAppUnderTestName)
	createAppBundle(t, buildDir, IOSTestInstrumentation)
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "release")
	envStub := setupEnvExporterStub(t)
	setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN it fails and does not export paths
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if len(envStub.exported) != 0 {
		t.Fatalf("expected no exports, got %v", envStub.exported)
	}
}

func TestCopyIOSArtifacts_DebugWithoutSimulator(t *testing.T) {
	// GIVEN a debug device build output
	workDir := setupWorkingDir(t)
	buildProductsPath := filepath.Join(workDir, IOSBuildProductsPath)
	deviceDebugDir := filepath.Join(buildProductsPath, "Debug-iphoneos")
	if err := os.MkdirAll(deviceDebugDir, 0755); err != nil {
		t.Fatalf("mkdir device debug dir: %v", err)
	}
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "debug")
	setupEnvExporterStub(t)
	setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN it fails with invalid combo
	if err == nil || !errors.Is(err, errInvalidBuildFlags) {
		t.Fatalf("expected invalid build flags error, got %v", err)
	}
}

func TestCopyIOSArtifacts_ReleaseWithSimulator(t *testing.T) {
	// GIVEN a release simulator build output
	workDir := setupWorkingDir(t)
	buildProductsPath := filepath.Join(workDir, IOSBuildProductsPath)
	releaseSimDir := filepath.Join(buildProductsPath, "Release-iphonesimulator")
	if err := os.MkdirAll(releaseSimDir, 0755); err != nil {
		t.Fatalf("mkdir release simulator dir: %v", err)
	}
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "release")
	setupEnvExporterStub(t)
	setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN it fails with invalid combo
	if err == nil || !errors.Is(err, errInvalidBuildFlags) {
		t.Fatalf("expected invalid build flags error, got %v", err)
	}
}

func TestCopyIOSArtifacts_ZipFailure(t *testing.T) {
	// GIVEN a valid build but zip runner fails
	workDir := setupWorkingDir(t)
	buildProductsPath, buildDir := createBuildProducts(t, workDir, IOSReleaseBuildDirName)
	createAppBundle(t, buildDir, IOSAppUnderTestName)
	createAppBundle(t, buildDir, IOSTestInstrumentation)
	createXCTestRun(t, buildProductsPath, "Runner_1.xctestrun")
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "release")
	envStub := setupEnvExporterStub(t)
	setupZipRunnerStub(t, fmt.Errorf("zip failed"))

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN it fails and does not export the zip path
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if _, ok := envStub.exported[IOSBuildExportsZipPathEnvKey]; ok {
		t.Fatalf("expected no zip export on failure")
	}
}

func TestCopyIOSArtifacts_SelectsFirstXCTestRun(t *testing.T) {
	// GIVEN multiple xctestrun files
	workDir := setupWorkingDir(t)
	buildProductsPath, buildDir := createBuildProducts(t, workDir, IOSReleaseBuildDirName)
	createAppBundle(t, buildDir, IOSAppUnderTestName)
	createAppBundle(t, buildDir, IOSTestInstrumentation)
	createXCTestRun(t, buildProductsPath, "b.xctestrun")
	first := createXCTestRun(t, buildProductsPath, "a.xctestrun")
	artifactsPath := t.TempDir()
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	t.Setenv(build_constants.BuildType, "release")
	envStub := setupEnvExporterStub(t)
	setupZipRunnerStub(t, nil)

	// WHEN exporting iOS artifacts
	err := CopyIOSArtifacts(artifactsPath)

	// THEN the first sorted xctestrun is exported
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	expectedRunnerPath := filepath.Join(artifactsPath, filepath.Base(first))
	assertExportedPath(t, envStub.exported, IOSRunnerFilePathEnvKey, expectedRunnerPath)
}
