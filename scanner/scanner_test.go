package scanner

import (
	"fmt"
	"reflect"
	"stmt/token"
	"strings"
	"testing"
)

func formatTokens(tokens []*token.Token) string {
	var result strings.Builder
	for i, tok := range tokens {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(fmt.Sprintf(
			"{TokenType: %v, Lexeme: %q, Literal: %v, Line: %d}",
			tok.TokenType, tok.Lexeme, tok.Literal, tok.Line))
	}
	return "[" + result.String() + "]"
}

func TestScanner_ScanTokens(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   []*token.Token
	}{
		{
			name:   "white",
			source: " \r \t \n ",
			want: []*token.Token{
				token.New(token.EOF, "", nil, 2),
			},
		},
		{
			name:   "comment",
			source: "// this is a comment",
			want: []*token.Token{
				token.New(token.EOF, "", nil, 1),
			},
		},
		{
			name: "comment and single-character tokens",
			source: `// this is a comment
				(( )){} // grouping stuff
				`,
			want: []*token.Token{
				token.New(token.LEFT_PAREN, "(", nil, 2),
				token.New(token.LEFT_PAREN, "(", nil, 2),
				token.New(token.RIGHT_PAREN, ")", nil, 2),
				token.New(token.RIGHT_PAREN, ")", nil, 2),
				token.New(token.LEFT_BRACE, "{", nil, 2),
				token.New(token.RIGHT_BRACE, "}", nil, 2),
				token.New(token.EOF, "", nil, 3),
			},
		},
		{
			name: "comment, one or two character tokens and single-character tokens",
			source: `// this is a comment
				(( )){} // grouping stuff
				!*+-/=<> <= == // operators
				`,
			want: []*token.Token{
				token.New(token.LEFT_PAREN, "(", nil, 2),
				token.New(token.LEFT_PAREN, "(", nil, 2),
				token.New(token.RIGHT_PAREN, ")", nil, 2),
				token.New(token.RIGHT_PAREN, ")", nil, 2),
				token.New(token.LEFT_BRACE, "{", nil, 2),
				token.New(token.RIGHT_BRACE, "}", nil, 2),
				token.New(token.BANG, "!", nil, 3),
				token.New(token.STAR, "*", nil, 3),
				token.New(token.PLUS, "+", nil, 3),
				token.New(token.MINUS, "-", nil, 3),
				token.New(token.SLASH, "/", nil, 3),
				token.New(token.EQUAL, "=", nil, 3),
				token.New(token.LESS, "<", nil, 3),
				token.New(token.GREATER, ">", nil, 3),
				token.New(token.LESS_EQUAL, "<=", nil, 3),
				token.New(token.EQUAL_EQUAL, "==", nil, 3),
				token.New(token.EOF, "", nil, 4),
			},
		},
		{
			name:   "identifier",
			source: "gaoshuo",
			want: []*token.Token{
				token.New(token.IDENTIFIER, "gaoshuo", nil, 1),
				token.New(token.EOF, "", nil, 1),
			},
		},
		{
			name:   "int literal",
			source: "1234",
			want: []*token.Token{
				token.New(token.INT_LITERAL, "1234", int64(1234), 1),
				token.New(token.EOF, "", nil, 1),
			},
		},
		{
			name:   "float literal",
			source: "12.34",
			want: []*token.Token{
				token.New(token.FLOAT_LITERAL, "12.34", 12.34, 1),
				token.New(token.EOF, "", nil, 1),
			},
		},
		{
			name:   "string literal",
			source: `"abc"`,
			want: []*token.Token{
				token.New(token.STRING_LITERAL, `"abc"`, "abc", 1),
				token.New(token.EOF, "", nil, 1),
			},
		},
		{
			name:   "keyword and",
			source: "and",
			want: []*token.Token{
				token.New(token.AND, "and", nil, 1),
				token.New(token.EOF, "", nil, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.source)
			if got := s.ScanTokens(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" \n ScanTokens() \n %v \n want \n %v \n", formatTokens(got), formatTokens(tt.want))
			}
		})
	}
}
