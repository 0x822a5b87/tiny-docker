package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cpu"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/manager"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

func Run(commands conf.RunCommands) error {
	conf.LoadRunConfig(commands)
	logrus.Infof("init config : %v", conf.GlobalConfig)
	var err error
	if err = setupFs(commands.Image); err != nil {
		logrus.Error(err, "error setup fs.")
		return err
	}
	parent := newParentProcess()
	setupEnv(parent)
	if err = setupCgroup(commands.Args, commands.Cfg); err != nil {
		return err
	}
	if err = parent.Start(); err != nil {
		logrus.Error("error start process: ", err)
		return err
	}

	if commands.Tty {
		logrus.Infof("Running {%s} in attach mode.", conf.GlobalConfig.ImageName())
		return parent.Wait()
	}

	logrus.Infof("Running {%s}, pid = {%d} exit.", conf.GlobalConfig.ImageName(), os.Getpid())
	return nil
}

func setupEnv(cmd *exec.Cmd) {
	util.AppendEnv(cmd, constant.MetaName, conf.GlobalConfig.Meta.Name)
	util.AppendEnv(cmd, constant.FsBasePath, conf.GlobalConfig.Fs.Root)

	util.AppendEnv(cmd, constant.FsReadLayerPath, conf.GlobalConfig.ReadPath())
	util.AppendEnv(cmd, constant.FsWriteLayerPath, conf.GlobalConfig.WritePath())
	util.AppendEnv(cmd, constant.FsWorkLayerPath, conf.GlobalConfig.WorkPath())
	util.AppendEnv(cmd, constant.FsMergeLayerPath, conf.GlobalConfig.MergePath())

	if conf.GlobalConfig.Cmd.Detach {
		util.AppendEnv(cmd, constant.DetachMode, "true")
	} else {
		util.AppendEnv(cmd, constant.DetachMode, "false")
	}
}

func setupCgroup(commands []string, cfg conf.CgroupConfig) error {
	pid := syscall.Getpid()
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
		return err
	}
	return cgroupManager.SetCpuMax(v.Quota, v.Period)
}

func newParentProcess() *exec.Cmd {
	commands := conf.GlobalConfig.Cmd
	args := []string{"init"}
	for _, arg := range commands.Args {
		args = append(args, arg)
	}
	cmd := exec.Command(constant.UnixProcSelfExe, args...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWNS |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWIPC,
		Unshareflags: unix.CLONE_NEWNS,
	}

	setTtyMode(cmd, commands.Tty)
	setDetachMode(cmd, commands.Detach, commands.Tty)

	cmd.Dir = conf.GlobalConfig.MergePath()
	cmd.Env = commands.UserEnv
	return cmd
}

func setTtyMode(cmd *exec.Cmd, tty bool) {
	cmd.SysProcAttr.Setctty = tty
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		nullFile, err := os.OpenFile("/dev/null", os.O_RDWR, 0666)
		if err != nil {
			logrus.Warnf("open /dev/null failed: %v", err)
		} else {
			cmd.Stdin = nullFile
			cmd.Stdout = nullFile
			cmd.Stderr = nullFile
		}
	}
}

func setDetachMode(cmd *exec.Cmd, detach bool, tty bool) {
	if detach {
		logrus.Infof("Running new process in detach mode.")
		cmd.SysProcAttr.Setsid = true
		cmd.SysProcAttr.Setctty = false
	} else {
		logrus.Infof("Running new process in attach mode.")
		cmd.SysProcAttr.Setsid = true
		cmd.SysProcAttr.Noctty = !tty
	}

	{
		// TODO delete me
		cmd.Stdout = os.Stdout
	}
}

func nullFileWithPanic() *os.File {
	nullFile, err := util.NullFile()
	if err != nil {
		logrus.Errorf("error open null file : %s", err.Error())
		panic(err)
	}
	return nullFile
}
