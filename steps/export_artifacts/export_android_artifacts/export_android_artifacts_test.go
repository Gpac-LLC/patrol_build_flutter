package export_android_artifacts

import (
	"os"
	"path/filepath"
	"strings"
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

func TestIsAndroidPlatform(t *testing.T) {
	if IsAndroidPlatform(build_constants.PlatformIOS) {
		t.Error("IsAndroidPlatform should return false for ios")
	}
	if !IsAndroidPlatform(build_constants.PlatformAndroid) {
		t.Error("IsAndroidPlatform should return true for android")
	}
	if !IsAndroidPlatform(build_constants.PlatformBoth) {
		t.Error("IsAndroidPlatform should return true for both")
	}
}

func TestAndroidApkPaths(t *testing.T) {
	testReleasePath, appReleasePath := AndroidApkPaths(true)
	if !strings.Contains(testReleasePath, "release") || !strings.Contains(appReleasePath, "release") {
		t.Error("AndroidApkPaths(true) should return release paths")
	}

	testDebugPath, appDebugPath := AndroidApkPaths(false)
	if !strings.Contains(testDebugPath, "debug") || !strings.Contains(appDebugPath, "debug") {
		t.Error("AndroidApkPaths(false) should return debug paths")
	}
}

func TestFindFirstApkInDir_NotFound(t *testing.T) {
	dir := t.TempDir()
	found, err := FindFirstApkInDir(dir)
	if err != nil {
		t.Fatalf("FindFirstApkInDir error: %v", err)
	}
	if found != "" {
		t.Errorf("expected empty, got %s", found)
	}
}

func TestFindFirstApkInDir(t *testing.T) {
	// Creates a temporary directory
	dir := t.TempDir()
	// Creates a temporary path of a file app-release.apk
	apkPath := filepath.Join(dir, "app-release.apk")
	// Tries to create a temporary file for app-release.apk
	if err := os.WriteFile(apkPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create test apk: %v", err)
	}
	found, err := FindFirstApkInDir(dir)
	if err != nil {
		t.Fatalf("FindFirstApkInDir error: %v", err)
	}
	if found != apkPath {
		t.Errorf("expected %s, got %s", apkPath, found)
	}
}

func TestFindFirstApkInDir_MultipleApks(t *testing.T) {
	dir := t.TempDir()
	apk1 := filepath.Join(dir, "app-release.apk")
	apk2 := filepath.Join(dir, "app-release-androidTest.apk")
	if err := os.WriteFile(apk1, []byte("one"), 0644); err != nil {
		t.Fatalf("failed to create first apk: %v", err)
	}
	if err := os.WriteFile(apk2, []byte("two"), 0644); err != nil {
		t.Fatalf("failed to create second apk: %v", err)
	}

	found, err := FindFirstApkInDir(dir)
	if err != nil {
		t.Fatalf("FindFirstApkInDir error: %v", err)
	}
	if found != apk1 && found != apk2 {
		t.Errorf("expected one of the APKs, got %s", found)
	}
	// The function should return the first one found by Walk, which is not guaranteed to be sorted,
	// but we can check that it is one of the two and document this behavior.
}

func TestCopyAndroidArtifacts_NoAndroid(t *testing.T) {
	setupEnvExporterStub(t)
	t.Setenv(build_constants.Platform, build_constants.PlatformIOS)
	err := CopyAndroidArtifacts(t.TempDir(), t.TempDir(), t.TempDir())
	if err != nil {
		t.Errorf("expected nil for ios platform, got %v", err)
	}
}

func TestCopyAndroidArtifacts_NoApks(t *testing.T) {
	setupEnvExporterStub(t)
	t.Setenv(build_constants.Platform, build_constants.PlatformAndroid)
	t.Setenv("BUILD_TYPE", "debug")
	artifactsPath := t.TempDir()
	testPath := t.TempDir()
	appPath := t.TempDir()
	err := CopyAndroidArtifacts(artifactsPath, testPath, appPath)
	if err != nil {
		t.Errorf("expected nil when no APKs found, got %v", err)
	}
}

func TestCopyAndroidArtifacts_Success(t *testing.T) {
	setupEnvExporterStub(t)
	t.Setenv(build_constants.Platform, build_constants.PlatformAndroid)
	t.Setenv("BUILD_TYPE", "debug")
	testDir := t.TempDir()
	appDir := t.TempDir()
	testApk := filepath.Join(testDir, "app-debug-androidTest.apk")
	appApk := filepath.Join(appDir, "app-debug.apk")
	if err := os.WriteFile(testApk, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test apk: %v", err)
	}
	if err := os.WriteFile(appApk, []byte("app"), 0644); err != nil {
		t.Fatalf("failed to create app apk: %v", err)
	}
	artifactsPath := t.TempDir()
	err := CopyAndroidArtifacts(artifactsPath, testDir, appDir)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	files, _ := os.ReadDir(artifactsPath)
	if len(files) != 2 {
		var names []string
		for _, f := range files {
			names = append(names, f.Name())
		}
		t.Errorf("expected 2 files copied, got %d: %v", len(files), names)
	}
	found := map[string]bool{}
	for _, f := range files {
		found[f.Name()] = true
	}
	if !found["app-debug.apk"] || !found["app-debug-androidTest.apk"] {
		t.Errorf("expected app-debug.apk and app-debug-androidTest.apk, got %v", found)
	}
}

func TestCopyAndroidArtifactsFromEnv_Release(t *testing.T) {
	// GIVEN a release build with APKs
	setupEnvExporterStub(t)
	t.Setenv(build_constants.Platform, build_constants.PlatformAndroid)
	t.Setenv("TEST_BUILD_TYPE", "release")
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
	testDir := filepath.Join(workDir, AndroidTestPath, ReleaseFolder)
	appDir := filepath.Join(workDir, AndroidAppPath, ReleaseFolder)
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	testApk := filepath.Join(testDir, "app-release-androidTest.apk")
	appApk := filepath.Join(appDir, "app-release.apk")
	if err := os.WriteFile(testApk, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test apk: %v", err)
	}
	if err := os.WriteFile(appApk, []byte("app"), 0644); err != nil {
		t.Fatalf("failed to create app apk: %v", err)
	}

	// WHEN exporting using env-derived paths
	err = CopyAndroidArtifactsFromEnv()

	// THEN it succeeds and copies artifacts
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	entries, err := os.ReadDir(AndroidArtifactsPath)
	if err != nil {
		t.Fatalf("read artifacts dir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(entries))
	}
}
