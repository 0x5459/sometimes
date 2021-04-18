package vm

import (
	"errors"
	"fmt"
	"sometimes/vm/value"
)

type VM struct {
	operandStack *OperandStack
	frames       *FrameStack
	pc           Ptr
	program      *Program
}

func New(program *Program, operandStackCap, frameStackCap int) *VM {
	frames := NewFrameStack(frameStackCap)
	// todo temp
	frames.Push(&Frame{
		Local:   NewLocal(100),
		RetAddr: 0,
	})
	return &VM{
		operandStack: NewOperandStack(operandStackCap),
		frames:       frames,
		pc:           program.Entry(),
		program:      program,
	}
}

func (vm *VM) fetch() (ins Instruction, exist bool) {
	if ins, exist = vm.program.FetchInstruction(vm.pc); exist {
		vm.pc++
	}
	return
}

func (vm *VM) Execute() {
	for ins, exist := vm.fetch(); exist && !vm.frames.IsEmpty(); ins, exist = vm.fetch() {
		fmt.Println(ins.Op(), ins)
		switch instr := ins.(type) {
		case *InstrPush:
			v, _ := vm.program.GetConst(instr.DataID)
			vm.operandStack.Push(v)
		case *InstrJmp:
			vm.pc = instr.Addr
		case *InstrJF:
			b := vm.operandStack.Pop().(*value.Boolean)
			if !b.Val {
				vm.pc = instr.Addr
			}
		case *InstrCall:
			vm.frames.Push(&Frame{
				Local:   NewLocal(instr.Arity),
				RetAddr: vm.pc,
			})
			// jump to function
			vm.pc = instr.Addr
		case *InstrRet:
			if vm.frames.IsEmpty() {
				panic(errors.New("can't return from top-level"))
			}
			frame := vm.frames.Pop()
			// jump to caller
			vm.pc = frame.RetAddr
		case *InstrLoad:
			v := vm.frames.Top().Local.Load(instr.Offset)
			vm.operandStack.Push(v)
		case *InstrStore:
			v := vm.operandStack.Pop()
			vm.frames.Top().Local.Store(instr.Offset, v)
		case BinaryArithInstruction:
			lhs := vm.operandStack.Pop()
			rhs := vm.operandStack.Pop()
			vm.operandStack.Push(arith(instr, lhs, rhs))
		case UnaryArithInstruction:
			x := vm.operandStack.Pop()
			vm.operandStack.Push(arith(instr, x, &value.Nil{}))
		case BinaryLogicInstruction:
			lhs := vm.operandStack.Pop()
			rhs := vm.operandStack.Pop()
			vm.operandStack.Push(&value.Boolean{Val: logic(instr, lhs, rhs)})
		case UnaryLogicInstruction:
			x := vm.operandStack.Pop()
			vm.operandStack.Push(&value.Boolean{Val: logic(instr, x, &value.Nil{})})
		}
	}

}

func (vm *VM) PrintOperandStack() {
	fmt.Print("[")
	for i := 0; i < vm.operandStack.top; i++ {
		fmt.Print(vm.operandStack.inner[i].String())
		fmt.Print(", ")
	}
	fmt.Println("]")
}
