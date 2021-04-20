package assembly

import (
	"fmt"
	"sometimes/hir"
	"sync/atomic"
)

type compileState struct {
	localIdx int
	locals   map[string]int // local varname -> localIdx
}

func newCompileState() *compileState {
	return &compileState{
		localIdx: 0,
		locals:   make(map[string]int),
	}
}

func (cs *compileState) StoreVar(b *hir.Binding) *AssemblyInstrStore {
	localIdx, ok := cs.locals[b.Name]
	if !ok {
		localIdx = cs.localIdx
		cs.locals[b.Name] = localIdx
		cs.localIdx++
	}
	return &AssemblyInstrStore{Offset: localIdx}

}

func (cs *compileState) IsLocalVar(b *hir.Binding) bool {
	_, ok := cs.locals[b.Name]
	return ok
}

func (cs *compileState) LoadVar(b *hir.Binding) *AssemblyInstrLoad {
	return &AssemblyInstrLoad{Offset: cs.locals[b.Name]}
}

func (cs *compileState) MaxLocals() int {
	return len(cs.locals)
}

type compileStateStack []*compileState

func (s *compileStateStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *compileStateStack) Last() *compileState {
	if s.IsEmpty() {
		panic("CompileStateStack must not be empty")
	}
	return (*s)[len(*s)-1]
}

func (s *compileStateStack) Push(cs *compileState) {
	*s = append(*s, cs)
}

func (s *compileStateStack) Pop() (cs *compileState, isExist bool) {
	if !s.IsEmpty() {
		idx := len(*s) - 1
		cs = (*s)[idx]
		*s = (*s)[:idx]
		isExist = true
	}
	return
}

/// Compile hir to assembly
type Compiler struct {
	labelGen            *LabelGen
	loopLabelStack      LoopLabelStack
	asm                 *AssemblyProgram
	hirProgram          *hir.Program
	states              compileStateStack
	constsDataIdMapping map[string]DataID
}

func NewCompiler(hirProgram *hir.Program) *Compiler {
	return &Compiler{
		labelGen:   NewLabelGen(),
		hirProgram: hirProgram,
		asm:        NewAssemblyProgram(),
		states: []*compileState{ // todo temp
			newCompileState(),
		},
		constsDataIdMapping: make(map[string]DataID),
	}
}

func (c *Compiler) Compile() *AssemblyProgram {
	c.saveConsts()

	funcs := c.hirProgram.Funcs()
	// save funcs to consts
	for _, f := range funcs {
		dataID := c.asm.Consts.insertConst(&hir.ValueFunc{
			FuncName:  f.Func.Name,
			MaxLoacls: 0, // MaxLocals compute after compile this func
		})
		c.constsDataIdMapping[f.Func.Name] = dataID
	}

	// call entry function
	entryFunc := c.hirProgram.EntryFunc()
	c.asm.Emit(&AssemblyInstrPush{DataID: c.FindConst(entryFunc.Func.Name)})
	c.asm.Emit(&AssemblyInstrCall{})

	for _, f := range funcs {
		c.compileExpr(f)
	}
	return c.asm
}

func (c *Compiler) saveConsts() {
	for constName, cnst := range c.hirProgram.Consts() {
		dataID := c.asm.Consts.insertConst(cnst)
		c.constsDataIdMapping[constName] = dataID
	}
}

func (c *Compiler) FindConst(name string) DataID {
	return c.constsDataIdMapping[name]
}

