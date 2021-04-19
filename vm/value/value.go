package value

import (
	"encoding/gob"
	"fmt"
	"strconv"
)

//go:generate stringer -type=Type -trimprefix=Type
type Type uint8

const (
	TypeInt Type = iota
	TypeFloat
	TypeBoolean
	TypeChar
	TypeNil
)

type Value interface {
	Type() Type
	String() string
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
)

func (*Int) Type() Type     { return TypeInt }
func (*Float) Type() Type   { return TypeFloat }
func (*Boolean) Type() Type { return TypeBoolean }
func (*Char) Type() Type    { return TypeChar }
func (*Nil) Type() Type     { return TypeNil }

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

func Equal(x, y Value) bool {
	switch a := x.(type) {
	case *Int:
		if b, ok := y.(*Int); ok {
			return a.Val == b.Val
		}
	case *Float:
		if b, ok := y.(*Float); ok {
			return a.Val == b.Val
		}
	case *Boolean:
		if b, ok := y.(*Boolean); ok {
			return a.Val == b.Val
		}
	case *Char:
		if b, ok := y.(*Char); ok {
			return a.Val == b.Val
		}
	case *Nil:
		_, ok := y.(*Nil)
		return ok
	}
	return false
}

func init() {
	gob.RegisterName("sometimes/vm/value.Int", &Int{})
	gob.RegisterName("sometimes/vm/value.Float", &Float{})
	gob.RegisterName("sometimes/vm/value.Boolean", &Boolean{})
	gob.RegisterName("sometimes/vm/value.Char", &Char{})
	gob.RegisterName("sometimes/vm/value.Nil", &Nil{})
}
