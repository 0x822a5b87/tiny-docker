package util

import (
	"fmt"
	"syscall"

	"github.com/sirupsen/logrus"
)

func KillProcessByPID(pid int, signal int) error {
	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		return fmt.Errorf("syscall kill process %d failed: %w", pid, err)
	}
	logrus.Infof("process %d killed successfully with signal %d", pid, signal)
	return nil
}