func (c *Compiler) compileExpr(expr hir.Expr) {
	switch e := expr.(type) {
	case *hir.ExprLiteral:
		c.asm.EmitPush(e.Val)
	case *hir.ExprVar:
		state := c.states.Last()
		var instr AssemblyInstruction
		if state.IsLocalVar(e.VarBinding) {
			instr = c.states.Last().LoadVar(e.VarBinding)
		} else {
			instr = &AssemblyInstrPush{DataID: c.FindConst(e.VarBinding.Name)}
		}
		c.asm.Emit(instr)
	case *hir.ExprMutate:
		if variable, ok := e.Lhs.(*hir.ExprVar); ok {
			c.compileExpr(e.Rhs)

			instr := c.states.Last().StoreVar(variable.VarBinding)
			c.asm.Emit(instr)
		} else {
			panic("mutate non-variable")
		}

	case *hir.ExprBinary:
		c.compileExpr(e.Lhs)
		c.compileExpr(e.Rhs)
		c.asm.Emit(hirBinaryOpToAssemblyInstr(e.Op))
	case *hir.ExprCall:
		for i := len(e.Args) - 1; i >= 0; i-- {
			c.compileExpr(e.Args[i])
		}
		c.compileExpr(e.Callee)
		c.asm.Emit(&AssemblyInstrCall{})
	case *hir.ExprFunction:
		c.states.Push(newCompileState())
		c.asm.Label(e.Func.Name)
		for _, arg := range e.Func.Args {
			instr := c.states.Last().StoreVar(arg)
			c.asm.Emit(instr)
		}
		c.compileExpr(e.Func.Body)
		state, _ := c.states.Pop()
		cnst := c.asm.Consts.GetConst(c.FindConst(e.Func.Name)).(*hir.ValueFunc)
		cnst.MaxLoacls = state.MaxLocals()
		state.MaxLocals()
	case *hir.ExprAnonFunction:
		panic("unimplement!")
	case *hir.ExprUnary:
		c.compileExpr(e.Expr)
		switch e.Op {
		case hir.OpNeg:
			c.asm.Emit(&AssemblyInstrNeg{})
		case hir.OpNot:
			c.asm.Emit(&AssemblyInstrNot{})
		}
	case *hir.ExprReturn:
		if e.Expr != nil {
			c.compileExpr(e.Expr)
		}
		c.asm.Emit(&AssemblyInstrRet{})
	case *hir.ExprIf:
		c.compileExpr(e.Cond)
		elseLabel, endifLabel := c.labelGen.NextIfLabel()
		jf := &AssemblyInstrJF{}
		if e.Else != nil {
			jf.Label = elseLabel
		} else {
			jf.Label = endifLabel
		}
		c.asm.Emit(jf)
		c.compileExpr(e.Body)
		if e.Else != nil {
			c.asm.Emit(&AssemblyInstrJmp{Label: endifLabel})
			c.asm.Label(elseLabel)
			c.compileExpr(e.Else)
		}
		c.asm.Label(endifLabel)
	case *hir.ExprLoop:
		loopStartLabel, loopEndLabel := c.labelGen.NextLoopLabel()
		c.loopLabelStack.StartLoop(loopStartLabel, loopEndLabel)
		c.asm.Label(loopStartLabel)
		c.compileExpr(e.Cond)
		c.asm.Emit(&AssemblyInstrJF{Label: loopEndLabel})
		c.compileExpr(e.Body)
		c.asm.Emit(&AssemblyInstrJmp{Label: loopStartLabel})
		c.asm.Label(loopEndLabel)
		c.loopLabelStack.EndLoop()
	case *hir.ExprBlock:
		for _, body := range e.Body {
			c.compileExpr(body)
		}
	case *hir.ExprBreak:
		c.compileExpr(e.Expr)
		_, loopEndLabel := c.loopLabelStack.CurrentLabel()
		c.asm.Emit(&AssemblyInstrJmp{Label: loopEndLabel})
	case *hir.ExprContinue:
		loopStartLabel, _ := c.loopLabelStack.CurrentLabel()
		c.asm.Emit(&AssemblyInstrJmp{Label: loopStartLabel})
	}
}

type LabelGen struct {
	ifID, loopID uint32
}

func NewLabelGen() *LabelGen {
	return &LabelGen{
		ifID:   0,
		loopID: 0,
	}
}

func (lg *LabelGen) NextIfLabel() (elseIf, endIf string) {
	ifID := lg.ifID
	atomic.AddUint32(&lg.ifID, 1)
	return fmt.Sprintf("else-%d", ifID), fmt.Sprintf("endif-%d", ifID)
}

func (lg *LabelGen) NextLoopLabel() (loopStart, loopEnd string) {
	atomic.AddUint32(&lg.loopID, 1)
	return fmt.Sprintf("loopStart-%d", lg.loopID), fmt.Sprintf("loopEnd-%d", lg.loopID)
}

type LoopLabelStack []struct{ loopStart, loopEnd string }

func (l *LoopLabelStack) StartLoop(loopStart, loopEnd string) {

	d := struct {
		loopStart string
		loopEnd   string
	}{loopStart: loopStart, loopEnd: loopEnd}
	*l = append(*l, d)
}

func (l *LoopLabelStack) EndLoop() {
	if len(*l) != 0 {
		idx := len(*l) - 1
		*l = (*l)[:idx]
	}
}

func (l *LoopLabelStack) CurrentLabel() (loopStart, loopEnd string) {
	label := (*l)[len(*l)-1]
	return label.loopStart, label.loopEnd
}

func hirBinaryOpToAssemblyInstr(bop hir.BinaryOp) AssemblyInstruction {
	switch bop {
	case hir.OpAdd:
		return &AssemblyInstrAdd{}
	case hir.OpSub:
		return &AssemblyInstrSub{}
	case hir.OpMul:
		return &AssemblyInstrMul{}
	case hir.OpDiv:
		return &AssemblyInstrDiv{}
	case hir.OpMod:
		return &AssemblyInstrMod{}
	case hir.OpEq:
		return &AssemblyInstrEq{}
	case hir.OpNE:
		return &AssemblyInstrNE{}
	case hir.OpGT:
		return &AssemblyInstrGT{}
	case hir.OpLT:
		return &AssemblyInstrLT{}
	case hir.OpGTE:
		return &AssemblyInstrGTE{}
	case hir.OpLTE:
		return &AssemblyInstrLTE{}
	}
	panic(fmt.Errorf("unsupport op: %d", bop))
}
