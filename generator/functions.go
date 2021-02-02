package genutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/kucjac/gentools/types"
)

// Indent defines line indentation.
type Indent struct{ size uint8 }

// FuncDef is the function wrapper used for easy function writing.
type FuncDef struct {
	tp          *types.Function
	contentFunc func(fc FuncCreator)
	written     bool
	imports     map[string]*Package
}

// Done checks if given function definition was already defined within given package.
func (f *FuncDef) Done() bool {
	return f.written
}

// InputParams
func (f *FuncDef) InputParams() FuncParams {
	return f.tp.In
}

func (f *FuncDef) Content() ([]byte, error) {
	return nil, nil
}

// SetComment sets the function comment.
func (f *FuncDef) SetComment(comment string) {
	f.tp.Comment = comment
}

// UnnamedFuncer is an interface that generates unnamed function type definition.
type UnnamedFuncer interface {
	In(name string, param types.Type, variadic ...bool)
	Out(name string, param types.Type)
	Type() *types.Function
}

// unnamedFunc is the function type creator.
type unnamedFunc struct {
	tp                        *types.Function
	err                       error
	namedInputs, namedOutputs bool
	declarations              map[string]struct{}
}

func newUnnamedFunc() *unnamedFunc {
	return &unnamedFunc{
		tp:           &types.Function{},
		namedInputs:  false,
		namedOutputs: false,
		declarations: map[string]struct{}{},
	}
}

// In generates function type inputs.
func (f *unnamedFunc) In(name string, paramType types.Type, variadic ...bool) {
	if f.err != nil {
		return
	}
	if len(f.tp.In) != 0 && name != "" && !f.namedInputs {
		f.err = fmt.Errorf("one of the previous function('%s') inputs is unnamed", f.tp.Name(false, ""))
		return
	}
	isVariadic := len(variadic) == 1 && variadic[0]
	if isVariadic && len(f.tp.In) != 0 && f.tp.Variadic {
		f.err = fmt.Errorf("function could only have one last variadic parameter")
		return
	}
	if isVariadic {
		f.tp.Variadic = true
		paramType = types.SliceOf(paramType)
	}
	if name != "" {
		if _, ok := f.declarations[name]; ok {
			f.err = fmt.Errorf("one of the previous function('%s') variables has exact name: %s", f.tp.Name(false, ""), name)
			return
		}
		f.namedInputs = true
		f.declarations[name] = struct{}{}
	}
	f.tp.In = append(f.tp.In, types.FuncParam{Name: name, Type: paramType})
}

// Out creates function type output.
func (f *unnamedFunc) Out(name string, paramType types.Type) {
	if f.err != nil {
		return
	}
	if len(f.tp.Out) != 0 && name != "" && !f.namedOutputs {
		f.err = fmt.Errorf("one of the previous function('%s') inputs is unnamed", f.tp.Name(false, ""))
		return
	}
	if name != "" {
		if _, ok := f.declarations[name]; ok {
			f.err = fmt.Errorf("one of the previous function('%s') variables has exact name: %s", f.tp.Name(false, ""), name)
			return
		}
		f.namedOutputs = true
		f.declarations[name] = struct{}{}
	}
	f.tp.Out = append(f.tp.Out, types.FuncParam{Name: name, Type: paramType})
}

// Type gets the types.Type of given function.
func (f *unnamedFunc) Type() *types.Function {
	return f.tp
}

// FuncCreator is a
type FuncCreator interface {
	In(name string, param types.Type, variadic ...bool)
	InFunc(name string, fn func(c UnnamedFuncer), variadic ...bool)
	Out(name string, param types.Type)
	ExecFn(funcT *types.Function, executor func(x Executor)) FuncExecution
	Var(name string, vartype types.Type)
	Const(name string, value string, kind types.Kind)
	P(args ...interface{})
	Q(input string) string
	Ind(size ...uint) Indent
	Ret(args ...interface{})
}

