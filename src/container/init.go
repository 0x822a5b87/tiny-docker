package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("init process command: %s, args: %v", command, args)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("mount proc err: %v", err)
		return err
	}
	argv := []string{command}
	if err = syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf("exec error : %s", err.Error())
	}
	return nil
}
