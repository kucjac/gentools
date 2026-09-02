# Public Type-Model Contract: Generic Support

## Scope

This contract defines additive behavior for the public `types` and `parser`
packages. It supplements, and does not replace, the existing public contracts
documented in [README.md](../../../README.md).

## Compatibility guarantees

- Existing exported identifiers, signatures, and documented non-generic output
  remain valid on the selected Go 1.27.1 toolchain.
- Existing named declarations continue to be represented by their established
  concrete Gentools type (`Struct`, `Interface`, `Alias`, or `Function`) rather
  than a new wrapper solely because the declaration is generic.
- Generic metadata is additive and discoverable through new documented
  accessors; callers that do not use it retain prior behavior.
- A resolved generic type use with concrete type arguments is distinguishable
  from its origin and from the same origin instantiated with another argument
  list.

## Required observations

| Source construct | Required public observation |
|------------------|-----------------------------|
| Generic type/function declaration | Its parameter names, declaration order, constraints, and owner are available. |
| Type parameter appearing in a field, parameter, result, receiver, or constraint | The parameter is represented as a named type reference associated with its declaration scope. |
| Instantiated named type | The origin and ordered type arguments are available and retain normal `Type` rendering/identity behavior. |
| Generic alias (Go 1.24+) | Alias identity, parameters, constraints, right-hand side, and type arguments are available without confusing it with a defined type. |
| Method on a generic receiver type | Receiver metadata, ordinary parameters/results, variadic status, and method name are retained. Go methods do not declare their own type parameters. |

## Failure behavior

- Source that the active Go toolchain cannot parse or type-check is returned as
  a normal package-loading failure with the compiler/tooling diagnostic intact.
- A generic representation unsupported by the active Gentools converter
  produces an actionable error that identifies the construct; it must not
  panic, silently discard type arguments, or return a successful partial value.
- Validation covers Go 1.27.1 only; historical toolchains are outside scope.

## Non-goals

- This feature does not expose a general expression-indexing API for every
  generic call or selector in a package body.
- This feature does not change generated protobuf code or introduce a code
  generator.
- This feature does not promise support for a future Go language change until
  it has a fixture and a supported CI toolchain entry.
