package daemon

import (
	"os/exec"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func StartDockerd(debug bool) error {
	cmd, err := newDaemonProcessCmd(debug)
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		logrus.Fatal("Failed to start mini dockerd: ", err)
		return err
	}
	logrus.Info("Starting docker daemon on pid: ", cmd.Process.Pid)
	return nil
}

// create mini-dockerd, always running on notty and detach mode.
func newDaemonProcessCmd(debug bool) (*exec.Cmd, error) {
	args := []string{constant.InitDaemon.String()}
	execPath, err := util.GetExecutableAbsolutePath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	cmd.SysProcAttr = &unix.SysProcAttr{}
	if err = configureDaemonProcessTerminalAndDaemonMode(cmd, debug); err != nil {
		return nil, err
	}

	return cmd, nil
}
