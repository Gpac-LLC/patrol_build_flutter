package export_artifacts_utils

import (
	"os"
	"path/filepath"
	"testing"
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
	SetEnvExporter(stub)
	t.Cleanup(func() {
		SetEnvExporter(nil)
	})
	return stub
}

func TestCopyFilesToFolder(t *testing.T) {
	stub := setupEnvExporterStub(t)
	srcDir := t.TempDir()
	dstDir := t.TempDir()
	file1 := filepath.Join(srcDir, "a.txt")
	file2 := filepath.Join(srcDir, "b.txt")
	if err := os.WriteFile(file1, []byte("foo"), 0644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("bar"), 0644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}
	files := []string{file1, file2}
	envKeys := []string{"TEST_ENV_A", "TEST_ENV_B"}
	if err := CopyFilesToFolder(files, dstDir, envKeys); err != nil {
		t.Fatalf("CopyFilesToFolder failed: %v", err)
	}
	for i, f := range files {
		base := filepath.Base(f)
		dstPath := filepath.Join(dstDir, base)
		if _, err := os.Stat(dstPath); err != nil {
			t.Errorf("file %s not copied", base)
		}
		// Check env variable is set
		val := os.Getenv(envKeys[i])
		if val != dstPath {
			t.Errorf("env %s expected %s, got %s", envKeys[i], dstPath, val)
		}
		if exported, ok := stub.exported[envKeys[i]]; !ok || exported != dstPath {
			t.Errorf("exporter expected %s=%s, got %s", envKeys[i], dstPath, exported)
		}
	}
}

func TestCopyFilesToFolder_Error(t *testing.T) {
	dstDir := t.TempDir()
	files := []string{"/nonexistent/file.txt"}
	envKeys := []string{"DUMMY_ENV"}
	if err := CopyFilesToFolder(files, dstDir, envKeys); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCopyFilesToFolder_LengthMismatch(t *testing.T) {
	dstDir := t.TempDir()
	files := []string{"/tmp/a.txt", "/tmp/b.txt"}
	envKeys := []string{"ENV_A"} // Mismatch
	if err := CopyFilesToFolder(files, dstDir, envKeys); err == nil {
		t.Error("expected error for length mismatch")
	}
}
