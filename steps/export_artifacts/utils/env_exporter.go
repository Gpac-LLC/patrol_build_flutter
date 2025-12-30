package export_artifacts_utils

import "github.com/bitrise-io/go-steputils/tools"

// EnvExporter exports key/value pairs into the environment store.
type EnvExporter interface {
	Export(key, value string) error
}

type envmanExporter struct{}

func (envmanExporter) Export(key, value string) error {
	return tools.ExportEnvironmentWithEnvman(key, value)
}

var envExporter EnvExporter = envmanExporter{}

// SetEnvExporter swaps the exporter used by CopyFilesToFolder. Pass nil to reset to the default.
func SetEnvExporter(exporter EnvExporter) {
	if exporter == nil {
		envExporter = envmanExporter{}
		return
	}
	envExporter = exporter
}

func exportEnv(key, value string) error {
	return envExporter.Export(key, value)
}
