package token

import (
	"strings"
	"testing"
)

func TestTokenTypeIsKeyword(t *testing.T) {
	tests := []struct {
		tk        Kind
		isKeyword bool
	}{
		{
			tk:        CONST,
			isKeyword: true,
		},
		{
			tk:        IF,
			isKeyword: true,
		},
		{
			tk:        MUL_ASSIGN,
			isKeyword: false,
		},
	}

	for _, testcase := range tests {
		want := testcase.isKeyword
		got := testcase.tk.IsKeyword()
		if want != got {
			t.Errorf("%s.IsKeyword() want %t; got %t.",
				strings.ToUpper(testcase.tk.String()), want, got)
		}
	}
}
