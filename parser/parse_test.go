package parser

import (
	"go/constant"
	"testing"

	"github.com/kucjac/gentools/types"
)

func TestZero(t *testing.T) {
	pkgs, err := LoadPackages(LoadConfig{
		Paths:      []string{"."},
		PkgNames:   nil,
		BuildFlags: nil,
		Verbose:    false,
	})
	if err != nil {
		t.Errorf("Loading packages failed: %v", err)
		return
	}

	thisPkg, ok := pkgs.PackageByPath("github.com/kucjac/gentools/types")
	if !ok {
		t.Error("This package not found by path")
		return
	}
	thisPkg.SetIdentifier("gentypes")

	_, ok = pkgs.TypeOf("types.Type", nil)
	if !ok {
		t.Error("TypeOf find types.Type failed")
	}
	_, ok = pkgs.TypeOf("types.Struct", nil)
	if !ok {
		t.Error("TypeOf find 'types.Struct' failed")
	}

	for _, expected := range []struct {
		PkgContext *types.Package
		Name       string
		Type       types.Type
	}{
		{
			Name:       "Type",
			Type:       thisPkg.MustGetType("Type"),
			PkgContext: thisPkg,
		},
		{
			Name: "gentypes.Type",
			Type: thisPkg.MustGetType("Type"),
		},
		{
			Name: "[]*gentypes.Struct",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Struct")}}
			}(),
		},
		{
			Name: "[]*Struct",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Struct")}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[]*Struct{}",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Struct")}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]*Interface",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Interface")}}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]*gentypes.Interface",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Interface")}}}
			}(),
		},
		{
			Name: "[3][]<- chan Interface",
			Type: func() types.Type {
				return &types.Array{ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{ArrayKind: types.KindSlice, Type: &types.Chan{Dir: types.RecvOnly, Type: thisPkg.MustGetType("Interface")}}}
			}(),
			PkgContext: thisPkg,
		},
		{
			Name: "[3][]chan <- *gentypes.Interface",
			Type: func() types.Type {
				return &types.Array{
					ArrayKind: types.KindArray, ArraySize: 3, Type: &types.Array{
						ArrayKind: types.KindSlice, Type: &types.Chan{
							Dir: types.SendOnly, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Interface")}},
					},
				}
			}(),
		},
		{
			Name:       "map[string][]*Struct",
			Type:       &types.Map{Key: &types.BuiltInType{BuiltInKind: types.KindString}, Value: &types.Array{ArrayKind: types.KindSlice, Type: &types.Pointer{PointedType: thisPkg.MustGetType("Struct")}}},
			PkgContext: thisPkg,
		},
		{
			Name: "[]gentypes.Kind",
			Type: &types.Array{ArrayKind: types.KindSlice, Type: thisPkg.MustGetType("Kind")},
		},
		{
			Name: "[]byte",
			Type: &types.Array{ArrayKind: types.KindSlice, Type: types.MustGetBuiltInType("byte")},
		},
	} {
		tp, ok := pkgs.TypeOf(expected.Name, expected.PkgContext)
		if !ok {
			t.Errorf("Package: TypeOf('%s') failed", expected.Name)
			continue
		}
		if !tp.Equal(expected.Type) {
			t.Errorf("TypeOf resulting type: '%s' is not equal to expected: '%s'", tp, expected.Type)
		}
	}
}

