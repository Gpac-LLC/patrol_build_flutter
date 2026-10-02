package export_artifacts

import (
	"os"

	build_constants "patrol_build_flutter/steps/build/constants"
	export_android_artifacts "patrol_build_flutter/steps/export_artifacts/export_android_artifacts"
	export_ios_artifacts "patrol_build_flutter/steps/export_artifacts/export_ios_artifacts"
	print "patrol_build_flutter/utils/print"
)

var exportAndroid = func() error {
	return export_android_artifacts.CopyAndroidArtifactsFromEnv()
}

var exportIOS = func() error {
	return export_ios_artifacts.CopyIOSArtifacts(export_ios_artifacts.IOSArtifactsPath)
}

type ExporterRunner struct{}

func (p *ExporterRunner) FindAndExportAndroid() error {
	return exportAndroid()
}

func (p *ExporterRunner) FindAndExportIOS() error {
	return exportIOS()
}

// FindAndExport runs platform-specific exports based on PLATFORM env.
func (p *ExporterRunner) FindAndExport() error {
	switch os.Getenv(build_constants.Platform) {
	case build_constants.PlatformAndroid:
		return p.FindAndExportAndroid()
	case build_constants.PlatformIOS:
		return p.FindAndExportIOS()
	case build_constants.PlatformBoth:
		if err := p.FindAndExportAndroid(); err != nil {
			return err
		}
		return p.FindAndExportIOS()
	default:
		print.Action("No valid platform selected for export")
		return nil
	}
}
