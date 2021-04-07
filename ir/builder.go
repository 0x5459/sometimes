package ir

import (
	"encoding/gob"
	"errors"
	"io"
	"sometimes/ir/value"
)

var ErrLabelNotExist = errors.New("label not exist")

type rawInstruction struct {
	Op   Op
	Args []Ptr
}

type Builder struct {
	rawInstructions []rawInstruction
	labels          map[string]Ptr
	consts          []value.Value
}

func NewBuilder() *Builder {
	return &Builder{}
}

func LoadBuilderFromBinary(r io.Reader) (b *Builder, err error) {
	dec := gob.NewDecoder(r)
	err = dec.Decode(b)
	return
}

func (b *Builder) Label(label string) {
	idx := len(b.rawInstructions)
	b.labels[label] = idx
}

func (b *Builder) GetLabelAddr(label string) (addr Ptr, exist bool) {
	addr, exist = b.labels[label]
	return
}

func (b *Builder) InsertJump(op Op, label string) error {
	addr, ok := b.labels[label]
	if !ok {
		return ErrLabelNotExist
	}
	b.rawInstructions = append(b.rawInstructions,
		rawInstruction{Op: op, Args: []Ptr{addr}})
	return nil
}

func (b *Builder) Insert(op Op, args ...value.Value) {
	argAddrs := make([]Ptr, len(args))
	for i, arg := range args {
		argAddrs[i] = b.insertConst(arg)
	}
	b.rawInstructions = append(b.rawInstructions,
		rawInstruction{Op: op, Args: argAddrs})
}

func (b *Builder) insertConst(v value.Value) (addr Ptr) {
	for i, c := range b.consts {
		if c == v {
			return i
		}
	}
	b.consts = append(b.consts, v)
	return len(b.consts) - 1
}

func (b *Builder) Output(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(b)
}

func (b *Builder) ToProgram() *Program {
	// return &Program{
	// 	Instructions: ,
	// }
	return nil
}
