package hir

import (
	"fmt"
	"strconv"
)

type Program struct {
	Code   []Expr
	consts map[string]Value
}

func (p *Program) FindConst(name string) (val Value, isExist bool) {
	val, isExist = p.consts[name]
	return
}

type Builder struct {
	code   []Expr
	consts map[string]Value
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Emit(e Expr) {
	b.code = append(b.code, e)
}

func (b *Builder) InsertConst(name string, val Value) {
	b.consts[name] = val
}

func (b *Builder) Build() *Program {
	return &Program{Code: b.code, consts: b.consts}
}

type Binding struct {
	Name string
	// Depth < 0 if dealing with a global.
	Depth, FuncDepth int
}

func NewBinding(name string, depth, funcDepth int) *Binding {
	return &Binding{
		Name:      name,
		Depth:     depth,
		FuncDepth: funcDepth,
	}
}

func NewGlobalBinding(name string) *Binding {
	return &Binding{
		Name:      name,
		Depth:     -1,
		FuncDepth: 0,
	}
}

func (b *Binding) IsUpvalue() bool {
	return b.Depth > b.FuncDepth
}

func (b *Binding) IsConst() bool {
	return b.Depth < 0
}

func (b *Binding) UpvalueDepth() (depth int, isUpvalue bool) {
	if b.IsUpvalue() {
		depth = b.Depth - b.FuncDepth
		isUpvalue = true
	}
	return
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
