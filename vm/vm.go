package vm

import (
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
	return &VM{
		operandStack: NewOperandStack(operandStackCap),
		frames:       frames,
		pc:           program.Entry,
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
	for ins, exist := vm.fetch(); exist; ins, exist = vm.fetch() {
		// fmt.Printf("op: %s, pc:%d\n", ins.Op().String(), vm.pc)
		// vm.PrintOperandStack()
		switch instr := ins.(type) {
		case *InstrPrint:
			for i := 0; i < instr.ArgLen; i++ {
				v := vm.operandStack.Pop()
				fmt.Print(v.String() + " ")
			}
			fmt.Println()
		case *InstrPush:
			v, _ := vm.program.GetConst(instr.DataID)
			vm.operandStack.Push(v)
		case *InstrDup:
			v := vm.operandStack.TopValue()
			vm.operandStack.Push(v.Clone())
		case *InstrJmp:
			vm.pc = instr.Addr
		case *InstrJF:
			b := vm.operandStack.Pop().(*value.Boolean)
			if !b.Val {
				vm.pc = instr.Addr
			}
		case *InstrCall:
			f := vm.operandStack.Pop().(*value.Func)
			vm.frames.Push(&Frame{
				Local:   NewLocal(f.MaxLocals),
				RetAddr: vm.pc,
			})
			// jump to function
			vm.pc = f.Addr
		case *InstrRet:
			frame := vm.frames.Pop()
			if vm.frames.IsEmpty() {
				return
			}
			// jump to caller
			vm.pc = frame.RetAddr
		case *InstrLoad:
			v := vm.frames.Top().Local.Load(instr.Offset)
			vm.operandStack.Push(v)
		case *InstrStore:
			v := vm.operandStack.Pop()
			vm.frames.Top().Local.Store(instr.Offset, v)
		case *InstrLoadPtr:
			vm.operandStack.Push(&value.Pointer{
				Addr:    instr.Offset,
				IsLocal: instr.IsLocal,
			})
		case *InstrLoadFromPtr:
			ptr := vm.operandStack.Pop().(*value.Pointer)
			if !ptr.IsLocal {
				panic("unimplement")
			}
			v := vm.frames.Top().Local.Load(ptr.Addr)
			vm.operandStack.Push(v)
		case *InstrStoreToPtr:
			v := vm.operandStack.Pop()
			ptr := vm.operandStack.Pop().(*value.Pointer)
			if !ptr.IsLocal {
				panic("unimplement")
			}
			vm.frames.Top().Local.Store(ptr.Addr, v)
		case BinaryArithInstruction:
			rhs := vm.operandStack.Pop()
			lhs := vm.operandStack.Pop()
			vm.operandStack.Push(arith(instr, lhs, rhs))
		case UnaryArithInstruction:
			x := vm.operandStack.Pop()
			vm.operandStack.Push(arith(instr, x, &value.Nil{}))
		case BinaryLogicInstruction:
			rhs := vm.operandStack.Pop()
			lhs := vm.operandStack.Pop()
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
