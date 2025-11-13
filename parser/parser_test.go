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
		{
			name: "assign",
			source: `
			a = "value";
			`,
			want: &ast.Assign{
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "a",
					Line:      2,
					Literal:   nil,
				},
				Value: &ast.Literal{
					Value: "value",
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
			got, err := p.Expression()
			if err != tt.err {
				t.Errorf("Expression() error = %v, want err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression() got = %v, want %v", got, tt.want)
				spew.Dump(got)
				spew.Dump(tt.want)
				return
			}
		})
	}
}

func TestParser_statement(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   ast.Expr
		err    error
	}{
		{
			name:   "print 1;",
			source: "print 1;",
			want: &ast.Print{
				Line: 1,
				Expression: &ast.Literal{
					Value: 1.0,
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
			got, err := p.statement()
			if err != tt.err {
				t.Errorf("statement() error = %v, want err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statement() got = %v, want %v", got, tt.want)
				spew.Dump(got)
				spew.Dump(tt.want)
				return
			}
		})
	}
}

func TestParser_declaration(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   ast.Expr
		err    error
	}{
		// todo
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.New(tt.source)
			tokens := s.ScanTokens()
			p := New(tokens)
			got, err := p.declaration()
			if err != tt.err {
				t.Errorf("declaration() error = %v, want err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("declaration() got = %v, want %v", got, tt.want)
				spew.Dump(got)
				spew.Dump(tt.want)
				return
			}
		})
	}
}
