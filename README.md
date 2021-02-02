# Golang Gentools

![Gopher](./gopher.png)

Package `gentools` contains Golang tools used for the .

The API is expected to be simpler and less verbose than `go/types` package, that looks more like `reflect`.


## Usage

This package allows reading and scanning the content of go files based on the non-runtime AST definitions.

The `astreflect.PackageMap` contains context of mapped packages with their structs and dependencies.
 
The map could get loaded by using `LoadPackages` function or `astreflect.PackageMap` method with the same name.
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
		// Paths: []string{"/home/user/golang/src/github.com/kucjac/astreflect", "./../mypackage"},
		
		// PkgNames should contain full package names to get. It could be set along with the 'Paths' field.
		PkgNames:   []string{"github.com/kucjac/astreflect"},
		BuildFlags: nil,
		Verbose:    false,
	})
	if err != nil {
		fmt.Printf("Err: Loading packages failed: %s\n", err)
		os.Exit(1)
	}
}
```


Based on loaded packages a developer can operate on the types loaded along with these packages.

```go
// Let's get the astreflect.Type interface and check it's content.
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
structType, ok := pkgs.TypeOf("types.StructType")
if !ok {
	fmt.Println("Err: getting types.StructType failed")
	os.Exit(1)
}

// Let's create a chan of slices of pointers to the given StructType.
newType := types.ChanOf(types.RecvOnly, types.SliceOf(types.PointerTo(structType)))
fmt.Println(sliceType) // chan <-[]*types.StructType
```

The package reads all methods and structure fields for loaded types along with the metadata like receiver
type (methods) and name, or parameter names.

```go
sliceType, ok := pkgs.TypeOf("types.ArrayType")
if !ok {
	fmt.Println("Err: getting types.StructType failed")
	os.Exit(1)
}

st := sliceType.(*types.StructType)
for i, sField := range st.Fields {
	fmt.Printf("Field: %d - %s\t%s\n", i, sField.Name, sField.Type)
}

for i, method := range st.Methods {
	fmt.Printf("Method: %d - %s\n", i, method)
}
```


