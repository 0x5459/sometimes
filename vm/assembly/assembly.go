package assembly

import (
	"errors"
	"sometimes/vm/value"
	"sync/atomic"
)

type Ptr = int

var ErrLabelNotExist = errors.New("label not exist")

type AssemblyProgram struct {
	Labels       map[string]Ptr
	Consts       *Consts
	Instructions []AssemblyInstruction
}

func NewAssemblyProgram() *AssemblyProgram {
	return &AssemblyProgram{
		Labels: make(map[string]Ptr),
		Consts: NewConsts(),
	}
}

func (ap *AssemblyProgram) Label(label string) {
	idx := len(ap.Instructions)
	ap.Labels[label] = idx
}

func (ap *AssemblyProgram) GetLabelAddr(label string) (addr Ptr, exist bool) {
	addr, exist = ap.Labels[label]
	return
}

func (ap *AssemblyProgram) EmitJump(label string) error {
	if _, ok := ap.GetLabelAddr(label); !ok {
		return ErrLabelNotExist
	}
	ap.Emit(&AssemblyInstrJmp{
		Label: label,
	})
	return nil
}

func (ap *AssemblyProgram) EmitJF(label string) error {
	if _, ok := ap.GetLabelAddr(label); !ok {
		return ErrLabelNotExist
	}
	ap.Emit(&AssemblyInstrJF{
		Label: label,
	})
	return nil
}

func (ap *AssemblyProgram) EmitPush(val value.Value) {
	dataID := ap.Consts.insertConst(val)
	ap.Emit(&AssemblyInstrPush{DataID: dataID})
}

func (ap *AssemblyProgram) Emit(assemblyInstr AssemblyInstruction) {
	ap.Instructions = append(ap.Instructions, assemblyInstr)
}

type Consts struct {
	dataID DataID
	Inner  map[DataID]value.Value
}

func NewConsts() *Consts {
	return &Consts{
		dataID: 0,
		Inner:  make(map[DataID]value.Value),
	}
}

func (c *Consts) fetchDataID() DataID {
	id := c.dataID
	atomic.AddUint32(&c.dataID, 1)
	return id
}

func (c *Consts) insertConst(val value.Value) DataID {
	for dataID, data := range c.Inner {
		if value.Equal(data, val) {
			return dataID
		}
	}
	dataID := c.fetchDataID()
	c.Inner[dataID] = val
	return dataID
}

func (c *Consts) MaxDataID() DataID {
	return c.dataID
}
