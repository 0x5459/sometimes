package vm

import "sometimes/vm/value"

type Local struct {
	locals []value.Value
}

func NewLocal(maxLen int) *Local {
	return &Local{locals: make([]value.Value, maxLen)}
}

func (l *Local) Load(idx int) value.Value {
	v := l.locals[idx]
	if v == nil {
		return &value.Nil{}
	}
	return v
}

func (l *Local) Store(idx int, v value.Value) {
	l.locals[idx] = v
}

type Frame struct {
	Local   *Local
	RetAddr Ptr
}

type frameNode struct {
	frame      *Frame
	prev, next *frameNode
}

// implemented by linked list
type FrameStack struct {
	head     *frameNode
	len, cap int
}

// NewFrameStack returns a new FrameStack.
func NewFrameStack(cap int) *FrameStack {
	return &FrameStack{
		head: &frameNode{},
		len:  0,
		cap:  cap,
	}
}

func (fs *FrameStack) IsEmpty() bool {
	return fs.head.next == nil
}

func (fs *FrameStack) Push(frame *Frame) {
	if fs.len >= fs.cap {
		panic(StackOverflow)
	}
	var tail *frameNode
	if fs.IsEmpty() {
		tail = fs.head
	} else {
		tail = fs.head.prev
	}

	tail.next = &frameNode{
		frame: frame,
		prev:  tail,
		next:  nil,
	}
	fs.head.prev = tail.next
	fs.len++
}

func (fs *FrameStack) Pop() *Frame {
	if fs.IsEmpty() {
		panic(StackOverflow)
	}
	tail := fs.head.prev
	tail.prev.next = nil
	fs.head.prev = tail.prev
	return tail.frame
}

func (fs *FrameStack) Top() *Frame {
	if fs.IsEmpty() {
		panic(StackOverflow)
	}
	return fs.head.prev.frame
}
