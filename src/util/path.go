package util

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

func GetExecutableAbsolutePath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("os.Executable failed: %v", err)
	}

	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("filepath.EvalSymlinks failed: %v", err)
	}

	return realPath, nil
}

func EnsureOpenFilePath(path string) (*os.File, error) {
	logDir := filepath.Dir(path)
	if err := EnsureFilePathExist(logDir); err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal("Failed to open log file: ", err)
		return nil, err
	}
	return logFile, err
}

func EnsureFileExists(path string) error {
	logDir := filepath.Dir(path)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Errorf("Failed to create path directory, err : %s, path : %s", err, path)
		return err
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, createErr := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
			if createErr != nil {
				logrus.Errorf("Failed to create file %s: %v", path, createErr)
				panic(err)
				return createErr
			}
			if closeErr := file.Close(); closeErr != nil {
				logrus.Warnf("Failed to close created file %s: %v", path, closeErr)
			}
			return nil
		}
		logrus.Errorf("Failed to stat file %s: %v", path, err)
		return err
	}

	return nil
}

func EnsureFilePathExist(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		logrus.Errorf("Failed to create path directory, err : %s, path : %s", err, path)
		return err
	}
	return nil
}

func EnsureDirectoryExists(path string) error {
	logrus.Debugf("Ensuring directory exists: {%s}", path)
	if err := os.MkdirAll(path, 0755); err != nil {
		logrus.Infof("Failed to create directory %s: %v", path, err)
		return err
	}

	logrus.Debugf("Directory %s is ready.", path)
	return nil
}

func MountOverlayFS(lowerDir, upperDir, workDir, mergedDir string) error {
	// 1. 校验所有路径是否存在（提前失败，避免挂载时出错）
	paths := map[string]string{
		"lowerdir": lowerDir,
		"upperdir": upperDir,
		"workdir":  workDir,
		"merged":   mergedDir,
	}
	for name, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("%s path %s does not exist", name, path)
		} else if err != nil {
			return fmt.Errorf("stat %s path %s failed: %v", name, path, err)
		}
	}

	overlayOpts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir)
	if err := syscall.Mount("overlay", mergedDir, "overlay", 0, overlayOpts); err != nil {
		return fmt.Errorf("mount overlay failed: %v (opts: %s)", err, overlayOpts)
	}

	logrus.Infof("mount overlay success: lower=%s, upper=%s, work=%s -> merged=%s",
		lowerDir, upperDir, workDir, mergedDir)
	return nil
}

// UnmountOverlayFS unmount OverlayFS
func UnmountOverlayFS(mergedDir string) error {
	if err := syscall.Unmount(mergedDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount overlay %s failed: %v", mergedDir, err)
	}
	logrus.Infof("unmount overlay success: %s", mergedDir)
	return nil
}
