package memory

import (
	"strconv"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/subsystem"
)

type MaxItem MaxValue

type MaxValue struct {
	Bytes int64
}

func (m *MaxValue) From(s string) error {
	bytes, err := subsystem.SizeToBytes(s)
	if err != nil {
		return err
	}
	m.Bytes = bytes
	return nil
}

func (m *MaxValue) Into() string {
	if m.Bytes == 0 {
		return constant.LiteralMax
	}
	return strconv.Itoa(int(m.Bytes))
}

func NewMaxValueSubsystem(data string) (*MaxValueSubsystem, error) {
	v := &MaxValue{}
	if err := v.From(data); err != nil {
		return nil, err
	}
	return &MaxValueSubsystem{
		value: v,
	}, nil
}

type MaxValueSubsystem struct {
	value *MaxValue
}

func (m *MaxValueSubsystem) Name() string {
	return constant.MemoryMax
}

func (m *MaxValueSubsystem) Get() (*MaxValue, error) {
	return m.value, nil
}

func (m *MaxValueSubsystem) Set(max MaxItem) error {
	m.value.Bytes = max.Bytes
	return nil
}

func (m *MaxValueSubsystem) Del(max MaxItem) error {
	m.value.Bytes = 0
	return nil
}

func (m *MaxValueSubsystem) Empty() bool {
	return m.value == nil || m.value.Bytes == 0
}

// add compiler check
var _ subsystem.Value = (*MaxValue)(nil)
var _ subsystem.Subsystem[MaxItem, *MaxValue] = (*MaxValueSubsystem)(nil)
