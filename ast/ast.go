package ast

import (
	"fmt"
	"sometimes/token"
	"strings"
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

// type (
// Type interface {
// 	// ty ensures that only `type` can be assigned to a Type.
// 	ty()
// }

// // [Len]Type, ..
// ArrayType struct {
// 	Type Type
// 	Len  Expr
// }

// // age int , struct { age int, }
// FieldDef struct {
// 	Ident Ident
// 	Type  Type
// }
// )

// func (*Ident) ty()     {}
// func (*ArrayType) ty() {}

type (
	Expr interface {
		Node
		// exprNode ensures that only expression nodes can be assigned to an Expr.
		exprNode()
		String() string
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
		Kind token.Kind // token.INT_LITERAL, token.FLOAT_LITERAL, token.CHAR_LITERAL, token.STRING_LITERAL, token.BOOLEAN_LITERAL
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
		Addr  Expr // arr
		Index Expr // 2+2
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
		Op   *token.Token // !
		Expr Expr         // b
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
		Op       *token.Token // =, +=, *=, ...
	}

	// return 1+1, ...
	ReturnExpr struct {
		*BaseExpr
		Ret Expr // optional
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
		RetExpr  Expr
	}

	// if Cond { Body } else { Else }
	IfExpr struct {
		*BaseExpr
		Cond Expr // condition
		Body *BlockExpr
		Else Expr // else expr; optional
	}

	// loop (Cond) { Body }
	LoopExpr struct {
		*BaseExpr
		Cond Expr // condition; optional
		Body *BlockExpr
	}

	// let a = 10, b=20, ...
	LetExpr struct {
		*BaseExpr
		Decls []ValueDecl
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
		Ident *Ident
		Value Expr
	}
)

func (v ValueDecl) String() string {
	return v.Ident.String() + " = " + v.Value.String()
}

// const A=10, B=1+1
type ConstDecl struct {
	*BaseNode
	Decls []ValueDecl
}

func (cd *ConstDecl) StartPos() Pos { return cd.startPos }
func (cd *ConstDecl) EndPos() Pos   { return cd.endPos }

func (cd *ConstDecl) String() string {
	var sb strings.Builder
	sb.WriteString("const ")
	for i := len(cd.Decls) - 2; i >= 0; i-- {
		sb.WriteString(cd.Decls[i].String())
		sb.WriteString(", ")
	}
	sb.WriteString(cd.Decls[len(cd.Decls)-1].String())
	sb.WriteRune(';')
	return sb.String()
}

// fn f(n) { xxx }
type FnDecl struct {
	*BaseNode
	FnName *Ident
	Args   []*Ident
	// Ret  Type
	Body *BlockExpr
}

func (fd *FnDecl) StartPos() Pos { return fd.startPos }
func (fd *FnDecl) EndPos() Pos   { return fd.endPos }

func (id *Ident) String() string {
	return id.Name
}

func (l *Literal) String() string {
	return l.Val
}

func (p *ParenExpr) String() string {
	return "(" + p.Inner.String() + ")"
}

func (i *IndexExpr) String() string {
	return i.Addr.String() + "[" + i.Index.String() + "]"
}

func (a *ArrayExpr) String() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i := len(a.Element) - 2; i >= 0; i-- {
		sb.WriteString(a.Element[i].String())
		sb.WriteString(", ")
	}

	sb.WriteString(a.Element[len(a.Element)-1].String())
	sb.WriteRune(']')
	return sb.String()
}

func (c *CallExpr) String() string {
	var sb strings.Builder
	sb.WriteString(c.Func.String())
	sb.WriteRune('(')
	for i := len(c.Args) - 2; i >= 0; i-- {
		sb.WriteString(c.Args[i].String())
		sb.WriteString(", ")
	}
	sb.WriteString(c.Args[len(c.Args)-1].String())
	sb.WriteRune(')')
	return sb.String()
}

func (u *UnaryExpr) String() string {
	return u.Op.Val + u.Expr.String()
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s%s%s)", b.Lhs.String(), b.Op.Val, b.Rhs.String())
}

func (a *AssignExpr) String() string {
	return a.Lhs.String() + "=" + a.Rhs.String()
}

func (r *ReturnExpr) String() string {
	s := "return "
	if r.Ret != nil {
		s += r.Ret.String()
	}
	return s
}

func (b *BreakExpr) String() string {
	s := "break"
	if b.Expr != nil {
		s += b.Expr.String()
	}
	return s
}

func (c *ContinueExpr) String() string {
	return "continue"
}

func (b *BlockExpr) String() string {
	var sb strings.Builder
	sb.WriteRune('{')
	sb.WriteRune('\n')
	for _, e := range b.ExprList {
		sb.WriteString(e.String())
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	if b.RetExpr != nil {
		sb.WriteString(b.RetExpr.String())
		sb.WriteRune('\n')
	}
	sb.WriteRune('}')
	return sb.String()
}

func (i *IfExpr) String() string {
	var sb strings.Builder
	sb.WriteString("if ")
	sb.WriteString(i.Cond.String())
	sb.WriteRune(' ')
	sb.WriteString(i.Body.String())
	if i.Else != nil {
		sb.WriteString(" else ")
		sb.WriteString(i.Else.String())
	}
	return sb.String()
}

func (l *LoopExpr) String() string {
	var sb strings.Builder
	sb.WriteString("loop ")
	sb.WriteString(l.Cond.String())
	sb.WriteRune(' ')
	sb.WriteString(l.Body.String())
	return sb.String()
}

func (l *LetExpr) String() string {
	var sb strings.Builder
	sb.WriteString("let ")
	for i := len(l.Decls) - 2; i >= 0; i-- {
		sb.WriteString(l.Decls[i].String())
		sb.WriteString(", ")
	}
	sb.WriteString(l.Decls[len(l.Decls)-1].String())
	sb.WriteRune(';')
	return sb.String()
}
