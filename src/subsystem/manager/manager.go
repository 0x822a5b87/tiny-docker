package manager

import (
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cgroup"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/cpu"
	"github.com/0x822a5b87/tiny-docker/src/subsystem/memory"
	"github.com/0x822a5b87/tiny-docker/src/util"
)

type CgroupManager struct {
	fs                 *CgroupFileSystem
	procsSubsystem     *cgroup.ProcsValueSubsystem
	cpuMaxSubsystem    *cpu.MaxValueSubsystem
	memoryMaxSubsystem *memory.MaxValueSubsystem
}

func NewCgroupManager(pid int) (*CgroupManager, error) {
	path := util.GenPidPath(pid)
	fs := NewCgroupFileSystem(path, true)
	err, data := fs.Read(constant.CgroupProcs)
	if err != nil {
		return nil, err
	}
	procsSubsystem, err := cgroup.NewProcsValueSubsystem(data)
	if err != nil {
		return nil, err
	}
	err = procsSubsystem.Set(cgroup.ProcsItem(pid))
	if err != nil {
		return nil, err
	}

	err, data = fs.Read(constant.MemoryMax)
	if err != nil {
		return nil, err
	}
	memoryMaxValueSubsystem, err := memory.NewMaxValueSubsystem(data)
	if err != nil {
		return nil, err
	}

	err, data = fs.Read(constant.CpuMax)
	if err != nil {
		return nil, err
	}
	cpuMaxSubsystem, err := cpu.NewCpuMaxValueSubsystem(data)
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
	err := Write(m.fs, m.procsSubsystem)
	if err != nil {
		return err
	}

	err = Write(m.fs, m.cpuMaxSubsystem)
	if err != nil {
		return err
	}

	err = Write(m.fs, m.memoryMaxSubsystem)
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
