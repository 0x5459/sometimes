package parser

import (
	"fmt"
	"sometimes/lexer"
	"testing"
)

func TestParseExpr(t *testing.T) {
	src := "1+2*3"
	tc := lexer.NewTokenCursor(lexer.NewSrcCursor([]byte(src)))
	p := Parser{
		tc: tc,
	}
	p.next()

	e := p.parseExpr()

	fmt.Println(e)
}
