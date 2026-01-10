package manager

import (
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cgroup"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cpu"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/memory"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

type CgroupManager struct {
	fs                 *CgroupFileSystem
	procsSubsystem     *cgroup.ProcsValueSubsystem
	cpuMaxSubsystem    *cpu.MaxValueSubsystem
	memoryMaxSubsystem *memory.MaxValueSubsystem
}

func NewCgroupManager(pid int) (*CgroupManager, error) {
	path := util.GenContainerCgroupPath(conf.GlobalConfig.Cmd.Id)
	fs := NewCgroupFileSystem(path, true)
	procsSubsystem, err := newSubsystem[*cgroup.ProcsValueSubsystem](fs, constant.CgroupProcs)
	if err != nil {
		return nil, err
	}
	err = procsSubsystem.Set(cgroup.ProcsItem(pid))

	cpuMaxSubsystem, err := newSubsystem[*cpu.MaxValueSubsystem](fs, constant.CpuMax)
	if err != nil {
		return nil, err
	}

	memoryMaxValueSubsystem, err := newSubsystem[*memory.MaxValueSubsystem](fs, constant.MemoryMax)
	if err != nil {
		return nil, err
	}

	return &CgroupManager{
		fs:                 fs,
		procsSubsystem:     procsSubsystem,
		cpuMaxSubsystem:    cpuMaxSubsystem,
		memoryMaxSubsystem: memoryMaxValueSubsystem,
	}, nil
}

func (m *CgroupManager) SetMemoryMax(memoryLimit int) error {
	item := memory.MaxItem{Bytes: int64(memoryLimit)}
	return m.memoryMaxSubsystem.Set(item)
}

func (m *CgroupManager) SetCpuMax(quota, period int) error {
	item := cpu.MaxItem{Quota: quota, Period: period}
	return m.cpuMaxSubsystem.Set(item)
}

func (m *CgroupManager) DelProcsPid(pid cgroup.ProcsItem) error {
	return m.procsSubsystem.Del(pid)
}

func (m *CgroupManager) Sync() error {
	err := Write(m.fs, m.cpuMaxSubsystem)
	if err != nil {
		return err
	}

	err = Write(m.fs, m.memoryMaxSubsystem)
	if err != nil {
		return err
	}

	err = Write(m.fs, m.procsSubsystem)
	if err != nil {
		return err
	}

	return nil
}

func (m *CgroupManager) readCgroupProcs() (*cgroup.ProcsValue, error) {
	err, data := m.fs.Read(constant.CgroupProcs)
	if err != nil {
		return nil, err
	}
	v := &cgroup.ProcsValue{}
	err = v.From(data)
	return v, err
}

func newSubsystem[T subsystem.BaseSubsystem](fs *CgroupFileSystem, name string) (T, error) {
	zeroSubsystem := subsystem.ZeroSubsystem{}
	err, data := fs.Read(name)
	if err != nil {
		return any(zeroSubsystem).(T), err
	}

	var ns T

	switch name {
	case constant.CgroupProcs:
		v, e := cgroup.NewProcsValueSubsystem(data)
		err = e
		ns = any(v).(T)
	case constant.MemoryMax:
		v, e := memory.NewMaxValueSubsystem(data)
		err = e
		ns = any(v).(T)
	case constant.CpuMax:
		v, e := cpu.NewCpuMaxValueSubsystem(data)
		err = e
		ns = any(v).(T)
	default:
		panic(fmt.Errorf("unknown subsystem {%s}", name))
	}

	if err != nil {
		logrus.Errorf("[newSubsystem]: {%s}", err)
		return any(zeroSubsystem).(T), err
	}

	return any(ns).(T), nil
}
