package hir

//go:generate stringer -type=ExprType -trimprefix=ExprType
type ExprType int

const (
	ExprTypeLiteral ExprType = iota
	ExprTypeConst            // decl const
	ExprTypeVar              // access binding
	ExprTypeBinding
	ExprTypeMutate
	ExprTypeBinary
	ExprTypeCall
	ExprTypeFunction
	ExprTypeAnonFunction
	ExprTypeUnary
	ExprTypeReturn
	ExprTypeIf
	ExprTypeLoop
	ExprTypeBlock
	ExprTypeBreak
	ExprTypeContinue
	ExprTypeArray
	ExprTypeSetElement
	ExprTypeGetElement
	ExprTypePrint
)

type Expr interface {
	ExprType() ExprType
}

type (
	ExprLiteral struct {
		Val Value
	}
	ExprConst struct {
	}

	ExprVar struct {
		VarBinding *Binding
	}

	ExprBinding struct {
		Binding *Binding
		Rhs     Expr
	}

	// like `a = a + 1`
	ExprMutate struct {
		Lhs, Rhs Expr
	}

	ExprBinary struct {
		Lhs, Rhs Expr
		Op       BinaryOp
	}

	ExprCall struct {
		Callee Expr
		Args   []Expr
	}

	ExprFunction struct {
		Func *Function
	}

	ExprAnonFunction struct {
		Func *Function
	}

	ExprUnary struct {
		Op   UnaryOp
		Expr Expr
	}

	ExprReturn struct {
		Expr Expr // optional
	}

	ExprIf struct {
		Cond Expr
		Body *ExprBlock
		Else Expr // optional
	}

	ExprLoop struct {
		Cond Expr
		Body *ExprBlock
	}

	ExprBlock struct {
		Body []Expr
	}

	ExprBreak struct {
		Expr Expr // optional
	}

	ExprContinue struct{}

	ExprArray struct {
		Arr   *ExprVar
		Exprs []Expr
	}

	ExprSetElement struct {
		ArrayAddr, Index, Value Expr
	}
	ExprGetElement struct {
		ArrayAddr, Index Expr
	}
	ExprPrint struct {
		Expr []Expr
	}
)

func (*ExprLiteral) ExprType() ExprType      { return ExprTypeLiteral }
func (*ExprVar) ExprType() ExprType          { return ExprTypeVar }
func (*ExprBinding) ExprType() ExprType      { return ExprTypeBinding }
func (*ExprMutate) ExprType() ExprType       { return ExprTypeMutate }
func (*ExprBinary) ExprType() ExprType       { return ExprTypeBinary }
func (*ExprCall) ExprType() ExprType         { return ExprTypeCall }
func (*ExprFunction) ExprType() ExprType     { return ExprTypeFunction }
func (*ExprAnonFunction) ExprType() ExprType { return ExprTypeAnonFunction }
func (*ExprUnary) ExprType() ExprType        { return ExprTypeUnary }
func (*ExprReturn) ExprType() ExprType       { return ExprTypeReturn }
func (*ExprIf) ExprType() ExprType           { return ExprTypeIf }
func (*ExprLoop) ExprType() ExprType         { return ExprTypeLoop }
func (*ExprBlock) ExprType() ExprType        { return ExprTypeBlock }
func (*ExprBreak) ExprType() ExprType        { return ExprTypeBreak }
func (*ExprContinue) ExprType() ExprType     { return ExprTypeContinue }
func (*ExprArray) ExprType() ExprType        { return ExprTypeArray }
func (*ExprSetElement) ExprType() ExprType   { return ExprTypeSetElement }
func (*ExprGetElement) ExprType() ExprType   { return ExprTypeGetElement }
func (*ExprPrint) ExprType() ExprType        { return ExprTypePrint }
