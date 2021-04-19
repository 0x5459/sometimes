package assembly

type DataID = uint32

type AssemblyInstruction interface {
	isAssemblyInstruction()
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
	AssemblyInstrLoad struct {
		Offset int
	}

	AssemblyInstrStore struct {
		Offset int
	}
)

func (*AssemblyInstrAdd) isAssemblyInstruction()   {}
func (*AssemblyInstrSub) isAssemblyInstruction()   {}
func (*AssemblyInstrMul) isAssemblyInstruction()   {}
func (*AssemblyInstrDiv) isAssemblyInstruction()   {}
func (*AssemblyInstrMod) isAssemblyInstruction()   {}
func (*AssemblyInstrNeg) isAssemblyInstruction()   {}
func (*AssemblyInstrEq) isAssemblyInstruction()    {}
func (*AssemblyInstrNE) isAssemblyInstruction()    {}
func (*AssemblyInstrGT) isAssemblyInstruction()    {}
func (*AssemblyInstrLT) isAssemblyInstruction()    {}
func (*AssemblyInstrGTE) isAssemblyInstruction()   {}
func (*AssemblyInstrLTE) isAssemblyInstruction()   {}
func (*AssemblyInstrNot) isAssemblyInstruction()   {}
func (*AssemblyInstrAnd) isAssemblyInstruction()   {}
func (*AssemblyInstrOr) isAssemblyInstruction()    {}
func (*AssemblyInstrJmp) isAssemblyInstruction()   {}
func (*AssemblyInstrJF) isAssemblyInstruction()    {}
func (*AssemblyInstrCall) isAssemblyInstruction()  {}
func (*AssemblyInstrRet) isAssemblyInstruction()   {}
func (*AssemblyInstrPush) isAssemblyInstruction()  {}
func (*AssemblyInstrLoad) isAssemblyInstruction()  {}
func (*AssemblyInstrStore) isAssemblyInstruction() {}
