package parser

import (
	"sometimes/lexer"
	"testing"
)

func TestParse(t *testing.T) {
	code := `
// aasdasd
if (a > 10) {
	
}
`
	parser := Parser{
		tc: lexer.NewTokenCursor(lexer.NewSrcCursor([]byte(code))),
	}

	parser.Parse()
}
