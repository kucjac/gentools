package astreflect

var _ Type = &FunctionType{}

// FunctionType is the function type used for getting.
type FunctionType struct {
	Comment     string
	PackagePath PkgPath
	Receiver    *Receiver
	FuncName    string
	In          []IOParam
	Out         []IOParam
	Variadic    bool
}

// Name implements Type interface.
func (f *FunctionType) Name(identified bool) string {
	if identified {
		if i := f.PackagePath.Identifier(); i != "" {
			return i + "." + f.FuncName
		}
	}
	return f.FuncName
}

// FullName implements Type interface.
func (f *FunctionType) FullName() string {
	return string(f.PackagePath) + "/" + f.FuncName
}

func (f *FunctionType) PkgPath() PkgPath {
	return f.PackagePath
}

func (f *FunctionType) Kind() Kind {
	return Func
}

func (f *FunctionType) Elem() Type {
	return nil
}
