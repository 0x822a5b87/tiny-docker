package cgroup

import (
	"strconv"
	"strings"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
)

type ProcsItem int

type ProcsValue struct {
	pids []ProcsItem
}

func (p *ProcsValue) From(s string) error {
	pidsStrArr := strings.Split(s, "\n")
	p.pids = make([]ProcsItem, len(pidsStrArr))
	for i, v := range pidsStrArr {
		pid, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		p.pids[i] = ProcsItem(pid)
	}
	return nil
}

func (p *ProcsValue) Into() string {
	pidStrArr := make([]string, len(p.pids))
	for i, pid := range p.pids {
		pidStrArr[i] = strconv.Itoa(int(pid))
	}
	return strings.Join(pidStrArr, "\n")
}

func NewProcsValueSubsystem(data string) (*ProcsValueSubsystem, error) {
	v := &ProcsValue{}
	if data != "" {
		err := v.From(data)
		if err != nil {
			return nil, err
		}
	}
	return &ProcsValueSubsystem{
		value: v,
	}, nil
}

type ProcsValueSubsystem struct {
	value *ProcsValue
}

func (p *ProcsValueSubsystem) Name() string {
	return constant.CgroupProcs
}

func (p *ProcsValueSubsystem) Get() (*ProcsValue, error) {
	return p.value, nil
}

func (p *ProcsValueSubsystem) Set(item ProcsItem) error {
	for pid := range p.value.pids {
		if int(item) == int(pid) {
			return nil
		}
	}
	p.value.pids = append(p.value.pids, item)
	return nil
}

func (p *ProcsValueSubsystem) Del(item ProcsItem) error {
	if p == nil || len(p.value.pids) == 0 {
		return constant.ErrProcsEmpty
	}

	newPids := make([]ProcsItem, 0, len(p.value.pids))
	found := false

	for _, pid := range p.value.pids {
		if pid == item {
			found = true
			continue
		}
		newPids = append(newPids, pid)
	}

	if !found {
		return constant.ErrProcsPidNotFound
	}

	p.value.pids = newPids
	return nil
}

func (p *ProcsValueSubsystem) Empty() bool {
	return len(p.value.pids) == 0
}

// add compiler check
var _ subsystem.Value = (*ProcsValue)(nil)
var _ subsystem.Subsystem[ProcsItem, *ProcsValue] = (*ProcsValueSubsystem)(nil)
