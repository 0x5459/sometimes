package vm

import (
	"fmt"
	"sometimes/ir"
	"sometimes/ir/value"
)

type VM struct {
	operandStack *OperandStack
	frames       *FrameStack
	pc           Ptr
	program      *ir.Program
}

func New(program *ir.Program, operandStackCap, frameStackCap int) *VM {
	return &VM{
		operandStack: NewOperandStack(operandStackCap),
		frames:       NewFrameStack(frameStackCap),
		pc:           program.Entry(),
		program:      program,
	}
}

func (vm *VM) fetch() (ins ir.Instruction, exist bool) {
	if ins, exist = vm.program.FetchInstruction(vm.pc); exist {
		vm.pc++
	}
	return
}

func (vm *VM) execute() {
	for ins, exist := vm.fetch(); exist && !vm.frames.IsEmpty(); ins, exist = vm.fetch() {
		switch instr := ins.(type) {
		case *ir.Push:
			vm.operandStack.Push(instr.Val)
		case *ir.Jmp:
			vm.pc = instr.Addr
		case *ir.JT:
			b := vm.operandStack.Pop().(*value.Boolean)
			if b.Val {
				vm.pc = instr.Addr
			}
		case *ir.JF:
			b := vm.operandStack.Pop().(*value.Boolean)
			if !b.Val {
				vm.pc = instr.Addr
			}
		case *ir.Call:
			vm.frames.Push(Frame{RetAddr: vm.pc})
			vm.pc = instr.Addr
		case *ir.Ret:
			retVal := vm.operandStack.Pop()
			vm.operandStack.
		case *ir.Eq:
		case *ir.NE:
		case *ir.GT:
		case *ir.LT:
		case *ir.GTE:
		case *ir.LTE:
		case *ir.Add:
		case *ir.Sub:
		case *ir.Mul:
		case *ir.Div:
		case *ir.Mod:
			lhs := vm.operandStack.Pop()
			rhs := vm.operandStack.Pop()
			vm.operandStack.Push(binaryNumberOp(instr.Op(), lhs, rhs))
		case *ir.Neg:
			v := vm.operandStack.Pop().(*value.Number)
			vm.operandStack.Push(&value.Number{Val: -v.Val})
		case *ir.Not:
			v := vm.operandStack.Pop().(*value.Boolean)
			vm.operandStack.Push(&value.Boolean{Val: !v.Val})
		case *ir.And:
			lhs := vm.operandStack.Pop().(*value.Boolean)
			rhs := vm.operandStack.Pop().(*value.Boolean)
			vm.operandStack.Push(&value.Boolean{Val: lhs.Val && rhs.Val})
		case *ir.Or:
			lhs := vm.operandStack.Pop().(*value.Boolean)
			rhs := vm.operandStack.Pop().(*value.Boolean)
			vm.operandStack.Push(&value.Boolean{Val: lhs.Val || rhs.Val})
		}
	}
}

func binaryNumberOp(op ir.Op, lhs, rhs value.Value) value.Value {
	if lhs.Type() != rhs.Type() || lhs.Type() != value.TypeNumber {
		panic(fmt.Errorf("unsupported operand type for `%s`: lhs: `%s` rhs: `%s`",
			op.String(), lhs.Type().String(), rhs.Type().String()))
	}

	x, y := lhs.(*value.Number), lhs.(*value.Number)

	switch op {
	case ir.OpEq:
		return &value.Boolean{Val: x.Val == y.Val}
	case ir.OpNE:
		return &value.Boolean{Val: x.Val != y.Val}
	case ir.OpGT:
		return &value.Boolean{Val: x.Val > y.Val}
	case ir.OpLT:
		return &value.Boolean{Val: x.Val < y.Val}
	case ir.OpGTE:
		return &value.Boolean{Val: x.Val >= y.Val}
	case ir.OpLTE:
		return &value.Boolean{Val: x.Val <= y.Val}
	case ir.OpAdd:
		return &value.Number{Val: x.Val + y.Val}
	case ir.OpSub:
		return &value.Number{Val: x.Val - y.Val}
	case ir.OpMul:
		return &value.Number{Val: x.Val * x.Val}
	case ir.OpDiv:
		return &value.Number{Val: x.Val / y.Val}
	case ir.OpMod:
		return &value.Number{Val: x.Val % y.Val}
	}
	return &value.Boolean{Val: false}
}
