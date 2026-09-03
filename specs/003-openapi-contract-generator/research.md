# Research: OpenAPI Contract Generator

## Decision: Generate OpenAPI 3 JSON with the standard library

The example emits a deterministic OpenAPI 3.0 JSON object using `encoding/json`
and sorted collections where necessary. JSON is directly reviewable and avoids
adding a YAML dependency to an educational example.

Alternatives considered: YAML or a full OpenAPI library were rejected because
they add dependencies and obscure the annotation/type-mapping example.

## Decision: Use Gentools metadata plus function comments

The generator loads the input package with `parser.LoadPackages` and
`WithComments: true`. Operation declarations are exported functions; their
comments contain one `@openapi` directive. Struct comments and fields come from
the existing `types.Struct`, `StructField`, and `StructTag` models.

## Decision: Small explicit annotation grammar

The supported directive is:

`@openapi method=GET route=/pets/{id} summary="Get a pet" response=200:Pet`

Required keys are `method`, `route`, `summary`, and `response`. A directive may
contain one response pairing in the first implementation. The example supports
GET; routes must begin with `/` and use named `{param}` segments. Unknown,
duplicate, malformed, or missing keys fail validation.

## Decision: Validate before writing and replace atomically

Generation builds the complete in-memory document, validates all operations,
models, references, and supported field shapes, then writes to a temporary file
in the destination directory and renames it. A validation error leaves an
existing destination unchanged and does not create a new destination.

## Supported model mapping

Exported struct fields use their `json` tag name; `json:"-"` excludes a field,
and absent tags use the Go field name. Primitive Go types map to OpenAPI scalar
types, pointers preserve the underlying schema, slices map to arrays, and
struct fields map to reusable component references. Unsupported maps,
interfaces, channels, functions, and unnamed shapes produce diagnostics naming
the model and field.
