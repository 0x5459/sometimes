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
	BaseNode struct {
		startPos, endPos Pos
	}
)

func (bn *BaseNode) StartPos() Pos { return bn.startPos }
func (bn *BaseNode) EndPos() Pos   { return bn.endPos }

func NewBaseNode(startPos, endPos Pos) *BaseNode {
	return &BaseNode{
		startPos: startPos,
		endPos:   endPos,
	}
}

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
		Ident Ident
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

	BaseExpr struct {
		*BaseNode
	}

	Ident struct {
		*BaseNode
		Name string
	}

	Literal struct {
		*BaseExpr
		Kind token.Kind // token.INT_LITERAL, token.FLOAT_LITERAL, token.CHAR_LITERAL, token.STRING_LITERAL
		Val  string     // 123, 3.14, 'a', "我的"
	}

	// (1+1), ...
	ParenExpr struct {
		*BaseExpr
		Inner Expr // 1+1
	}

	// arr[2+2], ...
	IndexExpr struct {
		*BaseExpr
		Addr  Ident // arr
		Index Expr  // 2+2
	}

	// [1+1, a(), 11], ...
	ArrayExpr struct {
		*BaseExpr
		Element []Expr
	}

	// say(1+1, 99), ...
	CallExpr struct {
		*BaseExpr
		Func Expr
		Args []Expr
	}

	// !b, ...
	UnaryExpr struct {
		*BaseExpr
		Op   token.Token // !
		Expr Expr        // b
	}

	// 1+1, ...
	BinaryExpr struct {
		*BaseExpr
		Lhs, Rhs Expr
		Op       *token.Token
	}

	// a = 10, ...
	AssignExpr struct {
		*BaseExpr
		Lhs, Rhs Expr
	}

	// return 1+1, ...
	ReturnExpr struct {
		*BaseExpr
		Ret Expr
	}

	// break 1+1, ...
	BreakExpr struct {
		*BaseExpr
		Expr Expr // optional
	}

	// continue
	ContinueExpr struct {
		*BaseExpr
	}

	// { a = 1+1; b = 2+2; }
	BlockExpr struct {
		*BaseExpr
		ExprList []Expr
	}

	// if Cond { Body } else { Else }
	IfExpr struct {
		*BaseExpr
		Cond Expr // condition
		Body Expr
		Else Expr // else expr; optional
	}

	// loop (Cond) { Body }
	LoopExpr struct {
		*BaseExpr
		Cond Expr // condition; optional
		Body BlockExpr
	}

	// let (a = 10, b=20, 1+1), ...
	LetExpr struct {
		*BaseExpr
		Decls    []ValueDecl
		LastExpr Expr // optional
	}

	// // type (
	// //   a = 100,
	// //   Person struct {
	// //      age int,
	// //   }
	// //)
	// TypeExpr struct {
	// 	baseExpr
	// 	Decls []TypeDecl
	// }
)

func (BaseNode) exprNode() {}

func NewBaseExpr(startPos, endPos Pos) *BaseExpr {
	return &BaseExpr{
		BaseNode: NewBaseNode(startPos, endPos),
	}
}

// declares
type (
	// let Ident = 100;
	ValueDecl struct {
		Ident Ident
		Value Expr
	}
)

// const (A=10, B=1+1)
type ConstDecl struct {
	BaseNode
	Decls []ValueDecl
}

func (cd *ConstDecl) StartPos() Pos { return cd.startPos }
func (cd *ConstDecl) EndPos() Pos   { return cd.endPos }

// fn f(n int) -> int { xxx }
type FnDecl struct {
	BaseNode
	Args []FieldDef
	Ret  Type
	Body BlockExpr
}

func (fd *FnDecl) StartPos() Pos { return fd.startPos }
func (fd *FnDecl) EndPos() Pos   { return fd.endPos }
