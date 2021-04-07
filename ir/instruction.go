package ir

import "sometimes/ir/value"

type Op uint8

// ops
const (
	OpAdd Op = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNeg // negate

	OpEq  // equal
	OpNE  // not equal
	OpGT  // greater than
	OpLT  // less than
	OpGTE // greater than equal
	OpLTE // less than equal

	OpNot
	OpAnd
	OpOr

	OpJmp // jump
	OpJT  // jump if true
	OpJF  // jump if false

	OpCall
	OpRet // return

	OpSwap // Swap the top two items on the stack.

	OpPop
	OpPush
	OpPopLocal
)

var (
	ops = [...]string{
		OpAdd: "add",
		OpSub: "sub",
		OpMul: "mul",
		OpDiv: "div",
		OpMod: "mod",
		OpNeg: "neg",

		OpEq:  "eq",
		OpNE:  "ne",
		OpGT:  "gt",
		OpLT:  "lt",
		OpGTE: "gte",
		OpLTE: "lte",

		OpNot: "not",
		OpAnd: "and",
		OpOr:  "or",

		OpJmp: "jmp",
		OpJT:  "jt",
		OpJF:  "jf",

		OpCall: "call",
		OpRet:  "ret",

		OpPop:      "pop",
		OpPush:     "push",
		OpPopLocal: "pop_local",
	}
	opByOpStr = make(map[string]Op)
)

func (op Op) String() string {
	return ops[op]
}

func OpByOpStr(opStr string) (op Op, exist bool) {
	op, exist = opByOpStr[opStr]
	return
}

func init() {
	for op, opStr := range ops {
		opByOpStr[opStr] = Op(op)
	}
}

// Instruction is one instruction executed by the vm
type Instruction interface {
	Op() Op
}

type (
	Add struct{}

	Sub struct{}
	Mul struct{}
	Div struct{}
	Mod struct{}
	Neg struct{}

	Eq  struct{}
	NE  struct{}
	GT  struct{}
	LT  struct{}
	GTE struct{}
	LTE struct{}

	Not struct{}
	And struct{}
	Or  struct{}

	Jmp struct {
		Addr Ptr
	}
	JT struct {
		Addr Ptr
	}
	JF struct {
		Addr Ptr
	}

	Call struct{ Addr Ptr }
	Ret  struct{}

	Pop struct{}

	Push struct {
		Val value.Value
	}
	PopLocal struct {
		Val value.Value
	}
)

func (*Add) Op() Op  { return OpAdd }
func (*Sub) Op() Op  { return OpSub }
func (*Mul) Op() Op  { return OpMul }
func (*Div) Op() Op  { return OpDiv }
func (*Mod) Op() Op  { return OpMod }
func (*Neg) Op() Op  { return OpNeg }
func (*Eq) Op() Op   { return OpEq }
func (*NE) Op() Op   { return OpNE }
func (*GT) Op() Op   { return OpGT }
func (*LT) Op() Op   { return OpLT }
func (*GTE) Op() Op  { return OpGTE }
func (*LTE) Op() Op  { return OpLTE }
func (*Not) Op() Op  { return OpNot }
func (*And) Op() Op  { return OpAnd }
func (*Or) Op() Op   { return OpOr }
func (*Jmp) Op() Op  { return OpJmp }
func (*JT) Op() Op   { return OpJT }
func (*JF) Op() Op   { return OpJF }
func (*Call) Op() Op { return OpCall }
func (*Ret) Op() Op  { return OpRet }
func (*Pop) Op() Op  { return OpPop }
func (*Push) Op() Op { return OpPush }
