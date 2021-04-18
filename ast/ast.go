package ast

import (
	"sometimes/token"
)

type (
	Pos interface {
		Line() int
		Col() int
	}
	// Node represents an AST node.
	Node interface {
		StartPos() Pos
		EndPos() Pos
	}
	baseNode struct {
		startPos, endPos Pos
	}
)

func (bn *baseNode) StartPos() Pos { return bn.startPos }
func (bn *baseNode) EndPos() Pos   { return bn.endPos }

type (
	Type interface {
		// ty ensures that only `type` can be assigned to a Type.
		ty()
	}

	// [Len]Type, ..
	ArrayType struct {
		Type Type
		Len  Expr
	}

	// age int , struct { age int, }
	FieldDef struct {
		ident Ident
		Type  Type
	}
)

func (*Ident) ty()     {}
func (*ArrayType) ty() {}

type (
	Expr interface {
		Node
		// exprNode ensures that only expression nodes can be assigned to an Expr.
		exprNode()
	}

	baseExpr struct {
		*baseNode
	}

	Ident struct {
		*baseNode
		Name string
	}

	Literal struct {
		baseExpr
		Kind token.Kind // token.INT_LITERAL, token.FLOAT_LITERAL, token.CHAR_LITERAL, token.STRING_LITERAL
		Val  string     // 123, 3.14, 'a', "我的"
	}

	// (1+1), ...
	ParenExpr struct {
		baseExpr
		inner Expr // 1+1
	}

	// arr[2+2], ...
	IndexExpr struct {
		baseExpr
		Addr  Ident // arr
		index Expr  // 2+2
	}

	// [1+1, a(), 11], ...
	ArrayExpr struct {
		baseExpr
		Element []Expr
	}

	// say(1+1, 99), ...
	CallExpr struct {
		baseExpr
		Func Expr
		Args []Expr
	}

	// !b, ...
	UnaryExpr struct {
		baseExpr
		Op   token.Token // !
		Expr Expr        // b
	}

	// 1+1, ...
	BinaryExpr struct {
		baseExpr
		Lhs, Rhs Expr
		Op       token.Token
	}

	// a = 10, ...
	AssignExpr struct {
		baseExpr
		Lhs, Rhs Expr
	}

	// return 1+1, ...
	ReturnExpr struct {
		baseExpr
		Ret Expr
	}

	// break 1+1, ...
	BreakExpr struct {
		baseExpr
		Expr Expr // optional
	}

	// continue
	ContinueExpr struct {
		baseExpr
	}

	// { a = 1+1; b = 2+2; }
	BlockExpr struct {
		baseExpr
		ExprList []Expr
	}

	// if Cond { Body } else { Else }
	IfExpr struct {
		baseExpr
		Cond Expr // condition
		Body Expr
		Else Expr // else expr; optional
	}

	// loop (Cond) { Body }
	LoopExpr struct {
		baseExpr
		Cond Expr // condition; optional
		Body BlockExpr
	}

	// let (a = 10, b=20, 1+1), ...
	LetExpr struct {
		baseExpr
		Decls    []ValueDecl
		LastExpr Expr // optional
	}

	// const (A=10, B=1+1)
	ConstExpr struct {
		baseExpr
		Decls []ValueDecl
	}

	// type (
	//   a = 100,
	//   Person struct {
	//      age int,
	//   }
	//)
	TypeExpr struct {
		baseExpr
		Decls []TypeDecl
	}
)

func (baseNode) exprNode() {}

// declares
type (
	// let Ident = 100;
	ValueDecl struct {
		Ident Ident
		Value Expr
	}
)

// fn f(n int) -> int { xxx }
type FnDecl struct {
	baseNode
	Args []FieldDef
	Ret  Type
	Body BlockExpr
}

func (fd *FnDecl) StartPos() Pos { return fd.startPos }
func (fd *FnDecl) EndPos() Pos   { return fd.endPos }
