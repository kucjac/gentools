package types

import (
	"testing"
)

func TestPackageMap_TypeOf(t *testing.T) {
	pkgMap := PackageMap{}

	const (
		testPkgPath        = "mytesting.com/package/pkg"
		testPkgIdentifier  = "pkg"
		testPtrTypeName    = "TestingPtrType"
		testIntTypeName    = "TestingIntType"
		testStringTypeName = "TestingStringType"
	)
	testPkg := NewPackage(testPkgPath, testPkgIdentifier)
	pkgMap[testPkgPath] = testPkg

	ptrAlias := &Alias{AliasName: testPtrTypeName, Type: PointerTo(Int32)}
	intAlias := &Alias{AliasName: testIntTypeName, Type: Int32}
	strAlias := &Alias{AliasName: testStringTypeName, Type: String}
	err := testPkg.NewNamedType(testPtrTypeName, ptrAlias)
	if err != nil {
		t.Fatalf("creating new named type failed: %f", err)
	}
	err = testPkg.NewNamedType(testIntTypeName, intAlias)
	if err != nil {
		t.Fatalf("creating new named type failed: %f", err)
	}
	err = testPkg.NewNamedType(testStringTypeName, strAlias)
	if err != nil {
		t.Fatalf("creating new named type failed: %f", err)
	}

	t.Run("BuiltIn", func(t *testing.T) {
		testCases := []Type{
			Int, Int8, Int16, Int32, Int64,
			Uint, Uint8, Uint16, Uint32, Uint64,
			String,
			Float32, Float64,
			Byte, Rune,
			Bool,
			Complex64, Complex128,
		}

		t.Run("Simple", func(t *testing.T) {
			testCaseKinds := []Kind{
				KindInt, KindInt8, KindInt16, KindInt32, KindInt64,
				KindUint, KindUint8, KindUint16, KindUint32, KindUint64,
				KindString,
				KindFloat32, KindFloat64,
				KindUint8, KindInt32,
				KindBool,
				KindComplex64, KindComplex128,
			}
			for i, tc := range testCases {
				t.Run(tc.String(), testPackagaMapTypeOfBuiltInSimple(pkgMap, tc, testCaseKinds, i))
			}
		})
		t.Run("Pointer", func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.String(), testPackageMapTypeOfBuiltInPointer(pkgMap, tc))
			}
		})
		t.Run("Array", func(t *testing.T) {
			for _, tc := range testCases {
				t.Run("[3]"+tc.String(), testPackageMapTypeOfBuiltInArray(tc, pkgMap))
			}
		})
		t.Run("Slice", func(t *testing.T) {
			for _, tc := range testCases {
				t.Run("[]"+tc.String(), testPackageMapTypeOfBuiltInSlice(tc, pkgMap))
			}
		})
		t.Run("Map", func(t *testing.T) {
			for _, tc := range testCases {
				t.Run("map[string]"+tc.String(), testPackageMapTypeOfBuiltInMap(tc, pkgMap))
			}
		})
		t.Run("Chan", testPackageMapTypeOfBuiltInChans(testCases, pkgMap))
	})

	t.Run("CustomPackage", func(t *testing.T) {
		t.Run("FullPath", testPackageMapTypeOfCustomFullPackage(testPkgPath, testPtrTypeName, pkgMap, ptrAlias))
		t.Run("Identifier", testPackageMapTypeOfCustomIdentifier(testPkgIdentifier, testPtrTypeName, pkgMap, ptrAlias))
		t.Run("Context", testPackageMapTypeOfCustomContext(pkgMap, testPtrTypeName, testPkg, ptrAlias))
		t.Run("AliasZero", func(t *testing.T) {
			t.Run("IntBased", testPackageMapTypeOfCustomAliasInt(testIntTypeName, pkgMap, testPkg, intAlias))
			t.Run("PointerBased", testPackageMapTypeOfCustomAliasPointer(testPtrTypeName, pkgMap, testPkg, ptrAlias))
			t.Run("StringBased", testPackageMapTypeOfCustomALiasString(testStringTypeName, pkgMap, testPkg, strAlias))
		})
	})
}

