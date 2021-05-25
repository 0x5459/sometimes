package assembly

import "fmt"

type DataID = uint32

type AssemblyInstruction interface {
	isAssemblyInstruction()
	String() string
}

type (
	AssemblyInstrAdd struct{}

	AssemblyInstrSub struct{}
	AssemblyInstrMul struct{}
	AssemblyInstrDiv struct{}
	AssemblyInstrMod struct{}
	AssemblyInstrNeg struct{}
	AssemblyInstrEq  struct{}
	AssemblyInstrNE  struct{}
	AssemblyInstrGT  struct{}
	AssemblyInstrLT  struct{}
	AssemblyInstrGTE struct{}
	AssemblyInstrLTE struct{}

	AssemblyInstrNot struct{}
	AssemblyInstrAnd struct{}
	AssemblyInstrOr  struct{}

	AssemblyInstrJmp struct {
		Label string
	}
	AssemblyInstrJF struct {
		Label string
	}

	AssemblyInstrCall struct{}

	AssemblyInstrRet struct{}

	AssemblyInstrPush struct {
		DataID DataID
	}
	AssemblyInstrDup  struct{}
	AssemblyInstrLoad struct {
		Offset int
	}
	AssemblyInstrLoadPtr struct {
		Offset  int
		IsLocal bool
	}
	AssemblyInstrStore struct {
		Offset int
	}

	AssemblyInstrLoadFromPtr struct{}
	AssemblyInstrStoreToPtr  struct{}

	AssemblyInstrPrint struct {
		ArgLen int
	}
)

func (*AssemblyInstrAdd) isAssemblyInstruction()         {}
func (*AssemblyInstrSub) isAssemblyInstruction()         {}
func (*AssemblyInstrMul) isAssemblyInstruction()         {}
func (*AssemblyInstrDiv) isAssemblyInstruction()         {}
func (*AssemblyInstrMod) isAssemblyInstruction()         {}
func (*AssemblyInstrNeg) isAssemblyInstruction()         {}
func (*AssemblyInstrEq) isAssemblyInstruction()          {}
func (*AssemblyInstrNE) isAssemblyInstruction()          {}
func (*AssemblyInstrGT) isAssemblyInstruction()          {}
func (*AssemblyInstrLT) isAssemblyInstruction()          {}
func (*AssemblyInstrGTE) isAssemblyInstruction()         {}
func (*AssemblyInstrLTE) isAssemblyInstruction()         {}
func (*AssemblyInstrNot) isAssemblyInstruction()         {}
func (*AssemblyInstrAnd) isAssemblyInstruction()         {}
func (*AssemblyInstrOr) isAssemblyInstruction()          {}
func (*AssemblyInstrJmp) isAssemblyInstruction()         {}
func (*AssemblyInstrJF) isAssemblyInstruction()          {}
func (*AssemblyInstrCall) isAssemblyInstruction()        {}
func (*AssemblyInstrRet) isAssemblyInstruction()         {}
func (*AssemblyInstrPush) isAssemblyInstruction()        {}
func (*AssemblyInstrDup) isAssemblyInstruction()         {}
func (*AssemblyInstrLoad) isAssemblyInstruction()        {}
func (*AssemblyInstrStore) isAssemblyInstruction()       {}
func (*AssemblyInstrLoadFromPtr) isAssemblyInstruction() {}
func (*AssemblyInstrStoreToPtr) isAssemblyInstruction()  {}
func (*AssemblyInstrLoadPtr) isAssemblyInstruction()     {}
func (*AssemblyInstrPrint) isAssemblyInstruction()       {}

func (*AssemblyInstrAdd) String() string         { return "Add" }
func (*AssemblyInstrSub) String() string         { return "Sub" }
func (*AssemblyInstrMul) String() string         { return "Mul" }
func (*AssemblyInstrDiv) String() string         { return "Div" }
func (*AssemblyInstrMod) String() string         { return "Mod" }
func (*AssemblyInstrNeg) String() string         { return "Neg" }
func (*AssemblyInstrEq) String() string          { return "Eq" }
func (*AssemblyInstrNE) String() string          { return "NE" }
func (*AssemblyInstrGT) String() string          { return "GT" }
func (*AssemblyInstrLT) String() string          { return "LT" }
func (*AssemblyInstrGTE) String() string         { return "GTE" }
func (*AssemblyInstrLTE) String() string         { return "LTE" }
func (*AssemblyInstrNot) String() string         { return "Not" }
func (*AssemblyInstrAnd) String() string         { return "And" }
func (*AssemblyInstrOr) String() string          { return "Or" }
func (jmp *AssemblyInstrJmp) String() string     { return fmt.Sprintf("Jmp %s", jmp.Label) }
func (jf *AssemblyInstrJF) String() string       { return fmt.Sprintf("JF %s", jf.Label) }
func (*AssemblyInstrCall) String() string        { return "Call" }
func (*AssemblyInstrRet) String() string         { return "Ret" }
func (p *AssemblyInstrPush) String() string      { return fmt.Sprintf("Push @%d", p.DataID) }
func (*AssemblyInstrDup) String() string         { return "Dup" }
func (l *AssemblyInstrLoad) String() string      { return fmt.Sprintf("Load %d", l.Offset) }
func (s *AssemblyInstrStore) String() string     { return fmt.Sprintf("Store %d", s.Offset) }
func (*AssemblyInstrLoadFromPtr) String() string { return "LoadFromPtr" }
func (*AssemblyInstrStoreToPtr) String() string  { return "StoreToPtr" }
func (lp *AssemblyInstrLoadPtr) String() string {
	s := ""
	if lp.IsLocal {
		s = "local"
	}
	return fmt.Sprintf("Load%sPtr #%d", s, lp.Offset)
}
func (lp *AssemblyInstrPrint) String() string { return fmt.Sprintf("Print %d", lp.ArgLen) }
