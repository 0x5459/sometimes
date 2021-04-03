package token

import (
	"fmt"
	"strings"
)

type (
	Kind int
	Pos  interface {
		Line() int
		Col() int
	}
	Token struct {
		Kind             Kind
		Val              string
		StartPos, EndPos Pos
	}
)

/// NewToken 返回一个新的 Token.
func NewToken(tk Kind, val string, startPos, endPos Pos) *Token {
	return &Token{
		Kind:     tk,
		Val:      val,
		StartPos: startPos,
		EndPos:   endPos,
	}
}

func (token *Token) String() string {
	return fmt.Sprintf(
		"{Kind: '%s', Val: '%s', StartPos: {Line: %d, Col: %d}, EndPos: {Line: %d, Col: %d}}",
		strings.ToUpper(token.Kind.String()),
		token.Val,
		token.StartPos.Line(), token.StartPos.Col(),
		token.EndPos.Line(), token.EndPos.Col(),
	)
}

// Token 类型列表
const (
	// Special tokens
	ILLEGAL Kind = iota
	EOF
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT          // main
	INT_LITERAL    // 12345
	FLOAT_LITERAL  // 123.45
	CHAR_LITERAL   // 'a'
	STRING_LITERAL // "abc"
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	LAND  // &&
	LOR   // ||
	ARROW // ->
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ    // !=
	LEQ    // <=
	GEQ    // >=
	DEFINE // :=

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operator_end

	keyword_beg
	// Keywords
	BREAK
	CASE
	CONST
	CONTINUE

	DEFAULT
	ELSE
	FOR

	FN
	IF
	IMPORT

	RETURN

	STRUCT
	SWITCH
	TYPE
	LET
	keyword_end
)

var (
	tokens = [...]string{
		ILLEGAL: "ILLEGAL",

		EOF:     "EOF",
		COMMENT: "COMMENT",

		IDENT:          "IDENT",
		INT_LITERAL:    "INT",
		FLOAT_LITERAL:  "FLOAT",
		CHAR_LITERAL:   "CHAR",
		STRING_LITERAL: "STRING",

		ADD: "+",
		SUB: "-",
		MUL: "*",
		QUO: "/",
		REM: "%",

		ADD_ASSIGN: "+=",
		SUB_ASSIGN: "-=",
		MUL_ASSIGN: "*=",
		QUO_ASSIGN: "/=",
		REM_ASSIGN: "%=",

		LAND:  "&&",
		LOR:   "||",
		ARROW: "->",
		INC:   "++",
		DEC:   "--",

		EQL:    "==",
		LSS:    "<",
		GTR:    ">",
		ASSIGN: "=",
		NOT:    "!",

		NEQ:    "!=",
		LEQ:    "<=",
		GEQ:    ">=",
		DEFINE: ":=",

		LPAREN: "(",
		LBRACK: "[",
		LBRACE: "{",
		COMMA:  ",",
		PERIOD: ".",

		RPAREN:    ")",
		RBRACK:    "]",
		RBRACE:    "}",
		SEMICOLON: ";",
		COLON:     ":",

		BREAK:    "break",
		CASE:     "case",
		CONST:    "const",
		CONTINUE: "continue",

		DEFAULT: "default",
		ELSE:    "else",
		FOR:     "for",

		FN:     "func",
		IF:     "if",
		IMPORT: "import",

		RETURN: "return",

		STRUCT: "struct",
		SWITCH: "switch",
		TYPE:   "type",
		LET:    "let",
	}

	keywords  = make(map[string]Kind)
	operators = make(map[string]Kind)
)

func init() {

	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}

	for i := operator_beg + 1; i < operator_end; i++ {
		operators[tokens[i]] = i
	}
}

func (tk Kind) String() string {
	return tokens[tk]
}

/// 如果'tk'是一个关键字,则返回正确.
func (tk Kind) IsKeyword() bool {
	return tk > keyword_beg && tk < keyword_end
}

/// 如果'tk'是一个运算符,则返回正确.
func (tk Kind) IsOperator() bool {
	return tk > operator_beg && tk < operator_end
}

/// 如果'tk'是一个字面量,则返回正确.
func (tk Kind) IsLiteral() bool {
	return tk > literal_beg && tk < literal_end
}

func Operator(t string) (Kind, bool) {
	tk, ok := operators[t]
	return tk, ok
}

func Keyword(ident string) (Kind, bool) {
	tk, ok := keywords[ident]
	return tk, ok
}

func IsOperatorStart(c byte) bool {
	for op := range operators {
		if c == op[0] {
			return true
		}
	}
	return false
}

func IsOperatorContinue(c byte) bool {
	for op := range operators {
		if strings.ContainsRune(op[1:], rune(c)) {
			return true
		}
	}
	return false
}
