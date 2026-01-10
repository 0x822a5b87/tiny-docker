package daemon

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cpu"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/manager"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/creack/pty"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

func RunContainerCmd(commands conf.RunCommands) error {
	// NOTE THAT `runContainer` ONLY RUNS IN DAEMON PROCESS.
	conf.LoadRunConfig(commands)
	var err error
	data, err := conf.GlobalConfig.String()
	if err != nil {
		return err
	}
	logrus.Infof("run container init config : %s", string(data))
	if err = setupFs(commands.Image); err != nil {
		logrus.Error(err, "error setup fs.")
		return err
	}

	if err = setupCgroup(os.Getpid(), commands.Args, commands.Cfg); err != nil {
		return err
	}

	parent, err := newContainerCmd()
	if err != nil {
		return err
	}

	if !(commands.Tty && !commands.Detach) {
		if err = parent.Start(); err != nil {
			logrus.Error("error start process: ", err)
			return err
		}
		if err = SendContainerInitRequest(parent.Process.Pid); err != nil {
			logrus.Error("error send init request: ", err)
			return err
		}
	}

	if commands.Detach {
		err = SendWaitRequest(entity.WaitRequest{
			Id:  conf.GlobalConfig.Cmd.Id,
			Pid: parent.Process.Pid,
		})
		if err != nil {
			logrus.Errorf("send wait request error : %s", err.Error())
			return err
		}
	}

	if commands.Tty && !commands.Detach {
		err = parent.Wait()
		if err != nil {
			logrus.Errorf("error wait for container in -it: %v", err)
			return err
		}
		// exit container
		if err = SendStopCurrentRequest(); err != nil {
			logrus.Error("error send stop request: ", err)
			return err
		}
	}

	return nil
}

func RunContainer(command string, args []string) error {
	logrus.Infof("init container command: {%s}, args: {%v}", command, args)
	var err error
	_ = setupDetachMode()
	if err = setupUnionFsFromEnv(); err != nil {
		return err
	}
	logrus.Info("setup layer success.")
	if err = setupMount(); err != nil {
		return err
	}
	logrus.Info("setup mount success.")
	path, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	logrus.Infof("running command {%s} with args {%s}", path, args)
	if err = syscall.Exec(path, args, os.Environ()); err != nil {
		logrus.Errorf("exec error : %s", err.Error())
	}
	return nil
}

func setupCgroup(pid int, commands []string, cfg conf.CgroupConfig) error {
	cgroupManager, err := manager.NewCgroupManager(pid)
	if err != nil {
		return err
	}
	logrus.Printf("cgroup pid = {%d}, command = {%s}", pid, commands)
	err = setConf(cgroupManager, cfg)
	if err != nil {
		return err
	}
	return cgroupManager.Sync()
}

// init union fs for layer
func setupFs(image string) error {
	readPath := conf.GlobalConfig.ReadPath()
	writePath := conf.GlobalConfig.WritePath()
	workPath := conf.GlobalConfig.WorkPath()
	mergePath := conf.GlobalConfig.MergePath()

	var err error
	if err = util.EnsureDirectoryExists(readPath); err != nil {
		logrus.Errorf("ensure directory error : %s", err.Error())
		return err
	}
	if err = util.EnsureDirectoryExists(writePath); err != nil {
		logrus.Errorf("ensure directory error : %s", err.Error())
		return err
	}
	if err = util.EnsureDirectoryExists(workPath); err != nil {
		logrus.Errorf("ensure directory error : %s", err.Error())
		return err
	}
	if err = util.EnsureDirectoryExists(mergePath); err != nil {
		logrus.Errorf("ensure directory error : %s", err.Error())
		return err
	}

	if err = util.UnTar(image, readPath); err != nil {
		logrus.Errorf("untar image error : %s", err.Error())
		return err
	}

	return nil
}

func setConf(cgroupManager *manager.CgroupManager, cfg conf.CgroupConfig) error {
	err := setMemoryLimit(cgroupManager, cfg)
	if err != nil {
		return err
	}

	err = setCpuShares(cgroupManager, cfg)
	if err != nil {
		return err
	}

	return nil
}

func setMemoryLimit(cgroupManager *manager.CgroupManager, cfg conf.CgroupConfig) error {
	memoryLimit, _ := subsystem.SizeToBytes(cfg.MemoryLimit)
	return cgroupManager.SetMemoryMax(int(memoryLimit))
}

func setCpuShares(cgroupManager *manager.CgroupManager, cfg conf.CgroupConfig) error {
	v := cpu.MaxValue{}
	err := v.From(cfg.CpuShares)
	if err != nil {
		logrus.Errorf("set cpu shares error : %s", err.Error())
		return err
	}
	return cgroupManager.SetCpuMax(v.Quota, v.Period)
}

func newContainerCmd() (*exec.Cmd, error) {
	args := []string{constant.InitContainer.String()}
	commands := conf.GlobalConfig.Cmd
	for _, arg := range commands.Args {
		args = append(args, arg)
	}
	cmd := exec.Command(constant.UnixProcSelfExe, args...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS |
			unix.CLONE_NEWPID |
			unix.CLONE_PIDFD |
			unix.CLONE_NEWNS |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWIPC,
		Unshareflags: unix.CLONE_NEWNS,
	}

	cmd.Dir = conf.GlobalConfig.MergePath()
	cmd.Env = commands.UserEnv
	cmd.Env = append(cmd.Env, conf.GlobalConfig.InnerEnv...)

	if err := configureContainerProcessTerminalAndDaemonMode(cmd, commands.Tty, commands.Detach); err != nil {
		return nil, err
	}

	return cmd, nil
}

func configureContainerProcessTerminalAndDaemonMode(cmd *exec.Cmd, interactive bool, detach bool) error {
	if interactive && detach {
		return constant.ErrProcessTerminalAndDaemonMode
	}
	logrus.Infof("Running new process in interactive mode {%v}, detach mode {%v}.", interactive, detach)
	// create a new session in any mode
	if detach {
		// -d: terminal and running as daemon
		cmd.SysProcAttr.Setctty = false
		cmd.Stdin = nil
		logFile, err := util.EnsureOpenFilePath(getContainerLogFilePath(conf.GlobalConfig.Cmd.Id))
		if err != nil {
			logrus.Errorf("ensure log file error : %s", err.Error())
			return err
		}
		cmd.Stdout = logFile
		cmd.Stderr = logFile
	} else if interactive {
		err := setPty(cmd)
		if err != nil {
			logrus.Errorf("set pty error : %s", err.Error())
			return err
		}
	} else {
		cmd.SysProcAttr.Setctty = false
		cmd.Stdin = nil
	}
	return nil
}

func setPty(cmd *exec.Cmd) error {
	// -it 模式核心配置
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Setctty = true
	cmd.SysProcAttr.Ctty = 0
	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		logrus.Errorf("error init ptmx: %v", err)
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	if err = SendContainerInitRequest(cmd.Process.Pid); err != nil {
		logrus.Errorf("send init request error : %s", err.Error())
		return err
	}

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err = pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func getContainerLogFilePath(id string) string {
	fileRoot := conf.RuntimeDockerdContainerLog.Get()
	return filepath.Join(fileRoot, id, constant.ContainerLogFile)
}

func nullFileWithPanic() *os.File {
	nullFile, err := util.NullFile()
	if err != nil {
		logrus.Errorf("error open null file : %s", err.Error())
		panic(err)
	}
	return nullFile
}
