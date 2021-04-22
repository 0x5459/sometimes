package value

import (
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
)

type Ptr = int

//go:generate stringer -type=Type -trimprefix=Type
type Type uint8

const (
	TypeInt Type = iota
	TypeFloat
	TypeBoolean
	TypeChar
	TypeNil
	TypeFunc
	TypePointer
)

type Value interface {
	Type() Type
	String() string
	Clone() Value
}

type NumberValue interface {
	Value
	isNumber()
}

type (
	Int struct {
		Val int
	}
	Float struct {
		Val float64
	}
	Boolean struct {
		Val bool
	}
	Char struct {
		Val rune
	}
	Nil struct {
	}
	Func struct {
		Addr      Ptr
		MaxLocals int
	}

	Pointer struct {
		Addr    Ptr
		IsLocal bool
	}
)

func (*Int) Type() Type     { return TypeInt }
func (*Float) Type() Type   { return TypeFloat }
func (*Boolean) Type() Type { return TypeBoolean }
func (*Char) Type() Type    { return TypeChar }
func (*Nil) Type() Type     { return TypeNil }
func (*Func) Type() Type    { return TypeFunc }
func (*Pointer) Type() Type { return TypePointer }

func (x *Int) Clone() Value     { return &Int{Val: x.Val} }
func (x *Float) Clone() Value   { return &Float{Val: x.Val} }
func (x *Boolean) Clone() Value { return &Boolean{Val: x.Val} }
func (x *Char) Clone() Value    { return &Char{Val: x.Val} }
func (x *Nil) Clone() Value     { return &Nil{} }
func (x *Func) Clone() Value {
	return &Func{
		Addr:      x.Addr,
		MaxLocals: x.MaxLocals,
	}
}
func (x *Pointer) Clone() Value {
	return &Pointer{
		Addr:    x.Addr,
		IsLocal: x.IsLocal,
	}
}

func (*Int) isNumber()   {}
func (*Float) isNumber() {}

func (n *Int) String() string {
	return strconv.Itoa(n.Val)
}

func (n *Float) String() string {
	return fmt.Sprintf("%f", n.Val)
}

func (b *Boolean) String() string {
	if b.Val {
		return "true"
	} else {
		return "false"
	}
}

func (c *Char) String() string {
	return string(c.Val)
}

func (n *Nil) String() string {
	return "<nil>"
}

func (f *Func) String() string {
	return fmt.Sprintf("Func #%d", f.Addr)
}

func (p *Pointer) String() string {
	var sb strings.Builder
	if p.IsLocal {
		sb.WriteString("Local")
	}
	sb.WriteString("Ptr(")
	sb.WriteString(strconv.Itoa(p.Addr))
	sb.WriteRune(')')
	return sb.String()
}

func init() {
	gob.RegisterName("sometimes/vm/value.Int", &Int{})
	gob.RegisterName("sometimes/vm/value.Float", &Float{})
	gob.RegisterName("sometimes/vm/value.Boolean", &Boolean{})
	gob.RegisterName("sometimes/vm/value.Char", &Char{})
	gob.RegisterName("sometimes/vm/value.Nil", &Nil{})
	gob.RegisterName("sometimes/vm/value.Func", &Func{})
	gob.RegisterName("sometimes/vm/value.Pointer", &Pointer{})
}
