package subsystem

type Item any

type Value interface {
	From(string) error
	Into() string
}

type BaseSubsystem interface {
	Name() string
	Empty() bool
}

type Subsystem[I Item, V Value] interface {
	BaseSubsystem
	Get() (V, error)
	Set(I) error
	Del(I) error
}

type AnySubsystem Subsystem[Item, Value]

type zeroItem struct{}
type zeroValue struct{}

func (z zeroValue) From(s string) error {
	panic("can not invoke zero value")
}

func (z zeroValue) Into() string {
	panic("can not invoke zero value")
}

type ZeroSubsystem struct{}

func (z ZeroSubsystem) Name() string {
	panic("can not invoke zero subsystem")
}

func (z ZeroSubsystem) Empty() bool {
	panic("can not invoke zero subsystem")
}

func (z ZeroSubsystem) Get() (zeroValue, error) {
	panic("can not invoke zero subsystem")
}

func (z ZeroSubsystem) Set(i zeroItem) error {
	panic("can not invoke zero subsystem")
}

func (z ZeroSubsystem) Del(i zeroItem) error {
	panic("can not invoke zero subsystem")
}