// funcContent is the function content builder.
type funcContent struct {
	def                       *FuncDef
	methodType                types.Type
	err                       error
	namedInputs, namedOutputs bool
	content                   *contentBlock
	pkg                       *Package
	declarations              map[string]types.Declaration
}

// Receiver sets the function receiver definition.
func (f *funcContent) Receiver(name string, pointer bool) {
	r := &types.Receiver{Name: name, Type: f.methodType}
	if pointer {
		r.Type = types.PointerTo(r.Type)
	}
	f.def.tp.Receiver = r
}

// ExecFn implements FuncCreator interface.
func (f *funcContent) ExecFn(funcT *types.Function, executor func(x Executor)) FuncExecution {
	if f.err != nil {
		return FuncExecution{}
	}
	ex := newExecutor(f, funcT)
	executor(ex)
	if ex.err != nil {
		return FuncExecution{}
	}
	// Build the string representation.
	var buf bytes.Buffer
	buf.WriteString(ex.ft.Name(true, f.pkg.tp.Path))
	buf.WriteRune('(')
	if len(ex.args) != 0 {
		buf.WriteString(strings.Join(ex.args, ","))
	}
	buf.WriteRune(')')
	return FuncExecution{out: ex.outs, raw: buf.String()}
}

// InFunc creates new input argument with the inline function argument.
func (f *funcContent) InFunc(name string, fn func(c UnnamedFuncer), variadic ...bool) {
	if f.err != nil {
		return
	}
	ft := newUnnamedFunc()
	fn(ft)
	if ft.err != nil {
		f.err = ft.err
		return
	}
	f.In(name, ft.tp, variadic...)
	return
}

// In adds the next input parameter.
func (f *funcContent) In(name string, paramType types.Type, variadic ...bool) {
	if f.err != nil {
		return
	}
	if len(f.def.tp.In) != 0 && name != "" && !f.namedInputs {
		f.err = fmt.Errorf("one of the previous function('%s') inputs is unnamed", f.def.tp.Name(false, ""))
		return
	}
	isVariadic := len(variadic) == 1 && variadic[0]
	if isVariadic && len(f.def.tp.In) != 0 && f.def.tp.Variadic {
		f.err = fmt.Errorf("function could only have one last variadic parameter")
		return
	}
	if isVariadic {
		f.def.tp.Variadic = true
		paramType = types.SliceOf(paramType)
	}
	if name != "" {
		if _, ok := f.declarations[name]; ok {
			f.err = fmt.Errorf("one of the previous function('%s') variables has exact name: %s", f.def.tp.Name(false, ""), name)
			return
		}
		f.namedInputs = true
		f.declarations[name] = types.Declaration{Name: name, Type: paramType}
	}
	f.def.tp.In = append(f.def.tp.In, types.FuncParam{Name: name, Type: paramType})
	return
}

// Out adds the next output parameter.
func (f *funcContent) Out(name string, paramType types.Type) {
	if f.err != nil {
		return
	}
	if len(f.def.tp.Out) != 0 && name != "" && !f.namedOutputs {
		f.err = fmt.Errorf("one of the previous function('%s') inputs is unnamed", f.def.tp.Name(false, ""))
		return
	}
	if name != "" {
		if _, ok := f.declarations[name]; ok {
			f.err = fmt.Errorf("one of the previous function('%s') variables has exact name: %s", f.def.tp.Name(false, ""), name)
			return
		}
		f.namedOutputs = true
		f.declarations[name] = types.Declaration{Name: name, Type: paramType}
	}
	f.def.tp.Out = append(f.def.tp.Out, types.FuncParam{Name: name, Type: paramType})
}

// Var declares new variable with given name and type.
func (f *funcContent) Var(name string, tp types.Type) {
	if f.err != nil {
		return
	}
	if name == "" {
		f.err = errors.New("provided unnamed variable")
		return
	}
	if _, ok := f.declarations[name]; ok {
		f.err = fmt.Errorf("'%s' already declared within function(%s) definition", f.def.tp.Name(false, ""), name)
		return
	}

	f.P("var ", name, " ", tp.Name(true, f.pkg.tp.Path))
	f.declarations[name] = types.Declaration{Name: name, Type: tp}
	return
}