func testPackageMapTypeOfBuiltInChans(testCases []Type, pkgMap PackageMap) func(t *testing.T) {
	return func(t *testing.T) {
		for _, tc := range testCases {
			dirs := []struct {
				name    string
				chanDir ChanDir
			}{
				{
					name:    "<-chan",
					chanDir: RecvOnly,
				},
				{
					name:    "chan<-",
					chanDir: SendOnly,
				},
				{
					name:    "chan",
					chanDir: SendRecv,
				},
			}
			for _, dir := range dirs {
				t.Run(dir.chanDir.String()+"/"+tc.String(), testPackageMapTypeOfBuiltInChanDir(pkgMap, dir.name, dir.chanDir, tc))
			}
		}
	}
}

func testPackageMapTypeOfCustomContext(pkgMap PackageMap, testPtrTypeName string, testPkg *Package, ptrAlias *Alias) func(t *testing.T) {
	return func(t *testing.T) {
		tp, ok := pkgMap.TypeOf(testPtrTypeName+"{}", testPkg)
		if !ok {
			t.Fatalf("can't find type %s with context pkg", testPtrTypeName)
		}
		if !tp.Equal(ptrAlias) {
			t.Errorf("expected: %s but is: '%s'", ptrAlias, tp)
		}
	}
}

func testPackageMapTypeOfCustomALiasString(testStringTypeName string, pkgMap PackageMap, testPkg *Package, strAlias *Alias) func(t *testing.T) {
	return func(t *testing.T) {
		var testAliasType = testStringTypeName + "(\"\")"
		tp, ok := pkgMap.TypeOf(testAliasType, testPkg)
		if !ok {
			t.Fatalf("can't find type: '%s'", testAliasType)
		}
		if !tp.Equal(strAlias) {
			t.Errorf("expected: '%s' but is: '%s'", tp, strAlias)
		}
	}
}

func testPackageMapTypeOfCustomAliasPointer(testPtrTypeName string, pkgMap PackageMap, testPkg *Package, ptrAlias *Alias) func(t *testing.T) {
	return func(t *testing.T) {
		var testAliasType = testPtrTypeName + "(nil)"
		tp, ok := pkgMap.TypeOf(testAliasType, testPkg)
		if !ok {
			t.Fatalf("can't find type: '%s'", testAliasType)
		}
		if !tp.Equal(ptrAlias) {
			t.Errorf("expected: '%s' but is: '%s'", tp, ptrAlias)
		}
	}
}

func testPackageMapTypeOfCustomAliasInt(testIntTypeName string, pkgMap PackageMap, testPkg *Package, intAlias *Alias) func(t *testing.T) {
	return func(t *testing.T) {
		var testAliasType = testIntTypeName + "(0)"
		tp, ok := pkgMap.TypeOf(testAliasType, testPkg)
		if !ok {
			t.Fatalf("can't find type: '%s'", testAliasType)
		}
		if !tp.Equal(intAlias) {
			t.Errorf("expected: '%s' but is: '%s'", tp, intAlias)
		}
	}
}

func testPackageMapTypeOfCustomIdentifier(testPkgIdentifier string, testPtrTypeName string, pkgMap PackageMap, ptrAlias *Alias) func(t *testing.T) {
	return func(t *testing.T) {
		var testTypeIdentName = testPkgIdentifier + "." + testPtrTypeName + "{}"
		tp, ok := pkgMap.TypeOf(testTypeIdentName, nil)
		if !ok {
			t.Fatalf("can't find type with pkg identifier")
		}
		if !tp.Equal(ptrAlias) {
			t.Errorf("expected: %s but is: '%s'", ptrAlias, tp)
		}
	}
}

func testPackageMapTypeOfCustomFullPackage(testPkgPath string, testPtrTypeName string, pkgMap PackageMap, ptrAlias *Alias) func(t *testing.T) {
	return func(t *testing.T) {
		var testTypeFullName = testPkgPath + "." + testPtrTypeName + "{}"
		tp, ok := pkgMap.TypeOf(testTypeFullName, nil)
		if !ok {
			t.Fatalf("can't find type with pkg path")
		}
		if !tp.Equal(ptrAlias) {
			t.Errorf("expected: %s but is: '%s'", ptrAlias, tp)
		}
	}
}

