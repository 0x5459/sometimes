package vm

import (
	"encoding/gob"
	"fmt"
	"io"
	"sometimes/hir"
	"sometimes/vm/assembly"
	"sometimes/vm/value"
)

type Program struct {
	Instructions []Instruction
	Consts       []value.Value
	Entry        Ptr
}

func NewProgramFromBinary(r io.Reader) *Program {
	dec := gob.NewDecoder(r)
	var p Program
	if err := dec.Decode(&p); err != nil {
		panic(err)
	}
	return &p
}

func NewProgramFromAsm(asm *assembly.AssemblyProgram) *Program {
	instrs := make([]Instruction, len(asm.Instructions))
	for i, assemblyInstruction := range asm.Instructions {
		switch asmInstr := assemblyInstruction.(type) {
		case *assembly.AssemblyInstrAdd:
			instrs[i] = &InstrAdd{}
		case *assembly.AssemblyInstrSub:
			instrs[i] = &InstrSub{}
		case *assembly.AssemblyInstrMul:
			instrs[i] = &InstrMul{}
		case *assembly.AssemblyInstrDiv:
			instrs[i] = &InstrDiv{}
		case *assembly.AssemblyInstrMod:
			instrs[i] = &InstrMod{}
		case *assembly.AssemblyInstrNeg:
			instrs[i] = &InstrNeg{}
		case *assembly.AssemblyInstrEq:
			instrs[i] = &InstrEq{}
		case *assembly.AssemblyInstrNE:
			instrs[i] = &InstrNE{}
		case *assembly.AssemblyInstrGT:
			instrs[i] = &InstrGT{}
		case *assembly.AssemblyInstrLT:
			instrs[i] = &InstrLT{}
		case *assembly.AssemblyInstrGTE:
			instrs[i] = &InstrGTE{}
		case *assembly.AssemblyInstrLTE:
			instrs[i] = &InstrLTE{}
		case *assembly.AssemblyInstrNot:
			instrs[i] = &InstrNot{}
		case *assembly.AssemblyInstrAnd:
			instrs[i] = &InstrAnd{}
		case *assembly.AssemblyInstrOr:
			instrs[i] = &InstrOr{}
		case *assembly.AssemblyInstrJmp:
			instrs[i] = &InstrJmp{Addr: getAsmLabelAddr(asm, asmInstr.Label)}
		case *assembly.AssemblyInstrJF:
			instrs[i] = &InstrJF{Addr: getAsmLabelAddr(asm, asmInstr.Label)}
		case *assembly.AssemblyInstrCall:
			instrs[i] = &InstrCall{}
		case *assembly.AssemblyInstrRet:
			instrs[i] = &InstrRet{}
		case *assembly.AssemblyInstrPush:
			instrs[i] = &InstrPush{int(asmInstr.DataID)}
		case *assembly.AssemblyInstrLoad:
			instrs[i] = &InstrLoad{Offset: asmInstr.Offset}
		case *assembly.AssemblyInstrStore:
			instrs[i] = &InstrStore{Offset: asmInstr.Offset}
		case *assembly.AssemblyInstrDup:
			instrs[i] = &InstrDup{}
		case *assembly.AssemblyInstrLoadFromPtr:
			instrs[i] = &InstrLoadFromPtr{}
		case *assembly.AssemblyInstrLoadPtr:
			instrs[i] = &InstrLoadPtr{
				Offset:  asmInstr.Offset,
				IsLocal: asmInstr.IsLocal,
			}
		case *assembly.AssemblyInstrStoreToPtr:
			instrs[i] = &InstrStoreToPtr{}
		case *assembly.AssemblyInstrPrint:
			instrs[i] = &InstrPrint{ArgLen: asmInstr.ArgLen}
		}
	}

	consts := make([]value.Value, asm.Consts.MaxDataID())
	for id, v := range asm.Consts.Inner {
		if f, ok := v.(*hir.ValueFunc); ok {
			consts[id] = &value.Func{
				Addr:      getAsmLabelAddr(asm, f.FuncName),
				MaxLocals: f.MaxLoacls,
			}
		} else {
			consts[id] = hirValueToVmValue(v)
		}
	}
	return &Program{
		Instructions: instrs,
		Consts:       consts,
		Entry:        0,
	}
}

func getAsmLabelAddr(asm *assembly.AssemblyProgram, label string) Ptr {
	if addr, ok := asm.GetLabelAddr(label); ok {
		return addr
	}
	panic(fmt.Errorf("label: %s not exist", label))
}

func (p *Program) FetchInstruction(addr Ptr) (ins Instruction, exist bool) {
	if len(p.Instructions) > addr {
		ins, exist = p.Instructions[addr], true
	} else {
		ins, exist = nil, false
	}
	return
}

func (p *Program) GetConst(dataId int) (val value.Value, exist bool) {
	if len(p.Consts) > dataId {
		val, exist = p.Consts[dataId], true
	} else {
		val, exist = nil, false
	}
	return
}

func (p *Program) WriteBinary(w io.Writer) {
	enc := gob.NewEncoder(w)
	err := enc.Encode(p)
	if err != nil {
		panic(err)
	}
}

func hirValueToVmValue(hirVal hir.Value) value.Value {
	switch hv := hirVal.(type) {
	case *hir.ValueInt:
		return &value.Int{Val: hv.Val}
	case *hir.ValueFloat:
		return &value.Float{Val: hv.Val}
	case *hir.ValueBoolean:
		return &value.Boolean{Val: hv.Val}
	case *hir.ValueString:
		// todo vm value unsupport string
		return &value.Char{Val: rune(hv.Val[0])}
	case *hir.ValueNil:
		return &value.Nil{}
	}
	return &value.Nil{}
}
