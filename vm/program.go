package vm

import (
	"fmt"
	"sometimes/vm/assembly"
	"sometimes/vm/value"
)

type Program struct {
	instructions []Instruction
	consts       []value.Value
	entry        Ptr
}

func NewProgram(asm *assembly.AssemblyProgram) *Program {
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
			instrs[i] = &InstrGT{}
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
			instrs[i] = &InstrCall{Addr: getAsmLabelAddr(asm, asmInstr.Label)}
		case *assembly.AssemblyInstrRet:
			instrs[i] = &InstrRet{}
		case *assembly.AssemblyInstrPush:
			instrs[i] = &InstrPush{int(asmInstr.DataID)}
		case *assembly.AssemblyInstrLoad:
			instrs[i] = &InstrLoad{Offset: asmInstr.Offset}
		case *assembly.AssemblyInstrStore:
			instrs[i] = &InstrStore{Offset: asmInstr.Offset}
		}
	}

	consts := make([]value.Value, asm.Consts.MaxDataID())
	for id, v := range asm.Consts.Inner {
		consts[id] = v
	}
	return &Program{
		instructions: instrs,
		consts:       consts,
		entry:        0,
	}
}

func getAsmLabelAddr(asm *assembly.AssemblyProgram, label string) Ptr {
	if addr, ok := asm.GetLabelAddr(label); ok {
		return addr
	}
	panic(fmt.Errorf("label: %s not exist", label))
}

func (p *Program) Entry() Ptr {
	return p.entry
}

func (p *Program) FetchInstruction(addr Ptr) (ins Instruction, exist bool) {
	if len(p.instructions) > addr {
		ins, exist = p.instructions[addr], true
	} else {
		ins, exist = nil, false
	}
	return
}

func (p *Program) GetConst(dataId int) (val value.Value, exist bool) {
	if len(p.consts) > dataId {
		val, exist = p.consts[dataId], true
	} else {
		val, exist = nil, false
	}
	return
}
