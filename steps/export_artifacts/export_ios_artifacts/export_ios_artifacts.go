package export_ios_artifacts

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	build_constants "patrol_build_flutter/steps/build/constants"
	export_artifacts_utils "patrol_build_flutter/steps/export_artifacts/utils"
	print "patrol_build_flutter/utils/print"
)

var errInvalidBuildFlags = errors.New("invalid iOS build flags")

type zipFilesFunc func(zipPath string, inputPaths []string, executor export_artifacts_utils.CommandExecutor) (string, error)

var zipFiles zipFilesFunc = export_artifacts_utils.ZipFiles

func setZipFiles(fn zipFilesFunc) {
	if fn == nil {
		zipFiles = export_artifacts_utils.ZipFiles
		return
	}
	zipFiles = fn
}

// CopyIOSArtifacts exports iOS build artifacts into the artifacts folder and via envman.
func CopyIOSArtifacts(artifactsPath string) error {
	platform := os.Getenv(build_constants.Platform)
	if platform != build_constants.PlatformIOS && platform != build_constants.PlatformBoth {
		print.Action("No iOS builds were selected to build")
		return nil
	}

	buildType := os.Getenv(build_constants.BuildType)
	buildDirName, err := resolveBuildDirName(buildType)
	if err != nil {
		return err
	}

	buildProductsPath := IOSBuildProductsPath
	buildDir := filepath.Join(buildProductsPath, buildDirName)

	appUnderTest, err := findRequiredApp(buildDir, IOSAppUnderTestName)
	if err != nil {
		return err
	}
	testInstrumentation, err := findRequiredApp(buildDir, IOSTestInstrumentation)
	if err != nil {
		return err
	}

	xctestrunFiles, err := findXCTestRunFiles(buildProductsPath)
	if err != nil {
		return err
	}
	selectedXCTestRun := xctestrunFiles[0]

	if err := export_artifacts_utils.CreateFolder(artifactsPath); err != nil {
		return err
	}

	artifactFiles := []string{appUnderTest, testInstrumentation, selectedXCTestRun}
	artifactKeys := []string{IOSAppUnderTestPathEnvKey, IOSTestInstrumentationEnvKey, IOSRunnerFilePathEnvKey}
	if err := export_artifacts_utils.CopyFilesToFolder(artifactFiles, artifactsPath, artifactKeys); err != nil {
		return err
	}

	zipPath := filepath.Join(buildProductsPath, IOSExportsZipName)
	inputPaths := append([]string{filepath.Join(buildProductsPath, buildDirName)}, xctestrunFiles...)
	zipPath, err = zipFiles(zipPath, inputPaths, nil)
	if err != nil {
		return err
	}

	if err := export_artifacts_utils.CopyFilesToFolder([]string{zipPath}, artifactsPath, []string{IOSBuildExportsZipPathEnvKey}); err != nil {
		return err
	}

	return nil
}

func resolveBuildDirName(buildType string) (string, error) {
	switch buildType {
	case "release":
		if _, err := os.Stat(filepath.Join(IOSBuildProductsPath, "Release-iphonesimulator")); err == nil {
			return "", errInvalidBuildFlags
		}
		return IOSReleaseBuildDirName, nil
	case "debug":
		if _, err := os.Stat(filepath.Join(IOSBuildProductsPath, IOSDebugBuildDirName)); err != nil {
			if os.IsNotExist(err) {
				return "", errInvalidBuildFlags
			}
			return "", err
		}
		return IOSDebugBuildDirName, nil
	default:
		return "", fmt.Errorf("unsupported build type: %s", buildType)
	}
}

func findRequiredApp(buildDir, appName string) (string, error) {
	appPath := filepath.Join(buildDir, appName)
	info, err := os.Stat(appPath)
	if err != nil {
		return "", fmt.Errorf("missing iOS artifact %s: %w", appName, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("expected %s to be a directory", appName)
	}
	return appPath, nil
}

func findXCTestRunFiles(buildProductsPath string) ([]string, error) {
	pattern := filepath.Join(buildProductsPath, IOSXCTestRunGlobPattern)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("missing xctestrun file in %s", buildProductsPath)
	}
	sort.Strings(matches)
	return matches, nil
}