// Const creates a function constant definition.
func (f *funcContent) Const(name string, value string, kind types.Kind) {
	if f.err != nil {
		return
	}
	if name == "" {
		f.err = errors.New("provided unnamed constant")
		return
	}
	if _, ok := f.declarations[name]; ok {
		f.err = fmt.Errorf("'%s' already declared within function(%s) definition", f.def.tp.Name(false, ""), name)
		return
	}
	tp, ok := types.GetBuiltInType(kind.BuiltInName())
	if !ok {
		f.err = fmt.Errorf("invalid constant value kind: %v", kind)
		return
	}
	if kind == types.KindString {
		value = quoteAll(value)
	}
	f.P("const ", name, " ", kind.BuiltInName(), " = ", value)
	f.declarations[name] = types.Declaration{Name: name, Type: tp, Constant: true}
	return
}

// Ret adds the return statements.
func (f *funcContent) Ret(args ...interface{}) {
	var res []string
	var indents uint8
	for i := range args {
		switch ot := args[i].(type) {
		case types.Declaration:
			res = append(res, ot.Name)
		case Indent:
			indents += ot.size
		case nil:
			res = append(res, "nil")
		case FuncExecution:
			res = append(res, ot.raw)
		default:
			res = append(res, fmt.Sprint(ot))
		}
	}

	var sb strings.Builder
	for i := uint8(0); i < indents; i++ {
		sb.WriteRune('\t')
	}
	sb.WriteString("return ")
	sb.WriteString(strings.Join(res, ","))
	f.P(sb.String())
}

// Q surrounds given input with the quotes if it
func (f *funcContent) Q(input string) string {
	return quoteAll(input)
}

// P writes a new line of the function content.
func (f *funcContent) P(args ...interface{}) {
	var sb strings.Builder

	var result []interface{}
	for _, arg := range args {
		if i, ok := arg.(Indent); ok {
			for j := uint8(0); j < i.size; j++ {
				sb.WriteRune('\t')
			}
			continue
		}
		result = append(result, arg)
	}
	if sb.Len() != 0 {
		result = append([]interface{}{sb.String()}, result...)
	}
	f.content.P(result...)
}

// Ind creates an ident in the content paragraph.
func (f *funcContent) Ind(size ...uint) Indent {
	i := Indent{size: 1}
	if len(size) == 1 {
		i.size = uint8(size[0])
	}
	return i
}

func (f *funcContent) getDeclaration(name string) (types.Declaration, bool) {
	d, ok := f.declarations[name]
	if ok {
		return d, true
	}
	d, ok = f.pkg.tp.Declarations[name]
	if ok {
		return d, true
	}
	return types.Declaration{}, false
}

func (f *funcContent) StructVal(st *types.Struct, si func(s *structValuer)) {

}

// ExecOut is the result of executing the function.
type ExecOut []types.Declaration

type Executor interface {
	Q(string) string
	Arg(arg interface{})
	ArgAt(index int, arg interface{})
	ArgsByDeclType()
	ArgsByExecOut(execOut ExecOut)
	Out(name string)
	OutAt(index int, name string)
}

var _ Executor = (*funcExecutor)(nil)

type funcExecutor struct {
	ft                  *types.Function
	fc                  *funcContent
	args                []string
	outs                []types.Declaration
	nextArg, currentOut int
	err                 error
}

// Q wraps input string with quotation marks.
func (f *funcExecutor) Q(input string) string {
	return quoteAll(input)
}

// Arg sets up next argument.
func (f *funcExecutor) Arg(arg interface{}) {
	if !f.arg(f.nextArg, arg) {
		return
	}
	f.nextArg++
}

// ArgsDone used as a for loop flag for the arguments. Returns true if all arguments are set using the 'Arg' method.
func (f *funcExecutor) ArgsDone() bool {
	return f.nextArg > len(f.args)
}

