package ir

import "sometimes/ir/value"

type Ptr = int

type Program struct {
	instructions []Instruction
	consts       []value.Value
	entry        Ptr
}

func (p *Program) Entry() Ptr {
	return p.entry
}

func (p *Program) FetchInstruction(addr Ptr) (ins Instruction, exist bool) {
	if len(p.instructions) < addr {
		ins, exist = p.instructions[addr], true
	} else {
		ins, exist = nil, false
	}
	return
}

func (p *Program) GetConst(addr Ptr) (val value.Value, exist bool) {
	if len(p.consts) < addr {
		val, exist = p.consts[addr], true
	} else {
		val, exist = nil, false
	}
	return
}
