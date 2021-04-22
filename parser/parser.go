package parser

import (
	"fmt"
	"sometimes/ast"
	"sometimes/lexer"
	"sometimes/token"
)

type Parser struct {
	tc  *lexer.TokenCursor
	tok *token.Token
}

func (p *Parser) parse() {
	// consts := make([]*ast.ConstDecl, 0)
	// fns := make([]*ast.FnDecl, 0)
	// p.next()
	// switch p.tok.Kind {
	// case token.CONST:
	// 	consts = append(consts, parseConstsDecl(tc)...)
	// case token.FN:
	// 	fns = append(fns, parseFnDecl(tc))
	// }
}

func (p *Parser) next() {
	p.tok = p.tc.Next()
}

// // const (A=10, B=1+1)
// func (p *Parser) parseConstsDecl(tc *lexer.TokenCursor) *ast.ConstDecl {
// 	return &ast.ConstDecl{}
// }

// func (p *Parser) parseValueDecl() []*ast.ValueDecl {
// 	vals := make([]*ast.ValueDecl, 0)
// 	if p.tok.Kind == token.LPAREN {
// 		p.next()
// 		for p.tok.Kind != token.RPAREN && p.tok.Kind != token.EOF {
// 			vals = append(vals, p.parseSingleValueDecl())
// 		}
// 	} else {
// 		vals = append(vals, p.parseSingleValueDecl())
// 	}
// }

// func (p *Parser) parseSingleValueDecl() *ast.ValueDecl {
// 	ident := p.parseIdent()

// }

// func (p *Parser) parseUnaryExpr() ast.Expr {
// 	switch p.tok.Kind {
// 	case token.SUB, token.NOT:
// 		startPos := p.tok.StartPos
// 		p.next()
// 		expr := p.parseUnaryExpr()
// 		return &ast.UnaryExpr{
// 			BaseExpr: ast.NewBaseExpr(startPos, expr.EndPos()),

// 			Expr: expr,
// 		}
// 	}
// 	return p.parsePrimaryExpr()
// }

func (p *Parser) parseIdent() *ast.Ident {
	name := "_"
	if p.tok.Kind == token.IDENT {
		name = p.tok.Val
		p.next()
	} else {
		p.expect(token.IDENT)
	}
	return &ast.Ident{
		BaseNode: ast.NewBaseNode(p.tok.StartPos, p.tok.EndPos),
		Name:     name,
	}
}

func (p *Parser) error(pos token.Pos, msg string) {
	panic(fmt.Errorf("-> line %d, column %d\n%s", pos.Line(), pos.Col(), msg))
}

func (p *Parser) expect(kind token.Kind) {
	if p.tok.Kind != kind {
		p.error(p.tok.StartPos,
			fmt.Sprintf("expected %s, found %s", kind.String(), p.tok.Kind.String()))
	}
	p.next()
}
