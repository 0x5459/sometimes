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

func (p *Parser) Parse() (consts []*ast.ConstDecl, fns []*ast.FnDecl) {
	p.next()
	for p.tok.Kind != token.EOF {
		switch p.tok.Kind {
		case token.CONST:
			consts = append(consts, p.parseConstDecl())
		case token.FN:
			fns = append(fns, p.parseFnDecl())
		default:
			p.errorExpect("const' or 'fn")
		}
		p.next()
	}
	return
}

func (p *Parser) next() {
	tok := p.tc.Next()
	for tok.Kind == token.COMMENT {
		tok = p.tc.Next()
	}
	p.tok = tok
}

func (p *Parser) parseConstDecl() *ast.ConstDecl {
	startPos := p.tok.StartPos
	p.expect(token.CONST)
	var l []ast.ValueDecl
	for {
		l = append(l, p.parseValueDecl())
		if p.tok.Kind == token.SEMICOLON || p.tok.Kind == token.EOF {
			p.next()
			break
		} else {
			p.expect(token.COMMA)
		}
	}
	p.next()
	return &ast.ConstDecl{
		BaseNode: ast.NewBaseNode(startPos, p.tok.EndPos),
		Decls:    l,
	}
}

func (p *Parser) parseFnDecl() *ast.FnDecl {
	startPos := p.tok.StartPos
	p.expect(token.FN)
	fnName := p.parseIdent()
	p.expect(token.LPAREN)

	var params []*ast.Ident
	for p.tok.Kind != token.RPAREN && p.tok.Kind != token.EOF {
		params = append(params, p.parseIdent())
		if p.tok.Kind == token.RPAREN || p.tok.Kind == token.EOF {
			break
		} else {
			p.expect(token.COMMA)
		}
	}

	p.next() // eat ')'
	body := p.parseBlockExpr()

	return &ast.FnDecl{
		BaseNode: ast.NewBaseNode(startPos, body.EndPos()),
		FnName:   fnName,
		Args:     params,
		Body:     body,
	}
}
func (p *Parser) parseExpr() ast.Expr {
	switch p.tok.Kind {
	case token.IF:
		return p.parseIfExpr()
	case token.LOOP:
		return p.parseLoopExpr()
	case token.RETURN:
		return p.parseRetExpr()
	case token.LET:
		return p.parseLetExpr()
	case token.BREAK:
		return p.parseBreakExpr()
	case token.CONTINUE:
		return p.parseContinueExpr()
	case token.LBRACK:
		return p.parseArrayExpr()
	default:
		return p.parseBinaryExpr(token.LowestPrec + 1)
	}
}

func (p *Parser) parseUnaryExpr() ast.Expr {
	switch p.tok.Kind {
	case token.NOT, token.SUB:
		op := p.tok
		p.next()
		e := p.parseUnaryExpr()
		return &ast.UnaryExpr{
			BaseExpr: ast.NewBaseExpr(op.StartPos, p.tok.EndPos),
			Op:       op,
			Expr:     e,
		}

	}
	return p.parsePrimaryExpr()
}

func (p *Parser) parseBinaryExpr(prec int) ast.Expr {
	lhs := p.parseUnaryExpr()
	for {
		tokPrec := p.tok.Kind.Precedence()
		if tokPrec < prec {
			return lhs
		}
		op := p.tok
		p.next() // eat op
		rhs := p.parseBinaryExpr(tokPrec + 1)
		lhs = &ast.BinaryExpr{
			BaseExpr: ast.NewBaseExpr(lhs.StartPos(), rhs.EndPos()),
			Lhs:      lhs,
			Rhs:      rhs,
			Op:       op,
		}
	}
}

func (p *Parser) parseOperand() ast.Expr {
	switch p.tok.Kind {
	case token.IDENT:
		x := p.parseIdent()
		return x
	case token.INT_LITERAL, token.FLOAT_LITERAL, token.CHAR_LITERAL, token.STRING_LITERAL:
		x := &ast.Literal{
			BaseExpr: ast.NewBaseExpr(p.tok.StartPos, p.tok.EndPos),
			Kind:     p.tok.Kind,
			Val:      p.tok.Val,
		}
		p.next()
		return x
	case token.LPAREN:
		lparenPos := p.tok.StartPos
		p.next() // eat '('
		inner := p.parseExpr()
		p.expect(token.RPAREN) // eat ')'
		rparenPos := p.tok.EndPos
		return &ast.ParenExpr{
			BaseExpr: ast.NewBaseExpr(lparenPos, rparenPos),
			Inner:    inner,
		}
	}
	p.errorExpect("operand")
	return nil // never
}

func (p *Parser) parsePrimaryExpr() ast.Expr {
	x := p.parseOperand()
	for {
		switch p.tok.Kind {
		case token.LBRACK: // [
			x = p.parseIndexExpr(x)
		case token.LPAREN: // (
			x = p.parseCallExpr(x)
		case token.ASSIGN, token.ADD_ASSIGN, token.MUL_ASSIGN,
			token.QUO_ASSIGN, token.REM_ASSIGN, token.SUB_ASSIGN:
			if _, isIdent := x.(*ast.Ident); isIdent {
				p.next()
				rhs := p.parseExpr()
				p.expect(token.SEMICOLON)
				x = &ast.AssignExpr{
					BaseExpr: ast.NewBaseExpr(x.StartPos(), p.tok.EndPos),
					Lhs:      x,
					Rhs:      rhs,
				}
			}
		default:
			return x
		}
	}
}

