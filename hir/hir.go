package hir

import (
	"fmt"
	"strconv"
)

type Program struct {
	funcs         map[string]*ExprFunction
	entryFuncName string
	consts        map[string]Value
}

func (p *Program) FindConst(name string) (val Value, isExist bool) {
	val, isExist = p.consts[name]
	return
}

func (p *Program) FindFunc(name string) (f *ExprFunction, isExist bool) {
	f, isExist = p.funcs[name]
	return
}

func (p *Program) EntryFunc() *ExprFunction {
	return p.funcs[p.entryFuncName]
}

func (p *Program) Funcs() []*ExprFunction {
	funcs := make([]*ExprFunction, 0, len(p.funcs))
	for _, f := range p.funcs {
		funcs = append(funcs, f)
	}
	return funcs
}

func (p *Program) Consts() map[string]Value {
	return p.consts
}

type Builder struct {
	Funcs        []*ExprFunction
	EntryFuncIdx int
	consts       map[string]Value
}

func NewBuilder() *Builder {
	return &Builder{
		Funcs:        []*ExprFunction{},
		EntryFuncIdx: -1,
		consts:       make(map[string]Value),
	}
}

func (b *Builder) InsertConst(name string, val Value) {
	b.consts[name] = val
}

func (b *Builder) InsertFunc(f *ExprFunction, entryFunc bool) {
	if entryFunc {
		b.EntryFuncIdx = len(b.Funcs)
	}
	b.Funcs = append(b.Funcs, f)
}

func (b *Builder) Build() *Program {
	funcs := make(map[string]*ExprFunction, len(b.Funcs))

	for _, f := range b.Funcs {
		funcs[f.Func.Name] = f
	}

	var entryFuncName string

	if b.EntryFuncIdx == -1 {
		panic("empty program")
	}
	entryFuncName = (b.Funcs[b.EntryFuncIdx]).Func.Name
	return &Program{
		funcs:         funcs,
		entryFuncName: entryFuncName,
		consts:        b.consts,
	}
}

type FuncBuilder struct {
	funcName string
	funcBody []Expr
	args     []*Binding
}

func NewFuncBuilder(funcName string, args []*Binding) *FuncBuilder {
	return &FuncBuilder{
		funcName: funcName,
		funcBody: []Expr{},
		args:     args,
	}
}

func (b *FuncBuilder) Emit(e Expr) {
	b.funcBody = append(b.funcBody, e)
}

func (b *FuncBuilder) Build() *ExprFunction {
	return &ExprFunction{
		Func: &Function{
			Name: b.funcName,
			Body: &ExprBlock{
				Body: b.funcBody,
			},
			Args: b.args,
		},
	}
}

type Binding struct {
	Name string
}

func NewBinding(name string) *Binding {
	return &Binding{
		Name: name,
	}
}

type Function struct {
	Name string
	Body *ExprBlock
	Args []*Binding
}

type Value interface {
	isValue()
	String() string
}

type (
	ValueInt struct {
		Val int
	}
	ValueFloat struct {
		Val float64
	}
	ValueString struct {
		Val string
	}
	ValueBoolean struct {
		Val bool
	}
	ValueNil struct{}

	ValueFunc struct {
		FuncName  string
		MaxLoacls int
	}
)

func NewValueInt(v int) *ValueInt {
	return &ValueInt{Val: v}
}

func NewValueFloat(v float64) *ValueFloat {
	return &ValueFloat{Val: v}
}

func NewValueString(v string) *ValueString {
	return &ValueString{Val: v}
}

func NewValueBoolean(v bool) *ValueBoolean {
	return &ValueBoolean{Val: v}
}

func NewValueNil() *ValueNil {
	return &ValueNil{}
}

func (*ValueInt) isValue()     {}
func (*ValueFloat) isValue()   {}
func (*ValueString) isValue()  {}
func (*ValueBoolean) isValue() {}
func (*ValueNil) isValue()     {}
func (*ValueFunc) isValue()    {}

func (i *ValueInt) String() string {
	return strconv.Itoa(i.Val)
}
func (f *ValueFloat) String() string {
	return fmt.Sprintf("%f", f.Val)
}
func (s *ValueString) String() string {
	return s.Val
}
func (b *ValueBoolean) String() string {
	if b.Val {
		return "true"
	} else {
		return "false"
	}
}

func (*ValueNil) String() string {
	return "<nil>"
}

func (f *ValueFunc) String() string {
	return fmt.Sprintf("Func @%s", f.FuncName)
}

func ValueEqual(x, y Value) bool {
	switch a := x.(type) {
	case *ValueInt:
		if b, ok := y.(*ValueInt); ok {
			return a.Val == b.Val
		}
	case *ValueFloat:
		if b, ok := y.(*ValueFloat); ok {
			return a.Val == b.Val
		}
	case *ValueBoolean:
		if b, ok := y.(*ValueBoolean); ok {
			return a.Val == b.Val
		}
	case *ValueString:
		if b, ok := y.(*ValueString); ok {
			return a.Val == b.Val
		}
	case *ValueNil:
		_, ok := y.(*ValueNil)
		return ok
	case *ValueFunc:
		if b, ok := y.(*ValueFunc); ok {
			return a.FuncName == b.FuncName
		}
	}
	return false
}

type BinaryOp uint8

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpEq  // equal
	OpNE  // not equal
	OpGT  // greater than
	OpLT  // less than
	OpGTE // greater than equal
	OpLTE // less than equal
	OpAnd
	OpOr
)

type UnaryOp uint8

const (
	OpNeg UnaryOp = iota
	OpNot
)
