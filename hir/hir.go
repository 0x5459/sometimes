package hir

import (
	"fmt"
	"strconv"
)

type Program struct {
	EntryFunc *ExprFunction
	consts    map[string]Value
}

func (p *Program) FindConst(name string) (val Value, isExist bool) {
	val, isExist = p.consts[name]
	return
}

type FuncBuilder struct {
	maxLocals int
	function  *ExprFunction
	consts    map[string]Value
}

func NewFuncBuilder(funcName string, args []*Binding) *FuncBuilder {
	return &FuncBuilder{
		maxLocals: len(args),
		function: &ExprFunction{
			Func: &Function{
				Name: funcName,
				Body: &ExprBlock{},
			},
		},
	}
}

func (b *FuncBuilder) Emit(e Expr) {
	if _, ok := e.(*ExprVar); ok {
		b.maxLocals++
	}
	b.function.Func.Body.Body = append(b.function.Func.Body.Body, e)
}

func (b *FuncBuilder) InsertConst(name string, val Value) {
	b.consts[name] = val
}

func (b *FuncBuilder) Build() *Program {
	return &Program{EntryFunc: b.function, consts: b.consts}
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
