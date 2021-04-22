package hir

import (
	"fmt"
	"strconv"
)

type Ptr = int

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
