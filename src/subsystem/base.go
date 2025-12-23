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
