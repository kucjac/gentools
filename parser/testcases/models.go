package testcases

import (
	"context"
	"strconv"
	"time"
)

type Enumerated int

const (
	_ Enumerated = iota
	// EnumeratedOne defines a first enumerated type value.
	EnumeratedOne
	EnumeratedTwo
)

// FooID is the custom type wrapper on the Foo identifier.
type FooID int64

// UnmarshalText implements encoding.TextUnmarshaler interface.
func (f *FooID) UnmarshalText(in []byte) error {
	i, err := strconv.Atoi(string(in))
	if err != nil {
		return err
	}
	*f = FooID(i)
	return nil
}

// MarshalText implements encoding.TextMarshaler interface.
func (f FooID) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(int(f))), nil
}

// Foo is the test model that contains multiple field definitions.
type Foo struct {
	// ID is the foo field identifier.
	ID         FooID  `json:"id"`
	String     string `custom:"name"`
	CustomName string
	Bool       bool
	Enumerated Enumerated
	Slice      []string
	Float64    float64
	Duration   time.Duration
	Bar        *Bar
}

type NotEmpty interface {
	Call(ctx context.Context, options ...string) (n int, err error)
	InheritMe
}

// InheritMe is an interface that will be inherited.
type InheritMe interface {
	Inherited()
}

type Bar struct {
	Map      map[string]byte
	Time     time.Time
	Any      interface{}
	ChanIn   chan<- int
	ChanOut  <-chan int
	Chan     chan int
	Error    error
	NotEmpty NotEmpty
}