# Data Model: Generic Parser Compatibility

## Existing entities retained

| Entity | Responsibility | Compatibility rule |
|--------|----------------|--------------------|
| `types.Type` | Common reflection contract. | Existing methods and observable behavior remain unchanged. |
| `types.Struct`, `types.Interface`, `types.Alias`, `types.Function` | Existing named declaration and signature representations. | Do not add exported fields or replace these concrete values for existing inputs. |
| `types.Package` and `types.PackageMap` | Own resolved declarations and cross-package lookup. | Existing lookups retain their names and results for non-generic source. |

## New additive entities

| Entity | Core fields / relationships | Validation rules |
|--------|-----------------------------|------------------|
| `TypeParameter` | Name, constraint, declaration owner, ordinal position. Belongs to a named type, function, or method receiver. | Name is non-empty unless deliberately blank in source; ordinal is unique within its owner; constraint is retained as a Gentools type representation. |
| `GenericInfo` | Declared type parameters and, when applicable, receiver type parameters. Addressable from the owning named declaration or function through an additive accessor. | Generic info for a non-generic owner is absent/empty; declaration order is preserved. |
| `Instantiation` | Origin declaration/type plus ordered, concrete type arguments. Implements `types.Type` when a type use is instantiated. | Argument count matches origin parameters when that fact is available; origin identity is preserved; different argument lists are not equal. |
| `Constraint` representation | Interface/type-set metadata, including embedded terms, unions, and approximation terms. Referenced by `TypeParameter`. | Preserve resolved semantic identity; unsupported representations produce a controlled error rather than a nil type or panic. |

## Relationships

```text
Package
├── named declaration (Struct | Interface | Alias | Function)
│   └── GenericInfo
│       ├── TypeParameter -> Constraint
│       └── receiver TypeParameter(s), where present
└── declaration/field/signature type use
    └── Instantiation -> origin + ordered type arguments
```

## State and conversion rules

1. Scaffold named declarations before resolving fields, methods, aliases, and
   signatures, preserving the existing cycle-safe order.
2. During semantic conversion, recognize type parameters, named instances,
   aliases, constraints, and signatures before falling back to existing basic,
   composite, and named-type conversion.
3. Attach generic metadata after its owner is scaffolded, so recursive generic
   definitions can refer to their own parameters without a second package load.
4. Return an error result for an unsupported construct or unavailable
   toolchain-specific converter. Never store a partially initialized generic
   entity as a successful result.
5. For source with no generic declarations or uses, return the same public
   entity kinds and values as before.
