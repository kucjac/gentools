package testcases

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/kucjac/gentools/parser/testcases/imported"
)

func init() {

}

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

type FooSlice *[3]Foo

type FooAlias Foo

type FooPtrAlias *Foo

type (
	Weird int

	// WeirdStruct docs.
	WeirdStruct struct {
		// Name doc.
		Name string
	}
)

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

type ArrayWrapper [16]byte

type FuncWrapper func(w io.Writer) error

type MultiPointerInlineStruct ******struct {
	// Field test comment.
	Field string
}

func (f FuncWrapper) Do() {}

type (
	SizeCache       = imported.SizeCache
	WeakFields      = imported.WeakFields
	UnknownFields   = imported.UnknownFields
	ExtensionFields = imported.ExtensionFields
)
