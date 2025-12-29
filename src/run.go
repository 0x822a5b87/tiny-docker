package main

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/container"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cpu"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/manager"
	log "github.com/sirupsen/logrus"
)

func Run(commands RunCommands) error {
	parent := container.NewParentProcess(commands.Tty, commands.Commands, commands.UserEnv)
	err := initCgroup(parent, commands.Commands, commands.Cfg)
	if err != nil {
		return err
	}
	if err = parent.Start(); err != nil {
		log.Error(err)
	}
	err = parent.Wait()
	if err != nil {
		return err
	}
	os.Exit(-1)
	return nil
}

func initCgroup(parent *exec.Cmd, commands []string, cfg conf.CgroupConfig) error {
	pid := syscall.Getpid()
	cgroupManager, err := manager.NewCgroupManager(pid)
	if err != nil {
		return err
	}
	log.Printf("cgroup pid = {%d}, command = {%s}", pid, commands)
	err = setConf(cgroupManager, cfg)
	if err != nil {
		return err
	}
	return cgroupManager.Sync()
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
