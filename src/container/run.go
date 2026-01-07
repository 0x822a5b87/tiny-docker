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
	logrus.Infof("init config : %v", conf.GlobalConfig)
	conf.LoadRunConfig(commands)
	var err error
	if err = setupFs(commands.Image); err != nil {
		logrus.Error(err, "error setup fs.")
		return err
	}
	parent := newParentProcess(commands.Tty, commands.Args, commands.UserEnv)
	setupEnv(parent)
	if err = setupCgroup(commands.Args, commands.Cfg); err != nil {
		return err
	}
	if err = parent.Start(); err != nil {
		logrus.Error(err, "error start process.")
		return err
	}
	if err = parent.Wait(); err != nil {
		logrus.Error("error wait process:", err)
		return err
	}
	os.Exit(-1)
	return nil
}

func setupEnv(cmd *exec.Cmd) {
	util.AppendEnv(cmd, constant.MetaName, conf.GlobalConfig.Meta.Name)
	util.AppendEnv(cmd, constant.FsBasePath, conf.GlobalConfig.Fs.Root)

	util.AppendEnv(cmd, constant.FsReadLayerPath, conf.GlobalConfig.ReadPath())
	util.AppendEnv(cmd, constant.FsWriteLayerPath, conf.GlobalConfig.WritePath())
	util.AppendEnv(cmd, constant.FsWorkLayerPath, conf.GlobalConfig.WorkPath())
	util.AppendEnv(cmd, constant.FsMergeLayerPath, conf.GlobalConfig.MergePath())
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

func newParentProcess(tty bool, commands []string, env []string) *exec.Cmd {
	args := []string{"init"}
	for _, command := range commands {
		args = append(args, command)
	}
	cmd := exec.Command(constant.UnixProcSelfExe, args...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWNS |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWIPC,
		Unshareflags: unix.CLONE_NEWNS,
		Setctty:      tty,
		Setsid:       tty,
	}

	if tty {
		logrus.Info("Running new process in tty.")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Dir = conf.GlobalConfig.MergePath()
	cmd.Env = env
	return cmd
}
