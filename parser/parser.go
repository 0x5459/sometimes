package parser

import (
	"fmt"
	"sometimes/lexer"
)

// parseFunction parse a function
// fn f() -> int {}
func parseFunction(tc *lexer.TokenCursor)  {
	tc.Next() // eat 'fn'
	fmt.Print()
}


