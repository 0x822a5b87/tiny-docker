package daemon

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/handler"
	"github.com/sirupsen/logrus"
)

func RunDaemon() error {
	logrus.Info("Starting daemon process")
	_ = setupDetachMode()
	return handler.CreateUdsServer()
}

func setupDetachMode() error {
	detach := conf.DetachMode.GetBoolean()
	if !detach {
		return nil
	}
	sid, err := syscall.Setsid()
	logrus.Infof("Running user process {%d} in detach mode.", sid)
	return err
}

func pivotRoot(root string) error {
	if err := syscall.Mount("", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return constant.ErrMountRootFS.Wrap(err)
	}

	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return constant.ErrMountRootFS.Wrap(err)
	}

	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return constant.ErrMountRootFS.Wrap(err)
	}

	// 在使用 pivot_root 切换根目录的时候，需要两个目录：
	//
	// 1. new_root，这个是我们切换的目标目录；
	// 2. old_root，这个并不是指的我们的宿主机的根目录，而是一个目录用来mount根目录的。
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return constant.ErrMountRootFS.Wrap(err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return constant.ErrMountRootFS.Wrap(err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return constant.ErrMountRootFS.Wrap(err)
	}
	return os.Remove(pivotDir)
}

func setupMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Get current location error %v", err)
		return err
	}
	logrus.Infof("Current location is %s", pwd)
	err = pivotRoot(pwd)
	if err != nil {
		return err
	}

	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("mount proc error : %s", err.Error())
		return err
	}

	return syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}
