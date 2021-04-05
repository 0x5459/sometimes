package vm

// Pointer type is int alias.
type Ptr = int

type Frame struct {
	RetAddr Ptr // instruction Pointer
}

type VM struct {
	stack  VMStack
	frames FrameStack
}
