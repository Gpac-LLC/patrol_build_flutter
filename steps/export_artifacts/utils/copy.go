package export_artifacts_utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	print "patrol_install/utils/print"
)

// closeWithLog closes a file and logs an error if closing fails.
func closeWithLog(f io.Closer, name string) {
	if err := f.Close(); err != nil {
		print.Error(fmt.Sprintf("Error closing %s: %v", name, err))
	}
}

// CopyFilesToFolder copies each file in srcFiles to destFolder, sets env variable for each file using envKeys.
// Returns an error if any copy fails or if lengths do not match.
func CopyFilesToFolder(srcFiles []string, destFolder string, envKeys []string) error {
	if len(srcFiles) != len(envKeys) {
		return fmt.Errorf("number of files (%d) does not match number of env keys (%d)", len(srcFiles), len(envKeys))
	}
	for i, srcFile := range srcFiles {
		dst := filepath.Join(destFolder, filepath.Base(srcFile))
		src, err := os.Open(srcFile)
		if err != nil {
			print.Error(fmt.Sprintf("Error opening %s: %v", srcFile, err))
			return err
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			print.Error(fmt.Sprintf("Error creating %s: %v", dst, err))
			closeWithLog(src, srcFile)
			return err
		}

		if _, err := io.Copy(dstFile, src); err != nil {
			print.Error(fmt.Sprintf("Error copying %s to %s: %v", srcFile, dst, err))
			closeWithLog(src, srcFile)
			closeWithLog(dstFile, dst)
			return err
		}
		print.Success(fmt.Sprintf("Copied to %s", dst))
		// Close files explicitly to avoid too many open files in a loop
		closeWithLog(src, srcFile)
		closeWithLog(dstFile, dst)

		if err := exportEnv(envKeys[i], dst); err != nil {
			print.Error(fmt.Sprintf("Error exporting env by Envman %s: %v", envKeys[i], err))
			return err
		}
		print.Success(fmt.Sprintf("Artifact: %s exported into: %s \n", dst, envKeys[i]))

	}
	return nil
}
