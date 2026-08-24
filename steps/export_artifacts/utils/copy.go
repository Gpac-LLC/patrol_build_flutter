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
		print.Errorf("Error closing %s: %v", name, err)
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
		info, err := os.Lstat(srcFile)
		if err != nil {
			print.Errorf("Error opening %s: %v", srcFile, err)
			return err
		}

		if info.IsDir() {
			if err := copyDir(srcFile, dst); err != nil {
				print.Errorf("Error copying directory %s to %s: %v", srcFile, dst, err)
				return err
			}
		} else {
			if err := copyFile(srcFile, dst, info.Mode()); err != nil {
				print.Errorf("Error copying %s to %s: %v", srcFile, dst, err)
				return err
			}
		}

		print.Successf("Copied to %s", dst)

		if err := exportEnv(envKeys[i], dst); err != nil {
			print.Errorf("Error exporting env by Envman %s: %v", envKeys[i], err)
			return err
		}
		print.Successf("Artifact: %s exported into: %s \n", dst, envKeys[i])

	}
	return nil
}

func copyFile(srcPath, dstPath string, mode os.FileMode) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		closeWithLog(src, srcPath)
		return err
	}

	if _, err := io.Copy(dstFile, src); err != nil {
		closeWithLog(src, srcPath)
		closeWithLog(dstFile, dstPath)
		return err
	}

	closeWithLog(src, srcPath)
	closeWithLog(dstFile, dstPath)

	if err := os.Chmod(dstPath, mode.Perm()); err != nil {
		return err
	}

	return nil
}

func copyDir(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return os.MkdirAll(dstPath, info.Mode().Perm())
		}

		if d.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(target, dstPath)
		}

		return copyFile(path, dstPath, info.Mode())
	})
}
