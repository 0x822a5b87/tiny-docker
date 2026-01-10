package daemon

import (
	"os"
	"os/exec"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func StartDockerd(debug bool) error {
	initContext()
	cmd, err := newDaemonProcessCmd()
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
func newDaemonProcessCmd() (*exec.Cmd, error) {
	args := []string{constant.InitDaemon.String()}
	execPath, err := util.GetExecutableAbsolutePath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	cmd.SysProcAttr = &unix.SysProcAttr{}
	cmd.Env = os.Environ()

	if err = configureDaemonProcessTerminalAndDaemonMode(cmd, conf.GlobalConfig.InnerEnv); err != nil {
		return nil, err
	}

	return cmd, nil
}

func configureDaemonProcessTerminalAndDaemonMode(cmd *exec.Cmd, env []string) error {
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Setctty = false

	logFile, err := util.EnsureOpenFilePath(conf.RuntimeDockerdLogFile.Get())
	if err != nil {
		logrus.Fatal("Failed to open log file: ", err)
		return err
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	cmd.Env = env

	return nil
}

func initContext() {
	conf.LoadDaemonConfig()
	ensureFile(conf.RuntimeDockerdUdsFile.Get())
	ensureFile(conf.RuntimeDockerdUdsPidFile.Get())
	ensureFile(conf.RuntimeDockerdLogFile.Get())

	ensurePath(conf.RuntimeDockerdContainerStatus.Get())
	logrus.Infof("init dockerd uds file: {%s}", conf.RuntimeDockerdUdsFile.Get())
	logrus.Infof("init dockerd uds pid file: {%s}", conf.RuntimeDockerdUdsPidFile.Get())
	logrus.Infof("init dockerd log file: {%s}", conf.RuntimeDockerdLogFile.Get())
	logrus.Infof("init dockerd container path: {%s}", conf.RuntimeDockerdContainerStatus.Get())
}

func ensureFile(path string) {
	if err := util.EnsureFileExists(path); err != nil {
		panic(err)
	}
}

func ensurePath(path string) {
	if err := util.EnsureFilePathExist(path); err != nil {
		panic(err)
	}
}
