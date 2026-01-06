package util

import (
	"fmt"
	"os"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
)

func GenPidPath(pid int) string {
	return fmt.Sprintf("%s/%s", constant.CgroupBasePath, constant.DefaultContainerName)
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
		return fmt.Errorf("mount overlayfs failed: %v (opts: %s)", err, overlayOpts)
	}

	logrus.Infof("mount overlayfs success: lower=%s, upper=%s, work=%s -> merged=%s",
		lowerDir, upperDir, workDir, mergedDir)
	return nil
}

// UnmountOverlayFS 卸载 OverlayFS（清理时用）
func UnmountOverlayFS(mergedDir string) error {
	if err := syscall.Unmount(mergedDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount overlayfs %s failed: %v", mergedDir, err)
	}
	logrus.Infof("unmount overlayfs success: %s", mergedDir)
	return nil
}
