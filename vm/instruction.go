package vm

import "encoding/gob"

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

	OpPush
	OpDup

	OpLoad  // Push a copy of the local with the given offset on to the stack
	OpStore // Store value of stack top to local with the given offset

	OpLoadPtr
	OpLoadFromPtr
	OpStoreToPtr
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

	InstrCall struct{}

	InstrRet struct{}

	InstrPush struct {
		DataID DataID
	}
	InstrDup struct{}

	InstrLoad struct {
		Offset int
	}

	InstrLoadPtr struct {
		Offset  int
		IsLocal bool
	}

	InstrStore struct {
		Offset int
	}

	InstrLoadFromPtr struct{}
	InstrStoreToPtr  struct{}
)

func (*InstrAdd) Op() Op         { return OpAdd }
func (*InstrSub) Op() Op         { return OpSub }
func (*InstrMul) Op() Op         { return OpMul }
func (*InstrDiv) Op() Op         { return OpDiv }
func (*InstrMod) Op() Op         { return OpMod }
func (*InstrNeg) Op() Op         { return OpNeg }
func (*InstrEq) Op() Op          { return OpEq }
func (*InstrNE) Op() Op          { return OpNE }
func (*InstrGT) Op() Op          { return OpGT }
func (*InstrLT) Op() Op          { return OpLT }
func (*InstrGTE) Op() Op         { return OpGTE }
func (*InstrLTE) Op() Op         { return OpLTE }
func (*InstrNot) Op() Op         { return OpNot }
func (*InstrAnd) Op() Op         { return OpAnd }
func (*InstrOr) Op() Op          { return OpOr }
func (*InstrJmp) Op() Op         { return OpJmp }
func (*InstrJF) Op() Op          { return OpJF }
func (*InstrCall) Op() Op        { return OpCall }
func (*InstrRet) Op() Op         { return OpRet }
func (*InstrPush) Op() Op        { return OpPush }
func (*InstrDup) Op() Op         { return OpDup }
func (*InstrLoad) Op() Op        { return OpLoad }
func (*InstrStore) Op() Op       { return OpStore }
func (*InstrLoadPtr) Op() Op     { return OpLoadPtr }
func (*InstrLoadFromPtr) Op() Op { return OpLoadFromPtr }
func (*InstrStoreToPtr) Op() Op  { return OpStoreToPtr }

func init() {
	gob.RegisterName("sometimes/vm.InstrAdd", &InstrAdd{})
	gob.RegisterName("sometimes/vm.InstrSub", &InstrSub{})
	gob.RegisterName("sometimes/vm.InstrMul", &InstrMul{})
	gob.RegisterName("sometimes/vm.InstrDiv", &InstrDiv{})
	gob.RegisterName("sometimes/vm.InstrMod", &InstrMod{})
	gob.RegisterName("sometimes/vm.InstrNeg", &InstrNeg{})
	gob.RegisterName("sometimes/vm.InstrEq", &InstrEq{})
	gob.RegisterName("sometimes/vm.InstrNE", &InstrNE{})
	gob.RegisterName("sometimes/vm.InstrGT", &InstrGT{})
	gob.RegisterName("sometimes/vm.InstrLT", &InstrLT{})
	gob.RegisterName("sometimes/vm.InstrGTE", &InstrGTE{})
	gob.RegisterName("sometimes/vm.InstrLTE", &InstrLTE{})
	gob.RegisterName("sometimes/vm.InstrNot", &InstrNot{})
	gob.RegisterName("sometimes/vm.InstrAnd", &InstrAnd{})
	gob.RegisterName("sometimes/vm.InstrOr", &InstrOr{})
	gob.RegisterName("sometimes/vm.InstrJmp", &InstrJmp{})
	gob.RegisterName("sometimes/vm.InstrJF", &InstrJF{})
	gob.RegisterName("sometimes/vm.InstrCall", &InstrCall{})
	gob.RegisterName("sometimes/vm.InstrRet", &InstrRet{})
	gob.RegisterName("sometimes/vm.InstrPush", &InstrPush{})
	gob.RegisterName("sometimes/vm.InstrDup", &InstrDup{})
	gob.RegisterName("sometimes/vm.InstrLoad", &InstrLoad{})
	gob.RegisterName("sometimes/vm.InstrStore", &InstrStore{})
	gob.RegisterName("sometimes/vm.InstrLoadPtr", &InstrLoadPtr{})
	gob.RegisterName("sometimes/vm.InstrLoadFromPtr", &InstrLoadFromPtr{})
	gob.RegisterName("sometimes/vm.InstrStoreToPtr", &InstrStoreToPtr{})
}
