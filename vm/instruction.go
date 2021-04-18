package vm

type Ptr = int
type DataID = int

//go:generate stringer -type=Op -trimprefix=Op
type Op uint8

// ops
const (
	op_arith_start Op = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNeg // negate
	op_arith_end

	op_logic_start
	OpEq  // equal
	OpNE  // not equal
	OpGT  // greater than
	OpLT  // less than
	OpGTE // greater than equal
	OpLTE // less than equal
	OpNot
	OpAnd
	OpOr
	op_logic_end

	OpJmp // jump
	OpJF  // jump if false

	OpCall
	OpRet // return

	OpSwap // Swap the top two items on the stack.

	OpPush

	OpLoad // Push a copy of the local with the given offset on to the stack
	OpStore
)

// Instruction is one instruction executed by the vm
type Instruction interface {
	Op() Op
}

type (
	InstrAdd struct{}

	InstrSub struct{}
	InstrMul struct{}
	InstrDiv struct{}
	InstrMod struct{}
	InstrNeg struct{}

	InstrEq  struct{}
	InstrNE  struct{}
	InstrGT  struct{}
	InstrLT  struct{}
	InstrGTE struct{}
	InstrLTE struct{}

	InstrNot struct{}
	InstrAnd struct{}
	InstrOr  struct{}

	InstrJmp struct {
		Addr Ptr
	}
	InstrJF struct {
		Addr Ptr
	}

	InstrCall struct {
		Addr  Ptr
		Arity int // args length
	}
	InstrRet struct{}

	InstrPush struct {
		DataID DataID
	}

	InstrLoad struct {
		Offset int
	}

	InstrStore struct {
		Offset int
	}
)

func (*InstrAdd) Op() Op   { return OpAdd }
func (*InstrSub) Op() Op   { return OpSub }
func (*InstrMul) Op() Op   { return OpMul }
func (*InstrDiv) Op() Op   { return OpDiv }
func (*InstrMod) Op() Op   { return OpMod }
func (*InstrNeg) Op() Op   { return OpNeg }
func (*InstrEq) Op() Op    { return OpEq }
func (*InstrNE) Op() Op    { return OpNE }
func (*InstrGT) Op() Op    { return OpGT }
func (*InstrLT) Op() Op    { return OpLT }
func (*InstrGTE) Op() Op   { return OpGTE }
func (*InstrLTE) Op() Op   { return OpLTE }
func (*InstrNot) Op() Op   { return OpNot }
func (*InstrAnd) Op() Op   { return OpAnd }
func (*InstrOr) Op() Op    { return OpOr }
func (*InstrJmp) Op() Op   { return OpJmp }
func (*InstrJF) Op() Op    { return OpJF }
func (*InstrCall) Op() Op  { return OpCall }
func (*InstrRet) Op() Op   { return OpRet }
func (*InstrPush) Op() Op  { return OpPush }
func (*InstrLoad) Op() Op  { return OpLoad }
func (*InstrStore) Op() Op { return OpStore }