func TestParse(t *testing.T) {
	const testCasesPkgPath = "github.com/kucjac/gentools/parser/testcases"
	pkgs, err := LoadPackages(LoadConfig{PkgNames: []string{
		testCasesPkgPath,
		"encoding",
	}, Verbose: true})
	if err != nil {
		t.Errorf("Parsing packages failed: %v", err)
		return
	}

	pkg, ok := pkgs.PackageByPath(testCasesPkgPath)
	if !ok {
		t.Fatal("no test cases package found")
	}

	fooType, ok := pkg.GetType("Foo")
	if !ok {
		t.Fatal("type Foo not found")
	}

	fooStruct, ok := fooType.(*types.Struct)
	if !ok {
		t.Fatal("type Foo is not a *types.Struct")
	}

	enum, ok := pkg.GetType("Enumerated")
	if !ok {
		t.Fatal("type Enumerated is not found")
	}
	t.Run("Enumerated", func(t *testing.T) {
		enumAlias, ok := enum.(*types.Alias)
		if !ok {
			t.Fatal("type Enumerated is not an Alias")
		}
		if enumAlias.Kind() != types.KindInt {
			t.Fatal("type Enumerated kind is not a KindInt")
		}

		enumOne, ok := pkg.Declarations["EnumeratedOne"]
		if !ok {
			t.Fatal("no EnumeratedOne declaration found in the package")
		}
		if !enumOne.Type.Equal(enum) {
			t.Fatalf("EnumeratedOne type is not Enumerated: %v", enumOne.Type)
		}

		if enumOne.Comment != "EnumeratedOne defines a first enumerated type value.\n" {
			t.Fatalf("EnumeratedOne comment doesn't match: '%s'", enumOne.Comment)
		}

		if !enumOne.Constant {
			t.Fatal("EnumeratedOne should be a constant")
		}

		if enumOne.Val == nil {
			t.Fatal("EnumeratedOne constant value should be defined")
		}

		if enumOne.Val.Kind() != constant.Int || enumOne.Val.String() != "1" {
			t.Fatal("EnumeratedOne constant value should be of kind Integer and take value of 1")
		}

		v, ok := enumOne.ConstValue().(int)
		if !ok || v != int(1) {
			t.Fatalf("EnumeratedOne constant value doesn't match: %v", enumOne.ConstValue())
		}
	})

	fooID, ok := pkg.GetType("FooID")
	if !ok {
		t.Fatal("type FooID not found")
	}

	fooAlias, ok := pkg.GetType("FooAlias")
	if !ok {
		t.Fatal("type FooAlias not found")
	}

	t.Run("FooAlias", func(t *testing.T) {
		alias, isAlias := fooAlias.(*types.Alias)
		if !isAlias {
			t.Fatalf("type FooAlias is expected to be an alias but is: %T", fooAlias)
		}
		if alias.Type == nil {
			t.Fatal("type FooAlias type is nil")
		}
		if !alias.Type.Equal(fooType) {
			t.Fatal("foo alias type doesn't match Foo type")
		}
	})

	t.Run("FooPtrAlias", func(t *testing.T) {
		tp, ok := pkg.GetType("FooPtrAlias")
		if !ok {
			t.Fatal("FooPtrAlias is not found within packages")
		}

		alias, ok := tp.(*types.Alias)
		if !ok {
			t.Fatal("FooPtrAlias is expected to be an *types.Alias")
		}

		if !alias.Type.Equal(types.PointerTo(fooType)) {
			t.Fatal("FooPtrAlias is expected to be pointer to Foo type")
		}
	})

	t.Run("Weird", func(t *testing.T) {
		wd, ok := pkg.GetType("Weird")
		if !ok {
			t.Fatal("cant find type Weird")
		}

		alias, ok := wd.(*types.Alias)
		if !ok {
			t.Fatal("Weird is expected to be *types.Alias")
		}

		if !alias.Type.Equal(types.Int) {
			t.Fatal("Weird type is expected to be Int")
		}
	})

	t.Run("WeirdStruct", func(t *testing.T) {
		tp, ok := pkg.GetType("WeirdStruct")
		if !ok {
			t.Fatal("WeirdStruct is not found")
		}

		st, ok := tp.(*types.Struct)
		if !ok {
			t.Fatal("WeirdStruct is expected to be a structure")
		}

		if st.Comment != "WeirdStruct docs.\n" {
			t.Errorf("WeirdStruct comment doesn't match: %s", st.Comment)
		}

		if len(st.Fields) != 1 {
			t.Fatal("WeirdStruct is expected to contain one field")
		}

		nm := st.Fields[0]
		if nm.Comment != "Name doc.\n" {
			t.Errorf("WeirdStruct.Name field comment doesn't match: %s", nm.Comment)
		}
	})

	t.Run("FooID", func(t *testing.T) {
		if k := fooID.Kind(); k != types.KindInt64 {
			t.Errorf("type FooID is not of a KindInt64 but: %v", k)
		}

		tm, ok := pkgs.TypeOf("encoding.TextMarshaler", nil)
		if !ok {
			t.Fatalf("type encoding.TextMarshaler not found")
		}

		tu, ok := pkgs.TypeOf("encoding.TextUnmarshaler", nil)
		if !ok {
			t.Fatalf("type encoding.TextUnmarshaler not found")
		}

		tmInterface, ok := tm.(*types.Interface)
		if !ok {
			t.Fatalf("type encoding.TextMarshaler is not an interface, but: %T", tm)
		}

		tuInterface, ok := tu.(*types.Interface)
		if !ok {
			t.Fatalf("type encoding.TextUnmarshaler is not an interface, but: %T", tm)
		}

		// The non pointer type FooID should implement encoding.TextMarshaler.
		if !types.Implements(fooID, tmInterface) {
			t.Error("type FooID doesn't implement encoding.TextMarshaler interface")
		}

		// But non pointer FooID should not implement encoding.TextUnmarshaler.
		if types.Implements(fooID, tuInterface) {
			t.Error("type FooID should not implement encoding.TextUnmarshaler interface")
		}

		// The Pointer to FooID - *FooID, should implement encoding.TextUnmarshaler interface.
		ptrFooID := types.PointerTo(fooID)
		if !types.Implements(ptrFooID, tuInterface) {
			t.Error("pointer to FooID should implement encoding.TextUnmarshaler interface")
		}

		// And also it should implement encoding.TextMarshaler interface.
		if !types.Implements(ptrFooID, tmInterface) {
			t.Errorf("pointer to FooID should implement  encoding.TextMarshaler interface")
		}

		alias := fooID.(*types.Alias)
		for _, method := range alias.Methods {
			switch method.FuncName {
			case "UnmarshalText":
				if method.Comment != "UnmarshalText implements encoding.TextUnmarshaler interface.\n" {
					t.Errorf("fooID UnmarshalText comment doesn't match: %s", method.Comment)
				}
			case "MarshalText":
				if method.Comment != "MarshalText implements encoding.TextMarshaler interface.\n" {
					t.Errorf("fooID MarshalText comment doesn't match: %s", method.Comment)
				}
			}
		}
	})

	bar, ok := pkg.GetType("Bar")
	if !ok {
		t.Fatal("no Bar type found")
	}
	barStruct, ok := bar.(*types.Struct)
	if !ok {
		t.Fatal("type Bar is not a struct")
	}

	notEmpty, ok := pkg.GetType("NotEmpty")
	if !ok {
		t.Fatal("no NotEmpty type found")
	}
	notEmptyInterface, ok := notEmpty.(*types.Interface)
	if !ok {
		t.Fatalf("NotEmpty is expected to be *types.Interface but is: %T", notEmpty)
	}
	inheritMe, ok := pkg.GetType("InheritMe")
	if !ok {
		t.Fatal("no InheritMe interface found")
	}
	inheritMeInterface, ok := inheritMe.(*types.Interface)
	if !ok {
		t.Fatalf("InheritMe is expected to be *types.Interface but is: %T", inheritMe)
	}

	t.Run("InheritMe", func(t *testing.T) {
		i := inheritMeInterface
		// TODO: implement comment matching.
		if i.Comment != "InheritMe is an interface that will be inherited.\n" {
			t.Errorf("InheritMe comment not match: '%s'", i.Comment)
		}

		if len(i.Methods) != 1 {
			t.Fatalf("InheritMe should have exactly one method but have: %d", len(i.Methods))
		}

		if i.Methods[0].FuncName != "Inherited" {
			t.Errorf("InhertMe method is not 'Inherited', '%s'", i.Methods[0].FuncName)
		}
	})

	t.Run("NotEmpty", func(t *testing.T) {
		i := notEmptyInterface
		if len(i.Methods) != 2 {
			t.Fatalf("NotEmpty interface expected to have two methods but have: %d", len(i.Methods))
		}
		m := i.Methods[0]
		if m.FuncName != "Call" {
			t.Errorf("NotEmpty method name is not equal to Call: %v", m.FuncName)
		}
		if len(m.In) != 2 {
			t.Fatalf("NotEmpty Call method expected to have two argument but have: %v", len(m.In))
		}
		ctx := m.In[0]
		if ctx.Name != "ctx" {
			t.Errorf("first argument name is not equal to 'ctx': %v", ctx.Name)
		}
		ctxPkg := pkgs.MustGetByPath("context")
		ctxInterface := ctxPkg.MustGetType("Context")
		if !ctx.Type.Equal(ctxInterface) {
			t.Errorf("first argument type expected to be context.Context, but is: %v", ctx.Type)
		}

		options := m.In[1]
		if options.Name != "options" {
			t.Errorf("options argument name is not 'options': %v", options.Name)
		}
		if !options.Type.Equal(types.SliceOf(types.String)) {
			t.Errorf("options argument type is not '[]string' but: %v", options.Type)
		}
		if !m.Variadic {
			t.Error("method has '...' in last argument - expected to be variadic")
		}

		if len(m.Out) != 2 {
			t.Fatalf("Call output expected to have two variables returned but have: %d", len(m.Out))
		}

		n := m.Out[0]
		if n.Name != "n" {
			t.Errorf("Call first returned variable name is not 'n': %v", n.Name)
		}
		if n.Type != types.Int {
			t.Errorf("Call first returned variable type is not Int: %v", n.Type)
		}
		err := m.Out[1]
		if err.Name != "err" {
			t.Errorf("Call second returned variable name is not 'err': %v", err.Name)
		}
		if err.Type != types.Error {
			t.Errorf("Call second returned variable is not 'error': %v", err.Type)
		}

		inherited := i.Methods[1]
		if inherited.FuncName != "Inherited" {
			t.Errorf("NotEmpty should contain Inherited method")
		}
	})

	t.Run("Foo", func(t *testing.T) {
		for _, field := range fooStruct.Fields {
			switch field.Name {
			case "ID":
				ftName := field.Type.Name(false, "")
				if ftName != "FooID" {
					t.Errorf("field ID is not of type FooID, current type %s", ftName)
				}

				if field.Type.Kind() != types.KindInt64 {
					t.Errorf("FooID type is not of kind int64, %s", field.Type.Kind())
				}
				if tag := field.Tag.Get("json"); tag != "id" {
					t.Errorf("field 'ID', tag json is not equal to 'id', %v", tag)
				}
				if field.Comment != "ID is the foo field identifier.\n" {
					t.Errorf("field: 'ID', comment not match. Expected: 'ID is the foo field identifier.\\n' is '%s'", field.Comment)
				}
			case "String":
				ftName := field.Type.Name(false, "")
				if ftName != "string" {
					t.Errorf("field 'String' is not of type string, current type %s", ftName)
				}
				if field.Type.Kind() != types.KindString {
					t.Errorf("field 'String' type is not of kind string, %s", field.Type.Kind())
				}
				if tag := field.Tag.Get("custom"); tag != "name" {
					t.Errorf("field 'String', tag custom is not equal to 'name', %v", tag)
				}
			case "CustomName":
				ftName := field.Type.Name(false, "")
				if ftName != "string" {
					t.Errorf("field 'CustomName' is not of type string, current type %s", ftName)
				}
				if field.Type.Kind() != types.KindString {
					t.Errorf("field 'CustomName' type is not of kind string, %s", field.Type.Kind())
				}
			case "Bool":
				ftName := field.Type.Name(false, "")
				if ftName != "bool" {
					t.Errorf("field 'Bool' is not of type bool, current type %s", ftName)
				}
				if field.Type.Kind() != types.KindBool {
					t.Errorf("field 'Bool' type is not of kind bool, %s", field.Type.Kind())
				}
			case "Enumerated":
				if field.Type != enum {
					t.Errorf("field 'Enumerated' is not of type Enumerated, current type %s", field.Type)
				}
			case "Slice":
				sl, ok := field.Type.(*types.Array)
				if !ok {
					t.Errorf("field 'Slice' is not of types Array, %T", field.Type)
					continue
				}

				if sl.ArrayKind != types.KindSlice {
					t.Errorf("slice type is not types.KindSlice, is: %s", sl.ArrayKind)
				}
				if sl.ArraySize != 0 {
					t.Errorf("slice size should be zero but is: %v", sl.ArraySize)
				}

				if sl.Type != types.String {
					t.Errorf("slice base type is not a string, is: %s", sl.Type)
				}
			case "Float64":
				ftName := field.Type.Name(false, "")
				if ftName != "float64" {
					t.Errorf("field 'Float64' is not of type floa64, current type %s", ftName)
				}
				if field.Type.Kind() != types.KindFloat64 {
					t.Errorf("field 'Float64' type is not of kind bool, %s", field.Type.Kind())
				}
			case "Duration":
				a, ok := field.Type.(*types.Alias)
				if !ok {
					t.Errorf("expected field 'Duration' type to be an alias but is: %T", field.Type)
					continue
				}
				if a.AliasName != "Duration" {
					t.Errorf("alias name is expected to be Duration but is: %s", a.AliasName)
				}
				if a.Kind() != types.KindInt64 {
					t.Errorf("time.Duration kind expected to be int64, but is: '%v", a.Kind())
				}
				if zero := a.Zero(true, "github.com/kucjac/gentools/parser"); zero != "time.Duration(0)" {
					t.Errorf("time.Duration zero expected to be time.Duration(0), but is: %v", zero)
				}
			case "Bar":
				if field.Type.Kind() != types.KindPtr {
					t.Errorf("field 'Bar' expected to be of kind 'Ptr' but is: %v", field.Type.Kind())
				}
				tp := field.Type.Elem()

				if tp != barStruct {
					t.Errorf("field 'Bar' elem expected to be a type Bar, but is: %v", tp)
				}
			default:
				t.Fatalf("unknown field name: '%s'", field.Name)
			}
		}
	})

	t.Run("Bar", func(t *testing.T) {
		for _, field := range barStruct.Fields {
			switch field.Name {
			case "Map":
				if field.Type.Kind() != types.KindMap {
					t.Errorf("'Map' field expected to be of kind map but is: %v", field.Type.Kind())
					continue
				}
				mp := field.Type.(*types.Map)
				if mp.Key != types.String {
					t.Errorf("'Map' field key expected to be a String but is: %v", mp.Key)
				}
				if mp.Value != types.Byte {
					t.Errorf("'Map' field value expected to be a Byte but is: %v", mp.Value)
				}
			case "Time":
				timePkg, ok := pkgs.PackageByIdentifier("time")
				if !ok {
					t.Errorf("time package not found")
					continue
				}
				timeType, ok := timePkg.GetType("Time")
				if !ok {
					t.Errorf("no time.Time type found")
					continue
				}
				if !field.Type.Equal(timeType) {
					t.Errorf("field 'Time' type is not a time.Time but: %v", field.Type)
				}
			case "Any":
				if field.Type.Kind() != types.KindInterface {
					t.Errorf("field 'Any' expected to be of KindInterface but is: %v", field.Type.Kind())
					continue
				}
				i, ok := field.Type.(*types.Interface)
				if !ok {
					t.Errorf("field 'Any' type expected to be *types.Interface but is: %T", field.Type)
					continue
				}

				if !types.IsEmptyInterface(i) {
					t.Errorf("field 'Any' type is expected to be an empty interface")
				}
			case "ChanIn":
				if field.Type.Kind() != types.KindChan {
					t.Errorf("field 'ChanIn' type expected to be KindChan but is: %v", field.Type.Kind())
				}
				c, ok := field.Type.(*types.Chan)
				if !ok {
					t.Errorf("field type Type should be *types.Chan but is: %T", field.Type)
					continue
				}
				if c.Type != types.Int {
					t.Errorf("field 'ChanIn' type expected to be types.Int but is: %T", c.Type)
				}
				if c.Dir != types.SendOnly {
					t.Errorf("field 'ChanIn' channel direction expected to be: 'SendOnly' but is: '%v'", c.Dir)
				}

			case "ChanOut":
				if field.Type.Kind() != types.KindChan {
					t.Errorf("field 'ChanOut' type expected to be KindChan but is: %v", field.Type.Kind())
				}
				c, ok := field.Type.(*types.Chan)
				if !ok {
					t.Errorf("field type Type should be *types.Chan but is: %T", field.Type)
					continue
				}
				if c.Type != types.Int {
					t.Errorf("field 'ChanOut' type expected to be types.Int but is: %T", c.Type)
				}
				if c.Dir != types.RecvOnly {
					t.Errorf("field 'ChanOut' channel direction expected to be: 'RecvOnly' but is: '%v'", c.Dir)
				}

			case "Chan":
				if field.Type.Kind() != types.KindChan {
					t.Errorf("field 'Chan' type expected to be KindChan but is: %v", field.Type.Kind())
				}
				c, ok := field.Type.(*types.Chan)
				if !ok {
					t.Errorf("field type Type should be *types.Chan but is: %T", field.Type)
					continue
				}
				if c.Type != types.Int {
					t.Errorf("field 'Chan' type expected to be types.Int but is: %T", c.Type)
				}
				if c.Dir != types.SendRecv {
					t.Errorf("field 'Chan' channel direction expected to be: 'SendRecv' but is: '%v'", c.Dir)
				}

			case "Error":
				if field.Type.Kind() != types.KindInterface {
					t.Errorf("field 'Error' expected to be of KindInterface but is: %v", field.Type.Kind())
					continue
				}
				i, ok := field.Type.(*types.Interface)
				if !ok {
					t.Errorf("field 'Error' type expected to be *types.Interface but is: %T", field.Type)
					continue
				}

				if !field.Type.Equal(types.Error) {
					t.Errorf("field 'Error' type is expected to be an error interface, but is: %v", i)
				}
			case "NotEmpty":
				if field.Type.Kind() != types.KindInterface {
					t.Errorf("field 'NotEmpty' expected to be of KindInterface but is: %v", field.Type.Kind())
					continue
				}
				if !field.Type.Equal(notEmpty) {
					t.Errorf("field 'NotEmpty' type expected to be of type NotEmpty but is: %v", field.Type)
				}
			}
		}
	})
}
