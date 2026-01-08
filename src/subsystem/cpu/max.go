package cpu

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
)

type MaxItem MaxValue

type MaxValue struct {
	Quota  int
	Period int
}

func (m *MaxValue) From(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		m.noLimit()
		return nil
	}
	items := strings.Split(s, " ")
	if len(items) != 2 {
		return constant.ErrMalformedType
	}
	var err error = nil
	if items[0] == constant.LiteralMax {
		m.Quota = math.MaxInt
	} else {
		m.Quota, err = strconv.Atoi(items[0])
	}
	if err != nil {
		return err
	}
	m.Period, err = strconv.Atoi(items[1])
	if err != nil {
		return err
	}

	return nil
}

func (m *MaxValue) Into() string {
	var quota = ""
	var period = ""
	if m.Quota == math.MaxInt {
		quota = constant.LiteralMax
	} else {
		quota = strconv.Itoa(m.Quota)
	}
	period = strconv.Itoa(m.Period)
	return fmt.Sprintf("%s %s", quota, period)
}

func (m *MaxValue) noLimit() {
	m.Quota = math.MaxInt
	m.Period = constant.CpuPeriod
}

func NewCpuMaxValueSubsystem(data string) (*MaxValueSubsystem, error) {
	v := &MaxValue{
		Quota:  0,
		Period: 0,
	}
	err := v.From(data)
	if err != nil {
		return nil, err
	}
	return &MaxValueSubsystem{value: v}, nil
}

type MaxValueSubsystem struct {
	value *MaxValue
}

func (m *MaxValueSubsystem) Name() string {
	return constant.CpuMax
}

func (m *MaxValueSubsystem) Get() (*MaxValue, error) {
	return m.value, nil
}

func (m *MaxValueSubsystem) Set(max MaxItem) error {
	m.value.Quota = max.Quota
	m.value.Period = max.Period
	return nil
}

func (m *MaxValueSubsystem) Del(max MaxItem) error {
	m.value.Quota = math.MaxInt
	m.value.Period = 0
	return nil
}

func (m *MaxValueSubsystem) Empty() bool {
	return (m.value.Quota == math.MaxInt && m.value.Period == constant.CpuPeriod) ||
		(m.value.Quota == 0 && m.value.Period == 0)
}

// add compiler check
var _ subsystem.Value = (*MaxValue)(nil)
var _ subsystem.Subsystem[MaxItem, *MaxValue] = (*MaxValueSubsystem)(nil)