func (p *Parser) parseIndexExpr(addr ast.Expr) *ast.IndexExpr {
	p.expect(token.LBRACK)
	e := p.parseExpr()
	p.expect(token.RBRACK)
	return &ast.IndexExpr{
		BaseExpr: ast.NewBaseExpr(addr.StartPos(), p.tok.EndPos),
		Addr:     addr,
		Index:    e,
	}
}

func (p *Parser) parseCallExpr(f ast.Expr) *ast.CallExpr {
	p.expect(token.LPAREN)
	var args []ast.Expr
	for p.tok.Kind != token.RPAREN && p.tok.Kind != token.EOF {
		args = append(args, p.parseExpr())
		if p.tok.Kind == token.RPAREN || p.tok.Kind == token.EOF {
			break
		} else {
			p.expect(token.COMMA) // eat ','
		}
	}
	p.next()
	return &ast.CallExpr{
		BaseExpr: ast.NewBaseExpr(f.StartPos(), p.tok.EndPos),
		Func:     f,
		Args:     args,
	}
}

func (p *Parser) parseIfExpr() *ast.IfExpr {
	startPos := p.tok.StartPos
	p.expect(token.IF)
	cond := p.parseExpr()
	body := p.parseBlockExpr()

	var elseExpr ast.Expr
	if p.tok.Kind == token.ELSE {
		p.next() // eat else
		if p.tok.Kind == token.IF {
			elseExpr = p.parseIfExpr()
		} else {
			elseExpr = p.parseBlockExpr()
		}
	}
	return &ast.IfExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		Cond:     cond,
		Body:     body,
		Else:     elseExpr,
	}
}

func (p *Parser) parseBlockExpr() *ast.BlockExpr {
	startPos := p.tok.StartPos
	p.expect(token.LBRACE)

	var exprs []ast.Expr
	hasRetExpr := false
	for p.tok.Kind != token.RBRACE && p.tok.Kind != token.EOF {
		exprs = append(exprs, p.parseExpr())
		if p.tok.Kind == token.RBRACE || p.tok.Kind == token.EOF {
			hasRetExpr = true
			break
		} else {
			p.expect(token.SEMICOLON)
		}
	}
	p.next()
	var retExpr ast.Expr

	if hasRetExpr {
		retExpr = exprs[len(exprs)-1]
		exprs = exprs[:len(exprs)-1]
	}
	return &ast.BlockExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		ExprList: exprs,
		RetExpr:  retExpr,
	}
}

func (p *Parser) parseLoopExpr() *ast.LoopExpr {
	startPos := p.tok.StartPos
	p.expect(token.LOOP)
	cond := p.parseExpr()
	body := p.parseBlockExpr()
	return &ast.LoopExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		Cond:     cond,
		Body:     body,
	}
}

func (p *Parser) parseRetExpr() *ast.ReturnExpr {
	startPos := p.tok.StartPos
	p.expect(token.RETURN)
	e := p.parseExpr()
	p.expect(token.SEMICOLON)
	return &ast.ReturnExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		Ret:      e,
	}
}

func (p *Parser) parseLetExpr() *ast.LetExpr {
	startPos := p.tok.StartPos
	p.expect(token.LET)
	var l []ast.ValueDecl
	for {
		l = append(l, p.parseValueDecl())
		if p.tok.Kind == token.SEMICOLON || p.tok.Kind == token.EOF {
			p.next()
			break
		} else {
			p.expect(token.COMMA)
		}
	}

	return &ast.LetExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		Decls:    l,
	}
}

func (p *Parser) parseValueDecl() ast.ValueDecl {
	lhs := p.parseIdent()
	p.expect(token.ASSIGN)
	rhs := p.parseExpr()
	return ast.ValueDecl{
		Ident: lhs,
		Value: rhs,
	}
}

func (p *Parser) parseBreakExpr() *ast.BreakExpr {
	startPos := p.tok.StartPos
	p.expect(token.BREAK)
	var expr ast.Expr
	if p.tok.Kind != token.SEMICOLON {
		expr = p.parseExpr()
	}
	p.expect(token.SEMICOLON)
	return &ast.BreakExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		Expr:     expr,
	}
}

func (p *Parser) parseContinueExpr() *ast.ContinueExpr {
	startPos := p.tok.StartPos
	p.expect(token.CONTINUE)
	p.expect(token.SEMICOLON)
	return &ast.ContinueExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
	}
}

func (p *Parser) parseArrayExpr() *ast.ArrayExpr {
	startPos := p.tok.StartPos
	p.expect(token.LBRACK)
	var l []ast.Expr
	for p.tok.Kind != token.RBRACK && p.tok.Kind != token.EOF {
		l = append(l, p.parseExpr())
		if p.tok.Kind == token.RBRACK || p.tok.Kind == token.EOF {
			break
		} else {
			p.expect(token.COMMA)
		}
	}
	p.next()
	return &ast.ArrayExpr{
		BaseExpr: ast.NewBaseExpr(startPos, p.tok.EndPos),
		Element:  l,
	}
}

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

func (p *Parser) errorExpect(expected string) {
	p.error(p.tok.StartPos,
		fmt.Sprintf("expected '%s', found '%s'", expected, p.tok.Kind.String()))
}

func (p *Parser) expect(kind token.Kind) {
	if p.tok.Kind != kind {
		p.errorExpect(kind.String())
	}
	p.next()
}
