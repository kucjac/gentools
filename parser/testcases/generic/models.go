package generic

// Number exercises union and approximation constraints.
type Number interface {
	~int | ~int64
}

// Pair is a generic named type with nested uses of its parameters.
type Pair[T Number, U any] struct {
	First  T
	Second U
	Nested map[T][]*Pair[T, U]
}

// Transformer exercises a generic interface method.
type Transformer[T any] interface {
	Transform(T) T
}

// Identity is a generic function.
func Identity[T any](value T) T { return value }

// Collect is a variadic generic function.
func Collect[T Number](values ...T) []T { return values }

// Copy exercises generic receiver metadata and instantiated method results.
func (pair Pair[T, U]) Copy() Pair[T, U] { return pair }

// Uses exercises direct and nested instantiations.
type Uses struct {
	Direct Pair[int, string]
	Nested []map[string]Pair[int, string]
}
