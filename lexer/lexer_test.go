package lexer

import (
	"reflect"
	"sometimes/token"
	"strings"
	"testing"
)

func TestTokenCursorNext(t *testing.T) {
	startPos := &SrcPos{0, 0}
	tests := []struct {
		src  string
		want *token.Token
	}{
		{
			src:  "// 你麻痹",
			want: token.NewToken(token.COMMENT, " 你麻痹", startPos, &SrcPos{0, 11}),
		},
		{
			src: `/* 我操
真的牛逼i啊
🐂*/`,
			want: token.NewToken(token.COMMENT, ` 我操
真的牛逼i啊
🐂`, startPos, &SrcPos{2, 4}),
		},
		{
			src:  "T",
			want: token.NewToken(token.IDENT, "T", startPos, &SrcPos{0, 0}),
		},
		{
			src:  "taoyu",
			want: token.NewToken(token.IDENT, "taoyu", startPos, &SrcPos{0, 4}),
		},
		{
			src:  "tao_Yu8",
			want: token.NewToken(token.IDENT, "tao_Yu8", startPos, &SrcPos{0, 6}),
		},
		{
			src:  "1taoyu",
			want: token.NewToken(token.INT_LITERAL, "1", startPos, &SrcPos{0, 0}),
		},
		{
			src:  "if",
			want: token.NewToken(token.IF, "if", startPos, &SrcPos{0, 1}),
		},
		{
			src:  "314159",
			want: token.NewToken(token.INT_LITERAL, "314159", startPos, &SrcPos{0, 5}),
		},
		{
			src:  "3.14159",
			want: token.NewToken(token.FLOAT_LITERAL, "3.14159", startPos, &SrcPos{0, 6}),
		},
		{
			src:  ".123",
			want: token.NewToken(token.FLOAT_LITERAL, ".123", startPos, &SrcPos{0, 3}),
		},
		{
			src:  "8",
			want: token.NewToken(token.INT_LITERAL, "8", startPos, &SrcPos{0, 0}),
		},
		{
			src:  "&&",
			want: token.NewToken(token.LAND, "&&", startPos, &SrcPos{0, 1}),
		},
		{
			src:  "+",
			want: token.NewToken(token.ADD, "+", startPos, &SrcPos{0, 0}),
		},
		{
			src:  "&",
			want: token.NewToken(token.ILLEGAL, "&", startPos, &SrcPos{0, 0}),
		},
	}

	for _, testcase := range tests {
		tc := NewTokenCursor(NewSrcCursor([]byte(testcase.src)))
		got := tc.Next()
		if !reflect.DeepEqual(testcase.want, got) {
			t.Errorf("\n`%s`\n fisrt want want %s; \n              got %s",
				testcase.src, testcase.want.String(), got.String())
		}
	}
}

func TestLexer(t *testing.T) {
	tests := []struct {
		src  string
		want []*token.Token
	}{
		{
			src: `
fn main() -> int {
    return a+b;
}
`,
			want: []*token.Token{
				token.NewToken(token.IDENT, "fn", &SrcPos{1, 0}, &SrcPos{1, 1}),
				token.NewToken(token.IDENT, "main", &SrcPos{1, 3}, &SrcPos{1, 6}),
				token.NewToken(token.LPAREN, "(", &SrcPos{1, 7}, &SrcPos{1, 7}),
				token.NewToken(token.RPAREN, ")", &SrcPos{1, 8}, &SrcPos{1, 8}),
				token.NewToken(token.ARROW, "->", &SrcPos{1, 10}, &SrcPos{1, 11}),
				token.NewToken(token.IDENT, "int", &SrcPos{1, 13}, &SrcPos{1, 15}),
				token.NewToken(token.LBRACE, "{", &SrcPos{1, 17}, &SrcPos{1, 17}),
				token.NewToken(token.RETURN, "return", &SrcPos{2, 4}, &SrcPos{2, 9}),
				token.NewToken(token.IDENT, "a", &SrcPos{2, 11}, &SrcPos{2, 11}),
				token.NewToken(token.ADD, "+", &SrcPos{2, 12}, &SrcPos{2, 12}),
				token.NewToken(token.IDENT, "b", &SrcPos{2, 13}, &SrcPos{2, 13}),
				token.NewToken(token.SEMICOLON, ";", &SrcPos{2, 14}, &SrcPos{2, 14}),
				token.NewToken(token.RBRACE, "}", &SrcPos{3, 0}, &SrcPos{3, 0}),
			},
		},
	}

	for _, testcase := range tests {
		tc := NewTokenCursor(NewSrcCursor([]byte(testcase.src)))
		var got []*token.Token
		for next := tc.Next(); next.Kind != token.EOF; next = tc.Next() {
			got = append(got, next)
		}

		if !reflect.DeepEqual(testcase.want, got) {
			t.Errorf("\n`%s`\n want want:\n%s;\ngot want:\n%s", testcase.src,
				tokensToString(testcase.want...), tokensToString(got...))
		}

	}

}

func tokensToString(tokens ...*token.Token) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	for _, t := range tokens {
		sb.WriteString(t.String())
		sb.WriteString(",\n")
	}
	sb.WriteRune(']')
	return sb.String()
}
