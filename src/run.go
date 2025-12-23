package main

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/container"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cpu"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/manager"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, command string, cfg conf.CgroupConfig) error {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	pid := parent.Process.Pid
	cgroupManager, err := manager.NewCgroupManager(pid)
	if err != nil {
		return err
	}
	log.Println("cgroup path = ", pid)
	err = setConf(cgroupManager, cfg)
	if err != nil {
		return err
	}
	err = cgroupManager.Sync()
	if err != nil {
		return err
	}
	err = parent.Wait()
	if err != nil {
		return err
	}
	os.Exit(-1)
	return nil
}

func setConf(cgroupManager *manager.CgroupManager, cfg conf.CgroupConfig) error {
	memoryLimit, _ := subsystem.SizeToBytes(cfg.MemoryLimit)
	err := cgroupManager.SetMemoryMax(int(memoryLimit))
	if err != nil {
		return err
	}

	v := cpu.MaxValue{}
	err = v.From(cfg.CpuShares)
	if err != nil {
		return err
	}
	err = cgroupManager.SetCpuMax(v.Quota, v.Period)
	if err != nil {
		return err
	}
	return nil
}
