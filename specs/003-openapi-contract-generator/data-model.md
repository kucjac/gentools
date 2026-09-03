# Data Model: OpenAPI Contract Generator

## Operation annotation

An exported function comment contains one directive with method, route,
summary, and status/model response pairing. The function’s package and name
identify the source declaration in diagnostics. The method/route pair is a
unique key.

## Operation

An operation has an identifier, HTTP method, route, summary, route parameters,
and one success response referencing a model. It becomes a path item and
operation object in the contract.

## Model definition

A public named struct has a fully qualified identity and ordered exported
fields. Each field has a serialized name, description from its source comment,
and a supported schema shape. Definitions are keyed by qualified model name
and emitted once under `components.schemas`.

## Contract document

The document has `openapi: 3.0.3`, an `info` identity, `paths`, and reusable
`components.schemas`. JSON output is deterministic and only written after all
validation succeeds.

## Validation rules

Missing or malformed directives, invalid methods/routes, duplicate method-route
pairs, unresolved models, and unsupported field shapes are errors. Every error
names the operation or model/field declaration. Validation is transactional:
the destination is not created or replaced on failure.
