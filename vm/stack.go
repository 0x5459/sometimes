// Code generated by go2go; DO NOT EDIT.


//line stack.go2:1
package vm

//line stack.go2:1
import (
//line stack.go2:1
 "errors"
//line stack.go2:1
 "sometimes/ir/value"
//line stack.go2:1
)

//line stack.go2:9
type Ptr = int

//line stack.go2:12
type Frame struct {
	RetAddr Ptr
}

var StackOverflow = errors.New("stack overflow")

//line stack.go2:19
type OperandStack = instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue

//line stack.go2:23
type FrameStack = instantiate୦୦Stack୦vm୮aFrame

//line stack.go2:32
func NewOperandStack(cap int) *instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue {
	return instantiate୦୦NewStack୦sometimes୮dir୮dvalue୮aValue(cap)
}

//line stack.go2:37
func NewFrameStack(cap int) *instantiate୦୦Stack୦vm୮aFrame {
	return instantiate୦୦NewStack୦vm୮aFrame(cap)
}

//line stack.go2:39
type instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue struct {
//line stack.go2:26
 inner []value.Value

//line stack.go2:27
 top int
			cap int
}

//line stack.go2:53
func (s *instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue,) Push(v value.Value,

//line stack.go2:53
) {
	if s.top == s.cap {
		panic(StackOverflow)
	}
	s.inner = append(s.inner, v)
	s.top++
}

//line stack.go2:63
func (s *instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue,) Pop() value.Value {
	return s.PopN(1)
}

//line stack.go2:69
func (s *instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue,) PopN(n int) value.Value {
	if s.top < n {
		panic(StackOverflow)
	}
	s.top -= n
	return s.inner[s.top]
}

func (s *instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue,) IsEmpty() bool {
	return s.top == 0
}

//line stack.go2:79
type instantiate୦୦Stack୦vm୮aFrame struct {
//line stack.go2:26
 inner []Frame

//line stack.go2:27
 top int
			cap int
}

//line stack.go2:53
func (s *instantiate୦୦Stack୦vm୮aFrame,) Push(v Frame,

//line stack.go2:53
) {
	if s.top == s.cap {
		panic(StackOverflow)
	}
	s.inner = append(s.inner, v)
	s.top++
}

//line stack.go2:63
func (s *instantiate୦୦Stack୦vm୮aFrame,) Pop() Frame {
	return s.PopN(1)
}

//line stack.go2:69
func (s *instantiate୦୦Stack୦vm୮aFrame,) PopN(n int) Frame {
	if s.top < n {
		panic(StackOverflow)
	}
	s.top -= n
	return s.inner[s.top]
}

func (s *instantiate୦୦Stack୦vm୮aFrame,) IsEmpty() bool {
	return s.top == 0
}
//line stack.go2:43
func instantiate୦୦NewStack୦sometimes୮dir୮dvalue୮aValue(cap int) *instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue {
	return &instantiate୦୦Stack୦sometimes୮dir୮dvalue୮aValue{
		inner: make([]value.Value, 0, cap>>2),
		top:   0,
		cap:   cap,
	}
}
//line stack.go2:43
func instantiate୦୦NewStack୦vm୮aFrame(cap int) *instantiate୦୦Stack୦vm୮aFrame {
	return &instantiate୦୦Stack୦vm୮aFrame{
		inner: make([]Frame, 0, cap>>2),
		top:   0,
		cap:   cap,
	}
}

//line stack.go2:49
type Importable୦ int

//line stack.go2:49
var _ = errors.As

//line stack.go2:49
type _ value.Boolean
