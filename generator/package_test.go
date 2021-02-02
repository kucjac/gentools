package genutils

import (
	"fmt"
	"testing"

	"github.com/kucjac/gentools/parser"
	"github.com/kucjac/gentools/types"
)

func TestPackage(t *testing.T) {
	generatePath := "github.com/kucjac/gentools/generate"
	typesPath := "github.com/kucjac/gentools/types"
	pm, err := parser.LoadPackages(parser.LoadConfig{
		PkgNames: []string{generatePath, "fmt", typesPath},
		Verbose:  false,
	})
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	generatePkg := pm.MustGetByPath(generatePath)
	typesPkg := pm.MustGetByPath(typesPath)

	p, err := NewPackage(pm, "github.com/kucjac/gentools/gentest", "gentest")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	fmtPackage := pm.MustGetByPath("fmt")
	// Define the functions that would get executed.
	fmtErrorf := fmtPackage.MustFunction("Errorf")

	// type Package struct {
	pkgStruct, err := p.Struct("Package", func(d *StructDef) {
		// Fields
		// tp *generate.Package
		d.Field("tp", types.PointerTo(typesPkg.MustGetType("Package")))
		d.Field("pm", typesPkg.MustGetType("PackageMap"))
	})
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	// Methods definitions.
	err = pkgStruct.Method("Func", func(fc MethodCreator) {
		fc.Receiver("p", true)

		fc.In("name", types.String)
		fc.InFunc("contentFunc", func(ft UnnamedFuncer) {
			ft.In("c", types.PointerTo(generatePkg.MustGetType("FuncContent")))
		})

		fc.Out("", types.PointerTo(generatePkg.MustGetType("FuncDef")))
		fc.Out("", types.Error)

		// Generate content.
		fc.P("if _, ok := p.tp.Types[name]; ok {")
		fc.Ret(fc.Ind(), nil, fc.ExecFn(fmtErrorf, func(x Executor) {
			x.Arg(fc.Q("package: '%s' has already a type with a name: '%s'"))
			x.Arg("p.tp.Path")
			x.Arg("name")
		}))
		fc.P("}")

		fc.P("if _, ok := p.tp.Declarations[name]; ok {")
		fc.Ret(fc.Ind(), nil, fc.ExecFn(fmtErrorf, func(x Executor) {
			x.Arg(fc.Q("package: '%s' has already a type with a name: '%s'"))
			x.Arg("p.tp.Path")
			x.Arg("name")
		}))
		fc.P("}")

		fc.Ret("&FuncDef{tp: &types.Function{Pkg: p.tp, FuncName: name}, contentFunc: contentFunc}", nil)
	})
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	// Struct Method.
	err = pkgStruct.Method("p", true, "Struct", func(fc *funcContent) {
		fc.In("name", types.String)
		fc.InFunc("defFunc", func(ft UnnamedFuncer) {
			ft.In("d", generatePkg.MustGetType("StructDef"))
		})

		fc.StructVal(typesPkg.MustStruct("Struct"), func(s *structValuer) {
			s.Ptr()
			s.Field("Pkg", "p.tp")
			s.Field("TypeName", "name")
		})
	})
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

}
