package main

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/container"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
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
	setConf(cgroupManager, cfg)
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

func setConf(cgroupManager *manager.CgroupManager, cfg conf.CgroupConfig) {
	memoryLimit, _ := subsystem.SizeToBytes(cfg.MemoryLimit)
	cgroupManager.SetMemoryMax(int(memoryLimit))
}