func testPackageMapTypeOfBuiltInChanDir(pkgMap PackageMap, dirName string, chanDir ChanDir, tc Type) func(t *testing.T) {
	return func(t *testing.T) {
		tp, ok := pkgMap.TypeOf(dirName+" "+tc.Name(false, ""), nil)
		if !ok {
			t.Fatalf("cannot find type Of: %s", dirName+tc.Name(false, ""))
		}

		if !tp.Equal(ChanOf(chanDir, tc)) {
			t.Errorf("Type: %s is not a %s", tc, ChanOf(chanDir, tc))
		}
		if tp.Kind() != KindChan {
			t.Errorf("Expected kind to be chan but is: %s", tp.Kind())
		}
		if !tp.Elem().Equal(tc) {
			t.Errorf("Elem of chan: %s is not equal to %s", tp.Elem(), tc)
		}
	}
}

func testPackageMapTypeOfBuiltInMap(tc Type, pkgMap PackageMap) func(t *testing.T) {
	return func(t *testing.T) {
		typeToCheck := "map[string]" + tc.Name(false, "")
		tp, ok := pkgMap.TypeOf(typeToCheck, nil)
		if !ok {
			t.Fatalf("array of %s builtin type is not found", typeToCheck)
		}

		if !tp.Equal(MapOf(String, tc)) {
			t.Errorf("%s is not equal to %s", tp, MapOf(String, tc))
		}
	}
}

func testPackageMapTypeOfBuiltInSlice(tc Type, pkgMap PackageMap) func(t *testing.T) {
	return func(t *testing.T) {
		typeToCheck := "[]" + tc.Name(false, "")
		tp, ok := pkgMap.TypeOf(typeToCheck, nil)
		if !ok {
			t.Fatalf("array of %s builtin type is not found", typeToCheck)
		}

		if !tp.Equal(SliceOf(tc)) {
			t.Errorf("%s is not equal to %s", tp, SliceOf(tc))
		}
	}
}

func testPackageMapTypeOfBuiltInArray(tc Type, pkgMap PackageMap) func(t *testing.T) {
	return func(t *testing.T) {
		typeToCheck := "[3]" + tc.Name(false, "")
		tp, ok := pkgMap.TypeOf(typeToCheck, nil)
		if !ok {
			t.Fatalf("array of %s builtin type is not found", typeToCheck)
		}

		if !tp.Equal(ArrayOf(tc, 3)) {
			t.Errorf("%s is not equal to %s", tp, ArrayOf(tc, 3))
		}
	}
}

func testPackageMapTypeOfBuiltInPointer(pkgMap PackageMap, tc Type) func(t *testing.T) {
	return func(t *testing.T) {
		tp, ok := pkgMap.TypeOf("*"+tc.Name(false, ""), nil)
		if !ok {
			t.Fatalf("%s built in type is not found", tp)
		}

		if tp.Kind() != KindPtr {
			t.Fatalf("expected kind to be pointer but is: '%s'", tp.Kind())
		}
		elem := tp.Elem()

		if !elem.Equal(tc) {
			t.Errorf("%s is not equal to its base type", tc)
		}
	}
}

func testPackagaMapTypeOfBuiltInSimple(pkgMap PackageMap, tc Type, testCaseKinds []Kind, i int) func(t *testing.T) {
	return func(t *testing.T) {
		tp, ok := pkgMap.TypeOf(tc.Name(false, ""), nil)
		if !ok {
			t.Fatalf("%s built in type is not found", tp)
		}

		if !tp.Equal(tc) {
			t.Errorf("%s is not equal to its base type", tc)
		}
		if tp.Kind() != testCaseKinds[i] {
			t.Errorf("%s kind is not equal to %s", tp.Kind(), testCaseKinds[i])
		}
	}
}
