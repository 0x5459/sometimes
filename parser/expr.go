package parser

import (
	"sometimes/ast"
	"sometimes/token"
)

func (p *Parser) parseExpr() ast.Expr {
	lhs := p.parsePrimaryExpr()
	if lhs == nil {
		return nil
	}
	return p.parseBinOpRhs(token.LowestPrec, lhs)
}

func (p *Parser) parseBinOpRhs(exprPrec int, lhs ast.Expr) ast.Expr {
	for {
		tokPrec := p.tok.Kind.Precedence()
		if tokPrec < exprPrec {
			return lhs
		}

		binOp := p.tok
		p.next() // eat binOp
		rhs := p.parsePrimaryExpr()
		if rhs == nil {
			return nil
		}
		nextPrec := p.tok.Kind.Precedence()
		if tokPrec < nextPrec {
			rhs = p.parseBinOpRhs(tokPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = &ast.BinaryExpr{
			BaseExpr: ast.NewBaseExpr(lhs.StartPos(), rhs.EndPos()),
			Lhs:      lhs,
			Rhs:      rhs,
			Op:       binOp,
		}
	}
}

func (p *Parser) parsePrimaryExpr() ast.Expr {
	switch p.tok.Kind {
	case token.IDENT:
		x := p.parseIdent()
		return x
	case token.INT_LITERAL, token.FLOAT_LITERAL, token.CHAR_LITERAL, token.STRING_LITERAL:
		x := &ast.Literal{BaseExpr: ast.NewBaseExpr(p.tok.StartPos, p.tok.EndPos)}
		p.next()
		return x
	case token.LPAREN:
		lparenPos := p.tok.StartPos
		p.next()
		inner := p.parseExpr()
		p.expect(token.RPAREN)
		rparenPos := p.tok.EndPos
		return &ast.ParenExpr{
			BaseExpr: ast.NewBaseExpr(lparenPos, rparenPos),
			Inner:    inner,
		}
	}
	return nil
}
