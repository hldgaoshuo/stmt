package parser

import (
	"errors"
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
		want   ast.Node
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
		{
			name:   "get",
			source: "someObject.someProperty",
			want: &ast.Get{
				Object: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "someObject",
						Line:      1,
						Literal:   nil,
					},
				},
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "someProperty",
					Line:      1,
					Literal:   nil,
				},
			},
			err: nil,
		},
		{
			name:   "set",
			source: "someObject.someProperty = value;",
			want: &ast.Set{
				Object: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "someObject",
						Line:      1,
						Literal:   nil,
					},
				},
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "someProperty",
					Line:      1,
					Literal:   nil,
				},
				Value: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "value",
						Line:      1,
						Literal:   nil,
					},
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
			if !errors.Is(err, tt.err) {
				t.Errorf("Expression() error = %v, want err %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression() DeepEqual error")
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
		want   ast.Node
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
			if !errors.Is(err, tt.err) {
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
		want   ast.Node
		err    error
	}{
		{
			name: "class",
			source: `
			class Breakfast {
			}
			`,
			want: &ast.Class{
				Line: 2,
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "Breakfast",
					Line:      2,
					Literal:   nil,
				},
				Methods: nil,
			},
			err: nil,
		},
		{
			name: "class 2",
			source: `
			class Breakfast {
				cook() {
					print "Eggs a-fryin'!";
				}
			}
			`,
			want: &ast.Class{
				Line: 2,
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "Breakfast",
					Line:      2,
					Literal:   nil,
				},
				Methods: []*ast.Function{
					{
						Line: 3,
						Name: &token.Token{
							TokenType: token.IDENTIFIER,
							Lexeme:    "cook",
							Line:      3,
							Literal:   nil,
						},
						Params: nil,
						Body: &ast.Block{
							Line: 3,
							Declarations: []ast.Node{
								&ast.Print{
									Line: 4,
									Expression: &ast.Literal{
										Value: "Eggs a-fryin'!",
									},
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "class 3",
			source: `
			class Breakfast {
				cook() {
					print "Eggs a-fryin'!";
				}
				serve(who) {
					print "Enjoy your breakfast, " + who + ".";
				}
			}
			`,
			want: &ast.Class{
				Line: 2,
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "Breakfast",
					Line:      2,
					Literal:   nil,
				},
				Methods: []*ast.Function{
					{
						Line: 3,
						Name: &token.Token{
							TokenType: token.IDENTIFIER,
							Lexeme:    "cook",
							Line:      3,
							Literal:   nil,
						},
						Params: nil,
						Body: &ast.Block{
							Line: 3,
							Declarations: []ast.Node{
								&ast.Print{
									Line: 4,
									Expression: &ast.Literal{
										Value: "Eggs a-fryin'!",
									},
								},
							},
						},
					},
					{
						Line: 6,
						Name: &token.Token{
							TokenType: token.IDENTIFIER,
							Lexeme:    "serve",
							Line:      6,
							Literal:   nil,
						},
						Params: []*token.Token{
							{
								TokenType: token.IDENTIFIER,
								Lexeme:    "who",
								Line:      6,
								Literal:   nil,
							},
						},
						Body: &ast.Block{
							Line: 6,
							Declarations: []ast.Node{
								&ast.Print{
									Line: 7,
									Expression: &ast.Binary{
										Left: &ast.Binary{
											Left: &ast.Literal{
												Value: "Enjoy your breakfast, ",
											},
											Operator: &token.Token{
												TokenType: token.PLUS,
												Lexeme:    "+",
												Line:      7,
											},
											Right: &ast.Variable{
												Name: &token.Token{
													TokenType: token.IDENTIFIER,
													Lexeme:    "who",
													Line:      7,
													Literal:   nil,
												},
											},
										},
										Operator: &token.Token{
											TokenType: token.PLUS,
											Lexeme:    "+",
											Line:      7,
										},
										Right: &ast.Literal{
											Value: ".",
										},
									},
								},
							},
						},
					},
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
			got, err := p.declaration()
			if !errors.Is(err, tt.err) {
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
