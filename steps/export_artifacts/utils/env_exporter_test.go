package export_artifacts_utils

import "testing"

type envExporterSpy struct {
	called bool
	key    string
	value  string
}

func (s *envExporterSpy) Export(key, value string) error {
	s.called = true
	s.key = key
	s.value = value
	return nil
}

func TestExportEnv_UsesConfiguredExporter(t *testing.T) {
	spy := &envExporterSpy{}
	SetEnvExporter(spy)
	t.Cleanup(func() {
		SetEnvExporter(nil)
	})

	if err := exportEnv("TEST_KEY", "TEST_VALUE"); err != nil {
		t.Fatalf("exportEnv returned error: %v", err)
	}
	if !spy.called {
		t.Fatal("expected Export to be called on configured exporter")
	}
	if spy.key != "TEST_KEY" || spy.value != "TEST_VALUE" {
		t.Fatalf("expected Export(TEST_KEY, TEST_VALUE), got (%s, %s)", spy.key, spy.value)
	}
}

func TestSetEnvExporter_ResetToDefault(t *testing.T) {
	spy := &envExporterSpy{}
	SetEnvExporter(spy)
	SetEnvExporter(nil)

	if envExporter == nil {
		t.Fatal("expected default exporter after reset, got nil")
	}
	if _, ok := envExporter.(envmanExporter); !ok {
		t.Fatalf("expected default envman exporter after reset, got %T", envExporter)
	}
}
