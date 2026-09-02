# Research: Generic Parser Compatibility

## Go language and toolchain boundary

**Decision**: Support Go 1.27.1 only. It is the repository-selected current
toolchain and is the sole local and CI validation target.

**Rationale**: The repository has a single consumer, so supporting one current
toolchain avoids the resource cost and complexity of an historical matrix while
retaining generic aliases and methods on generic receiver types.

**Alternatives considered**:

- Retain historical toolchains. Rejected because they are not needed by the
  repository owner and repeated local builds exhaust available memory.

References: [Go 1.27 release notes](https://go.dev/doc/go1.27) and the [Go language specification](https://go.dev/ref/spec).

## Source-of-truth type information

**Decision**: Make `go/types` the single semantic source for generic metadata;
continue using the AST only to associate declarations, comments, and
instantiation sites with that information. Request syntax and type information
when loading packages whose declaration metadata needs to be mapped.

**Rationale**: `go/types` exposes `TypeParam`, `Named.TypeParams`,
`Named.TypeArgs`, `Signature.TypeParams`, and `Signature.RecvTypeParams`.
On current releases aliases are materialized and expose equivalent type
parameter/type-argument accessors. Centralizing translation avoids brittle
string parsing and matches the compiler's constraint and substitution rules.

**Alternatives considered**:

- Parse bracket syntax directly from source text. Rejected because it cannot
  reliably distinguish type arguments, array syntax, inferred arguments, or
  resolved cross-package names.
- Export raw `go/types` objects. Rejected because it would couple Gentools'
  public API to toolchain representation and break its own abstraction.

Reference: [go/types API documentation](https://pkg.go.dev/go/types).

## Public model compatibility

**Decision**: Introduce additive generic metadata types and accessors, plus an
instantiated-type representation for uses that have concrete type arguments.
Keep existing `Struct`, `Interface`, `Alias`, and `Function` representations
and their fields unchanged for established inputs.

**Rationale**: Adding exported fields to existing public structs can break
downstream unkeyed composite literals. New accessors and new types make generic
information queryable without changing existing construction, equality, name,
or string contracts for pre-generic callers.

**Alternatives considered**:

- Add type-parameter fields to every existing exported model struct. Rejected
  because it is source-incompatible for unkeyed literals.
- Replace generic named declarations with a wrapper type. Rejected because it
  would make the parser return an unexpected concrete type for callers that
  expect `*types.Struct` or `*types.Interface`.

## Go-version adapter boundary

**Decision**: Use the Go 1.27.1 `go/types` alias and generic APIs directly.

**Rationale**: No supported compiler lacks these APIs.

**Alternatives considered**:

- Keep version adapters. Rejected because they no longer provide supported
  behavior and obscure the single-toolchain conversion path.

## Dependency policy

**Decision**: Use the dependency versions selected by the Go 1.27.1 module
graph: `golang.org/x/tools v0.49.0` and `golang.org/x/mod v0.40.0`.

**Rationale**: Current tooling fixes the historical package-loader failures and
is appropriate now that older Go builds are outside scope.

**Alternatives considered**:

- Retain an old dependency set. Rejected because it is no longer constrained by
  historical compiler support.

## CI design

**Decision**: Use Go 1.27.1 across the existing Linux, macOS, and Windows CI
matrix, with a separate Linux lint job.

**Rationale**: Separating lint keeps failures attributable while one toolchain
keeps validation resource use bounded.

**Alternatives considered**:

- Run all lint versions on every OS/toolchain combination. Rejected unless the
  linter supports all versions; it multiplies failures unrelated to library
  behavior.
- Retain a historical version matrix. Rejected because historical Go releases
  are outside the repository's support boundary.
