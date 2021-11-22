package types

import (
	"fmt"
	"go/constant"
	"strconv"
)

// Declaration is the variable or constant declaration.
type Declaration struct {
	Comment  string
	Name     string
	Type     Type
	Constant bool
	Val      constant.Value
	Package  *Package
}

// ConstValue gets the basic value of given constant declaration type.
// The method panics if the Declaration is not a constant but variable.
// For selected field kind it returns following Value Types:
//	- KindString 												- string
//	- KindBool													- bool
//	- KindInt, KindInt8, KindInt16, KindInt32, KindInt64: 		- int
//	- KindUint, KindUint8, KindUint16, KindUint32, KindUint64: 	- uint
//	- KindFloat64, KindFloat32:									- float64
//	- KindComplex64, KindComplex128:							- complex128
func (d Declaration) ConstValue() interface{} {
	if !d.Constant {
		panic(fmt.Sprintf("declaration is not a constant: %s", d))
	}

	switch d.Type.Kind() {
	case KindString:
		return d.Val.String()
	case KindBool:
		v, _ := strconv.ParseBool(d.Val.String())
		return v
	case KindInt, KindInt8, KindInt16, KindInt32, KindInt64:
		i, _ := strconv.ParseInt(d.Val.String(), 10, 64)
		return int(i)
	case KindUint, KindUint8, KindUint16, KindUint32, KindUint64:
		u, _ := strconv.ParseUint(d.Val.String(), 10, 64)
		return uint(u)
	case KindFloat64, KindFloat32:
		v, _ := strconv.ParseFloat(d.Val.String(), 64)
		return v
	case KindComplex64, KindComplex128:
		v, _ := strconv.ParseComplex(d.Val.String(), 128)
		return v
	default:
		panic(fmt.Sprintf("Uknown constant type Kind: %s", d.Type.Kind()))
	}
}

// String provides a string visual form of the declaration.
func (d Declaration) String() string {
	if !d.Constant {
		return fmt.Sprintf("var %s.%s", d.Package.Identifier, d.Name)
	}
	return fmt.Sprintf("const %s.%s(%s)", d.Package.Identifier, d.Name, d.Val.String())
}
