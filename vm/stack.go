package vm

import (
	"errors"
	"sometimes/vm/value"
)

var StackOverflow = errors.New("stack overflow")

type OperandStack struct {
	inner []value.Value
	top   Ptr // stack top pointer
	cap   int
}

// NewOperandStack returns a new Stack.
func NewOperandStack(cap int) *OperandStack {
	return &OperandStack{
		inner: make([]value.Value, 0, cap>>2),
		top:   0,
		cap:   cap,
	}
}

// Push appends an element to the back of a stack.
// or panic if the stack is full.
func (s *OperandStack) Push(v value.Value) {
	if s.top == s.cap {
		panic(StackOverflow)
	}
	if s.top >= len(s.inner) {
		s.inner = append(s.inner, v)
	} else {
		s.inner[s.top] = v
	}
	s.top++
}

// Pop removes the last element from a stack and returns it,
// or panic if it is empty.
func (s *OperandStack) Pop() value.Value {
	return s.PopN(1)
}

// PopN removes the last n elements from a stack and returns the last removed element,
// or panic if it is empty.
func (s *OperandStack) PopN(n int) value.Value {
	if s.top < n {
		panic(StackOverflow)
	}
	// 考虑缩容问题
	s.top -= n
	return s.inner[s.top]
}

func (s *OperandStack) IsEmpty() bool {
	return s.top == 0
}

func (s *OperandStack) Get(idx int) value.Value {
	if idx < 0 && idx >= s.top {
		panic(StackOverflow)
	}
	return s.inner[idx]
}

func (s *OperandStack) TopIdx() Ptr {
	return s.top
}

func (s *OperandStack) TopValue() value.Value {
	return s.Get(s.top - 1)
}