// ArgAt sets the argument value at the index 'index'.
func (f *funcExecutor) ArgAt(index int, arg interface{}) {
	f.arg(index, arg)
}

func (f *funcExecutor) arg(index int, arg interface{}) bool {
	if index >= len(f.args)-1 {
		f.err = errors.New("arg index out of range")
		return false
	}
	switch at := arg.(type) {
	case types.Declaration:
		f.args[index] = at.Name
	case nil:
		f.args[index] = "nil"
	default:
		f.args[index] = fmt.Sprint(arg)
	}
	return true
}

// ArgsByDeclType matches current function and package content declared types to given argument types.
func (f *funcExecutor) ArgsByDeclType() {
	if f.err != nil {
		return
	}
	for i, in := range f.ft.In {
		for name, decl := range f.fc.declarations {
			if in.Type.Equal(decl.Type) {
				f.args[i] = name
			}
		}
	}
}

// ArgsByExecOut
func (f *funcExecutor) ArgsByExecOut(execOut ExecOut) {

}

// Out sets up the next out result name.
func (f *funcExecutor) Out(name string) {
	if f.err != nil {
		return
	}
	if !f.out(f.currentOut, name) {
		return
	}
	f.currentOut++
}

// OutAt sets up the name of the output at index.
func (f *funcExecutor) OutAt(index int, name string) {
	if f.err != nil {
		return
	}
	f.out(index, name)
}

func (f *funcExecutor) out(index int, name string) bool {
	if index >= len(f.outs) {
		f.err = errors.New("out result index out of range")
		return false
	}
	f.outs[index] = types.Declaration{Name: name, Type: f.ft.Out[index].Type}
	return true
}

func newExecutor(fc *funcContent, f *types.Function) *funcExecutor {
	return &funcExecutor{
		fc:   fc,
		ft:   f,
		args: make([]string, len(f.In)),
		outs: make([]types.Declaration, len(f.Out)),
	}
}

// FuncExecution is the function execution generator.
type FuncExecution struct {
	raw string
	out ExecOut
}

func (f FuncExecution) String() string {
	return f.raw
}

func (f FuncExecution) Outs() ExecOut {
	return f.out
}

// FuncParams is a slice wrapper over function or method parameters.
type FuncParams []types.FuncParam

// WriteTo writes parameter names with respect to provided packageContext. If the names are not defined, but it is required.
// i.e.: for named output it is
func (f FuncParams) WriteTo(w io.Writer, packageContext string, namesRequired bool) (err error) {
	for i, param := range f {
		if param.Name == "" && namesRequired {
			typeName := []rune(param.Type.Name(false, ""))
			for _, rn := range typeName {
				if unicode.IsLetter(rn) {
					param.Name = fmt.Sprintf("_param%0d_%s", i, string([]rune{rn}))
					break
				}
			}
		}
		if param.Name != "" {
			if _, err = fmt.Fprintf(w, "%s ", param.Name); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprint(w, param.Type.Name(true, packageContext)); err != nil {
			return err
		}
		if i != len(f)-1 {
			if _, err = fmt.Fprint(w, ", "); err != nil {
				return err
			}
		}
	}
	return nil
}

// MergeByType merges all func params by given type.
func (f FuncParams) MergeByType() FuncParams {
	type funcParam struct {
		index int
		tp    types.Type
		name  string
	}
	mt := map[string]*funcParam{}
	var i int
	for _, param := range f {
		if pm, ok := mt[param.Type.FullName()]; ok {
			if pm.name == "" && param.Name != "" {
				pm.name = param.Name
			}
			continue
		}
		mt[param.Type.FullName()] = &funcParam{index: i, tp: param.Type, name: param.Name}
		i++
	}
	params := make(FuncParams, len(mt))
	for _, v := range mt {
		params[v.index] = types.FuncParam{Name: v.name, Type: v.tp}
	}
	return params
}

func quoteAll(input string) string {
	var sb strings.Builder
	sb.WriteRune('"')
	sb.WriteString(input)
	sb.WriteRune('"')
	return sb.String()
}
