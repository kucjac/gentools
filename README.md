# Golang Gentools

[![GitHub license](https://img.shields.io/github/license/kucjac/gentools.svg)](https://github.com/kucjac/gentools/blob/master/LICENSE)
![Workflow Status](https://github.com/kucjac/gentools/actions/workflows/ci.yml/badge.svg)
[![GoReportCard example](https://goreportcard.com/badge/github.com/kucjac/gentools)](https://goreportcard.com/report/github.com/kucjac/gentools)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/kucjac/gentools)
[![GitHub release](https://img.shields.io/github/release/kucjac/gentools.svg)](https://GitHub.com/kucjac/gentools/releases/)

![Gopher](./gopher.png)

Package `gentools` contains Golang tools used for scanning and parsing golang code.
It is based on the golang `ast` and provides simple but very powerful functionality in terms of code scanning.

The API is expected to be simpler and less verbose than `go/types` package, that looks more like `reflect`.

### Parsing packages

This package allows reading and scanning the content of go files based on the non-runtime AST definitions.

The `types.PackageMap` contains context of mapped packages with their structs and dependencies.
 
The map could get loaded by using `LoadPackages` function or `types.PackageMap` method with the same name.
Both of these loads provided input packages, whereas the method gets only required that doesn't already exist in itself.

Example:

```go
package main

import (
	"fmt"
	"os"

	"github.com/kucjac/gentools/types"
	"github.com/kucjac/gentools/parser"
)

func main() {
	pkgs, err := parser.LoadPackages(parser.LoadConfig{
		// Paths should contain file system paths to related golang directories.			
		// Paths: []string{"/home/user/golang/src/github.com/kucjac/gentools", "./../mypackage"},
		
		// PkgNames should contain full package names to get. It could be set along with the 'Paths' field.
		PkgNames:   []string{"github.com/kucjac/gentools"},
		BuildFlags: nil,
		Verbose:    false,
	})
	if err != nil {
		fmt.Printf("Err: Loading packages failed: %s\n", err)
		os.Exit(1)
	}
}
```


### Extracting declarations

Based on loaded packages a developer can operate on the types loaded along with these packages.

```go
// Let's get the types.Type interface and check it's content.
t, ok := pkgs.TypeOf("types.Type", nil)
if !ok {
    fmt.Println("Err: getting types.Type failed")
    os.Exit(1)
}

// We can get the name of given type, with or without it's identifier.
fmt.Println(t.Name(true, "")) // types.Type
fmt.Println(t.Name(false, "")) // Type

// We can set up current working package context while getting the name. 
// This way the result should not contain the identifier.
// The package context could be an identifier or even better full package name.
fmt.Println(t.Name(true, "types")) // Type
```


While the packages got loaded not only selected packages were provided but also all dependency imports.

```go
mutexType, ok := pkgs.TypeOf("*sync.Mutex", nil)
if !ok {
	fmt.Println("Err: getting *sync.Mutex failed")
	os.Exit(1)
}

fmt.Println(mutexType.Name(true, "")) // *sync.Mutex
// We can dereference the pointer, slice, array, channel or wrapper by using 'Elem' method. 
mutexType = mutexType.Elem()
fmt.Println(mutexType.Name(true, "")) // sync.Mutex
```

The types allows to easily create and operate on given types with very simple API.

```go
structType, ok := pkgs.TypeOf("types.Struct")
if !ok {
	fmt.Println("Err: getting types.Struct failed")
	os.Exit(1)
}

// Let's create a chan of slices of pointers to the given Struct.
newType := types.ChanOf(types.RecvOnly, types.SliceOf(types.PointerTo(structType)))
fmt.Println(sliceType) // chan <-[]*types.Struct
```

The package reads all methods and structure fields for loaded types along with the metadata like receiver
type (methods) and name, or parameter names.

```go
sliceType, ok := pkgs.TypeOf("types.Array")
if !ok {
	fmt.Println("Err: getting types.Struct failed")
	os.Exit(1)
}

st := sliceType.(*types.Struct)
for i, sField := range st.Fields {
	fmt.Printf("Field: %d - %s\t%s\n", i, sField.Name, sField.Type)
}

for i, method := range st.Methods {
	fmt.Printf("Method: %d - %s\n", i, method)
}
```

Each package could also be extracted out of all packages. It is also possible to extract any declaration out of the specific package in a way like:
```go
const typesPkgName = "github.com/kucjac/gentools/types"
typesPkg, ok := pkgs.PackageByPath(typesPkgName)
if !ok {
	fmt.Printf("Err: types package: %s not parsed\n", typesPkgName)
	os.Exit(1)
}

fmt.Println(typesPkg.GetPkgPath()) // github.com/kucjac/gentools/types
fmt.Println(typesPkg.Identifier) // types

parsedPkg, ok := typesPkg.GetType("Package") // This gets parsed types.Package type.
if !ok {
	fmt.Println("Err: package types can't get the type 'Package'")
	os.Exit(1)
}

```

### Comparing types

Each type could be compared to the other in two ways: 

#### Compare equality:

Two types could be compared on their equality by calling the `Equal` method of one type and providing another as an argument.

If two types are exactly the same then the value would result as true.

i.e.:
```go
// Check if the type String is equal to the type Bool 
fmt.Println(types.String.Equal(types.Bool)) // false

// But exactly the same types would give a positive result i.e:
fmt.Println(types.String.Equal(types.String)) // true

// This also applies to any other type - slice, struct, alias, interface, map, pointer etc i.e.:
fmt.Println(types.SliceOf(types.String).Equal(types.SliceOf(types.String))) // true
```

It is important to see that the aliases would not be equal to the type they point to i.e.:

```go
type AliasOfInt int

// Load this package and compare types ...

fmt.Println(aliasOfIntType.Equal(types.Int)) // false

// But if we try to derefence this alias the value would be positive:

fmt.Println(aliasOfIntType.Elem().Equal(types.Int)) // true
```

#### Compare by Kind

In some situations if we don't want to check if the types are equal directly, we can compare the kind of the type.

As an example, we might want to check if the type that we're looking for is a pointer, and if we need to dereference it.

All dereferences are done by using `Elem()` method, which panics if there is nothing to dereference - you cannot dereference basic types, structs etc.

These are supported type kinds:

- **KindBool**
- **KindInt**
- **KindInt8**
- **KindInt16**
- **KindInt32**
- **KindInt64**
- **KindUint**
- **KindUint8**
- **KindUint16**
- **KindUint32**
- **KindUint64**
- **KindUintptr**
- **KindFloat32**
- **KindFloat64**
- **KindComplex64**
- **KindComplex128**
- **KindString**
- **KindArray**
- **KindChan**
- **KindFunc**
- **KindInterface**
- **KindMap**
- **KindPtr**
- **KindSlice**
- **KindStruct**
- **KindUnsafePointer**


In order to compare some types kinds simply use the `Kind` method of the type and compare it  like: 

```go
fmt.Println(type1.Kind() == type2.Kind()) // KindBool == KindBool = true 
```

This might be super helpful in terms of comparing the kind of aliased types. I.e. if we have alias type mentioned above:

```go
// AliasOfInt is a type wrapper over the `int` standard type.
type AliasOfInt int

// As previously known this would fail
fmt.Println(aliasOfIntType.Equal(types.Int)) // false

// But we can compare it's kind which would result
fmt.Println(aliasOfIntType.Kind() == types.KindInt) // true
```

#### Interface comparison

This package supports also interfaces comparison. The parser deeply checks all methods and provides a way to check if 
an type (struct, interface) implements also another interface.

i.e.:

```go
osFile, ok := pkgs.TypeOf("*os.File")
if !ok {
	fmt.Println("*os.File not found")
	os.Exit(1)
}


ioReader, ok := pkgs.TypeOf("io.Reader")
if !ok {
	fmt.Println("io.Reader not found")
	os.Exit(1)
}

// The function types.Implements allows, checking if a type implements provided interface. 
fmt.Println(types.Implements(osFile, io.Reader.(*types.Interface))) // true

ioReadCloser, ok := pkgs.TypeOf("io.ReadCloser")
if !ok {
	fmt.Println("io.ReadCloser not found")
	os.Exit(1)
}

// It could also be an interface - io.ReadCloser implements io.Reader interface.
fmt.Println(types.Implements(ioReadCloser, io.Reader.(*types.Interface))) // true
```