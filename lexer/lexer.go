package lexer

import (
	"sometimes/token"
	"strings"
	"unicode"
)

// Using Char as an alias for byte.
// Because our simple language only support ASCII
type Char = byte

type SrcPos struct {
	line, col int
}

func (s *SrcPos) Line() int {
	return s.line
}

func (s *SrcPos) Col() int {
	return s.col
}

type SrcCursor struct {
	src      []Char // source code
	rdOffset int    // read offset of src
	rdPos    SrcPos
}

func NewSrcCursor(src []Char) *SrcCursor {
	return &SrcCursor{
		src:      src,
		rdOffset: -1,
		rdPos: SrcPos{
			line: 0,
			col:  0,
		},
	}
}

/// Eof returns true if rdOffset is end of src
func (sc *SrcCursor) Eof() bool {
	return sc.rdOffset >= len(sc.src)-1
}

// Peek returns the next character without advancing the cursor.
// If the cursor is at EOF, Peek returns 0.
func (sc *SrcCursor) Peek() Char {
	return sc.PeekN(1)
}

func (sc *SrcCursor) PeekN(n int) Char {
	if sc.rdOffset+n >= len(sc.src) {
		return 0
	}
	return sc.src[sc.rdOffset+n]
}

// Next Moves to the next character.
// If the cursor is at EOF, Next returns 0.
func (sc *SrcCursor) Next() Char {
	if sc.Eof() {
		return 0
	}
	sc.rdOffset++
	ch := sc.src[sc.rdOffset]
	if ch == '\n' {
		sc.rdPos.line++
		sc.rdPos.col = 0
	} else {
		sc.rdPos.col++
	}
	return ch
}

// EatWhile eats item while predicate returns true or until the end of source is reached.
func (sc *SrcCursor) EatWhile(predicate func(c Char) bool) string {

	var v strings.Builder
	for ch := sc.Peek(); ch > 0 && predicate(ch); ch = sc.Peek() {
		v.WriteByte(ch)
		sc.Next()
	}
	return v.String()
}

/// TokenCursor 代表一个token游标的源代码.
type TokenCursor struct {
	sc *SrcCursor
}

/// NewTokenCursor 返回一个新的token游标
func NewTokenCursor(sc *SrcCursor) *TokenCursor {
	return &TokenCursor{sc: sc}
}

func (tc *TokenCursor) Next() *token.Token {

	tc.eatWhitespace()
	startPos := tc.sc.rdPos

	if tc.sc.Eof() {
		return token.NewToken(token.EOF, "", &startPos, &startPos)
	}
	switch ch := tc.sc.Peek(); {
	case IsCommentStart(ch) && strings.ContainsRune("/*", rune(tc.sc.PeekN(2))):
		return tc.eatComment(&startPos)
	case IsIdentStart(ch):
		return tc.eatIdent(&startPos)
	case IsNumberLiteral(ch):
		return tc.eatNumberLiteral(&startPos)
	case token.IsOperatorStart(ch):
		return tc.eatOperator(&startPos)
	}

	// illegal want
	illegalStr := tc.sc.EatWhile(IsWhitespace)
	return token.NewToken(token.ILLEGAL, illegalStr, &startPos, tc.endPos())
}

func (tc *TokenCursor) endPos() token.Pos {
	ep := tc.sc.rdPos
	ep.col--
	return &ep
}

func (tc *TokenCursor) eatWhitespace() {
	tc.sc.EatWhile(IsWhitespace)
}

func (tc *TokenCursor) eatComment(startPos token.Pos) *token.Token {
	tc.sc.Next() // eat '/'

	switch tc.sc.Next() {
	case '/':
		// single-line comment; // xxx
		c := tc.sc.EatWhile(func(c Char) bool { return c != '\n' })
		return token.NewToken(token.COMMENT, c, startPos, tc.endPos())
	case '*':
		// multi line comment; /* xxx */
		var c strings.Builder
		for tc.sc.PeekN(2) != '/' {
			comment := tc.sc.EatWhile(func(c Char) bool { return c != '*' })
			c.WriteString(comment)
		}
		tc.sc.Next() // eat last '/'
		return token.NewToken(token.COMMENT, c.String(), startPos, tc.endPos())

	default:
		// Never reach
		return nil
	}

}

func (tc *TokenCursor) eatIdent(startPos token.Pos) *token.Token {
	idVal := tc.sc.EatWhile(IsIdentBody)
	if tk, ok := token.Keyword(idVal); ok {
		return token.NewToken(tk, idVal, startPos, tc.endPos())
	} else {
		return token.NewToken(token.IDENT, idVal, startPos, tc.endPos())
	}
}

func (tc *TokenCursor) eatNumberLiteral(startPos token.Pos) *token.Token {
	numVal := tc.sc.EatWhile(IsNumberLiteral)
	if strings.ContainsRune(numVal, '.') {
		return token.NewToken(token.FLOAT_LITERAL, numVal, startPos, tc.endPos())
	} else {
		return token.NewToken(token.INT_LITERAL, numVal, startPos, tc.endPos())
	}

}

func (tc *TokenCursor) eatOperator(startPos token.Pos) *token.Token {
	opFirstChar := tc.sc.Next()
	opContinueChars := tc.sc.EatWhile(token.IsOperatorContinue)
	opVal := string(opFirstChar) + opContinueChars
	if tk, ok := token.Operator(opVal); ok {
		return token.NewToken(tk, opVal, startPos, tc.endPos())

	}
	return token.NewToken(token.ILLEGAL, opVal, startPos, tc.endPos())

}

func IsWhitespace(c Char) bool {
	return unicode.IsSpace(rune(c))
}

func IsIdentStart(c Char) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func IsIdentBody(c Char) bool {
	return IsIdentStart(c) || (c >= '0' && c <= '9')
}

func IsNumberLiteral(c Char) bool {
	return (c >= '0' && c <= '9') || c == '.'
}

func IsCommentStart(c Char) bool {
	return c == '/'
}
