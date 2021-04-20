package assembly

import (
	"errors"
	"sometimes/hir"
	"strconv"
	"strings"
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

func (ap *AssemblyProgram) EmitPush(val hir.Value) {
	dataID := ap.Consts.insertConst(val)
	ap.Emit(&AssemblyInstrPush{DataID: dataID})
}

func (ap *AssemblyProgram) Emit(assemblyInstr AssemblyInstruction) {
	ap.Instructions = append(ap.Instructions, assemblyInstr)
}

func (ap *AssemblyProgram) String() string {
	/**
	type AssemblyProgram struct {
		Labels       map[string]Ptr
		Consts       *Consts
		Instructions []AssemblyInstruction
	}
	*/
	var sb strings.Builder

	for dataID, data := range ap.Consts.Inner {
		sb.WriteRune('@')
		sb.WriteString(strconv.Itoa(int(dataID)))
		sb.WriteString(" = ")
		sb.WriteString(data.String())
		sb.WriteRune('\n')
	}

	sb.WriteRune('\n')
	labels := make(map[Ptr]string, len(ap.Labels))
	for label, addr := range ap.Labels {
		labels[addr] = label
	}
	for addr, instr := range ap.Instructions {
		if label, ok := labels[addr]; ok {
			sb.WriteString(label)
			sb.WriteString(":\n")
		}
		sb.WriteString(instr.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}

type Consts struct {
	dataID DataID
	Inner  map[DataID]hir.Value
}

func NewConsts() *Consts {
	return &Consts{
		dataID: 0,
		Inner:  make(map[DataID]hir.Value),
	}
}

func (c *Consts) fetchDataID() DataID {
	id := c.dataID
	atomic.AddUint32(&c.dataID, 1)
	return id
}

func (c *Consts) insertConst(val hir.Value) DataID {
	for dataID, data := range c.Inner {
		if hir.ValueEqual(data, val) {
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

func (c *Consts) GetConst(dataID DataID) hir.Value {
	return c.Inner[dataID]
}
