package hir

type Program struct {
	funcs         map[string]*ExprFunction
	entryFuncName string
	consts        map[string]Value
}

func (p *Program) FindConst(name string) (val Value, isExist bool) {
	val, isExist = p.consts[name]
	return
}

func (p *Program) FindFunc(name string) (f *ExprFunction, isExist bool) {
	f, isExist = p.funcs[name]
	return
}

func (p *Program) EntryFunc() *ExprFunction {
	return p.funcs[p.entryFuncName]
}

func (p *Program) Funcs() []*ExprFunction {
	funcs := make([]*ExprFunction, 0, len(p.funcs))
	for _, f := range p.funcs {
		funcs = append(funcs, f)
	}
	return funcs
}

func (p *Program) Consts() map[string]Value {
	return p.consts
}

type Builder struct {
	Funcs        []*ExprFunction
	EntryFuncIdx int
	consts       map[string]Value
}

func NewBuilder() *Builder {
	return &Builder{
		Funcs:        []*ExprFunction{},
		EntryFuncIdx: -1,
		consts:       make(map[string]Value),
	}
}

func (b *Builder) InsertConst(name string, val Value) {
	b.consts[name] = val
}

func (b *Builder) InsertFunc(f *ExprFunction, entryFunc bool) {
	if entryFunc {
		b.EntryFuncIdx = len(b.Funcs)
	}
	b.Funcs = append(b.Funcs, f)
}

func (b *Builder) Build() *Program {
	funcs := make(map[string]*ExprFunction, len(b.Funcs))

	for _, f := range b.Funcs {
		funcs[f.Func.Name] = f
	}

	var entryFuncName string

	if b.EntryFuncIdx == -1 {
		panic("empty program")
	}
	entryFuncName = (b.Funcs[b.EntryFuncIdx]).Func.Name
	return &Program{
		funcs:         funcs,
		entryFuncName: entryFuncName,
		consts:        b.consts,
	}
}

type FuncBuilder struct {
	funcName string
	funcBody []Expr
	args     []*Binding
}

func NewFuncBuilder(funcName string, args []*Binding) *FuncBuilder {
	return &FuncBuilder{
		funcName: funcName,
		funcBody: []Expr{},
		args:     args,
	}
}

func (b *FuncBuilder) Emit(e Expr) {
	b.funcBody = append(b.funcBody, e)
}

func (b *FuncBuilder) Build() *ExprFunction {
	return &ExprFunction{
		Func: &Function{
			Name: b.funcName,
			Body: &ExprBlock{
				Body: b.funcBody,
			},
			Args: b.args,
		},
	}
}

type Binding struct {
	Name string
	Len  int
}

func NewBinding(name string) *Binding {
	return NewBindingWithLen(name, 1)
}

func NewBindingWithLen(name string, len int) *Binding {
	return &Binding{
		Name: name,
		Len:  len,
	}
}

type Function struct {
	Name string
	Body *ExprBlock
	Args []*Binding
}

type BinaryOp uint8

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpEq  // equal
	OpNE  // not equal
	OpGT  // greater than
	OpLT  // less than
	OpGTE // greater than equal
	OpLTE // less than equal
	OpAnd
	OpOr
	OpIndex
)

type UnaryOp uint8

const (
	OpNeg UnaryOp = iota
	OpNot
)
