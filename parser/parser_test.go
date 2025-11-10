package parser

import (
	"reflect"
	"stmt/ast"
	"stmt/scanner"
	"stmt/token"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestParser_expression(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   ast.Expr
		err    error
	}{
		{
			name:   "1+2",
			source: "1+2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: 1.0,
				},
				Operator: &token.Token{
					TokenType: token.PLUS,
					Lexeme:    "+",
					Line:      1,
					Literal:   nil,
				},
				Right: &ast.Literal{
					Value: 2.0,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.New(tt.source)
			tokens := s.ScanTokens()
			p := New(tokens)
			got, err := p.expression()
			if err != tt.err {
				t.Errorf("expression() error = %v, want err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expression() got = %v, want %v", got, tt.want)
				spew.Dump(got)
				spew.Dump(tt.want)
				return
			}
		})
	}
}
