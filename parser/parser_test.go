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

func TestParser_Expression(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   ast.Node
		err    error
	}{
		{
			name:   "false",
			source: "false",
			want: &ast.Literal{
				Value: false,
			},
			err: nil,
		},
		{
			name:   "true",
			source: "true",
			want: &ast.Literal{
				Value: true,
			},
			err: nil,
		},
		{
			name:   "nil",
			source: "nil",
			want: &ast.Literal{
				Value: nil,
			},
			err: nil,
		},
		{
			name:   "123",
			source: "123",
			want: &ast.Literal{
				Value: int64(123),
			},
			err: nil,
		},
		{
			name:   "12.3",
			source: "12.3",
			want: &ast.Literal{
				Value: 12.3,
			},
			err: nil,
		},
		{
			name:   `"abc"`,
			source: `"abc"`,
			want: &ast.Literal{
				Value: "abc",
			},
			err: nil,
		},
		{
			name:   "abc",
			source: "abc",
			want: &ast.Variable{
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "abc",
					Literal:   nil,
					Line:      1,
				},
			},
			err: nil,
		},
		{
			name:   "(123)",
			source: "(123)",
			want: &ast.Grouping{
				Expression: &ast.Literal{
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name:   "lucky()",
			source: "lucky()",
			want: &ast.Call{
				Callee: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "lucky",
						Literal:   nil,
						Line:      1,
					},
				},
				Arguments: nil,
			},
			err: nil,
		},
		{
			name:   "lucky(123)",
			source: "lucky(123)",
			want: &ast.Call{
				Callee: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "lucky",
						Literal:   nil,
						Line:      1,
					},
				},
				Arguments: []ast.Node{
					&ast.Literal{
						Value: int64(123),
					},
				},
			},
			err: nil,
		},
		{
			name:   "lucky(123, 12.3)",
			source: "lucky(123, 12.3)",
			want: &ast.Call{
				Callee: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "lucky",
						Literal:   nil,
						Line:      1,
					},
				},
				Arguments: []ast.Node{
					&ast.Literal{
						Value: int64(123),
					},
					&ast.Literal{
						Value: 12.3,
					},
				},
			},
			err: nil,
		},
		{
			name:   "-123",
			source: "-123",
			want: &ast.Unary{
				Operator: &token.Token{
					TokenType: token.MINUS,
					Lexeme:    "-",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name:   "4*2",
			source: "4*2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.STAR,
					Lexeme:    "*",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4/2",
			source: "4/2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.SLASH,
					Lexeme:    "/",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4+2",
			source: "4+2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.PLUS,
					Lexeme:    "+",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4-2",
			source: "4-2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.MINUS,
					Lexeme:    "-",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4>2",
			source: "4>2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.GREATER,
					Lexeme:    ">",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4>=2",
			source: "4>=2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.GREATER_EQUAL,
					Lexeme:    ">=",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4<2",
			source: "4<2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.LESS,
					Lexeme:    "<",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4<=2",
			source: "4<=2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.LESS_EQUAL,
					Lexeme:    "<=",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4==2",
			source: "4==2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.EQUAL_EQUAL,
					Lexeme:    "==",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4!=2",
			source: "4!=2",
			want: &ast.Binary{
				Left: &ast.Literal{
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.BANG_EQUAL,
					Lexeme:    "!=",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "true and true",
			source: "true and true",
			want: &ast.Logical{
				Left: &ast.Literal{
					Value: true,
				},
				Operator: &token.Token{
					TokenType: token.AND,
					Lexeme:    "and",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: true,
				},
			},
			err: nil,
		},
		{
			name:   "true or true",
			source: "true or true",
			want: &ast.Logical{
				Left: &ast.Literal{
					Value: true,
				},
				Operator: &token.Token{
					TokenType: token.OR,
					Lexeme:    "or",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Value: true,
				},
			},
			err: nil,
		},
		{
			name:   "abc=123",
			source: "abc=123",
			want: &ast.Assign{
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "abc",
					Literal:   nil,
					Line:      1,
				},
				Value: &ast.Literal{
					Value: int64(123),
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
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expression() result")
				spew.Dump(got)
				t.Errorf("want")
				spew.Dump(tt.want)
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
			name:   "print",
			source: "print 123;",
			want: &ast.Print{
				Line: 1,
				Expression: &ast.Literal{
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name:   "return",
			source: "return 123;",
			want: &ast.Return{
				Line: 1,
				Expression: &ast.Literal{
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name: "while",
			source: `
			while (i<2) {
				123;
			}
			`,
			want: &ast.While{
				Line: 2,
				Condition: &ast.Binary{
					Left: &ast.Variable{
						Name: &token.Token{
							TokenType: token.IDENTIFIER,
							Lexeme:    "i",
							Literal:   nil,
							Line:      2,
						},
					},
					Operator: &token.Token{
						TokenType: token.LESS,
						Lexeme:    "<",
						Literal:   nil,
						Line:      2,
					},
					Right: &ast.Literal{
						Value: int64(2),
					},
				},
				Body: &ast.Block{
					Line: 2,
					Declarations: []ast.Node{
						&ast.ExpressionStatement{
							Expression: &ast.Literal{
								Value: int64(123),
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "if",
			source: `
			if (i<2) {
				123;
			}
			`,
			want: &ast.If{
				Line: 2,
				Condition: &ast.Binary{
					Left: &ast.Variable{
						Name: &token.Token{
							TokenType: token.IDENTIFIER,
							Lexeme:    "i",
							Literal:   nil,
							Line:      2,
						},
					},
					Operator: &token.Token{
						TokenType: token.LESS,
						Lexeme:    "<",
						Literal:   nil,
						Line:      2,
					},
					Right: &ast.Literal{
						Value: int64(2),
					},
				},
				ThenBranch: &ast.Block{
					Line: 2,
					Declarations: []ast.Node{
						&ast.ExpressionStatement{
							Expression: &ast.Literal{
								Value: int64(123),
							},
						},
					},
				},
				ElseBranch: nil,
			},
			err: nil,
		},
		{
			name: "if 2",
			source: `
			if (i<2) {
				123;
			} else {
				"abc";
			}
			`,
			want: &ast.If{
				Line: 2,
				Condition: &ast.Binary{
					Left: &ast.Variable{
						Name: &token.Token{
							TokenType: token.IDENTIFIER,
							Lexeme:    "i",
							Literal:   nil,
							Line:      2,
						},
					},
					Operator: &token.Token{
						TokenType: token.LESS,
						Lexeme:    "<",
						Literal:   nil,
						Line:      2,
					},
					Right: &ast.Literal{
						Value: int64(2),
					},
				},
				ThenBranch: &ast.Block{
					Line: 2,
					Declarations: []ast.Node{
						&ast.ExpressionStatement{
							Expression: &ast.Literal{
								Value: int64(123),
							},
						},
					},
				},
				ElseBranch: &ast.Block{
					Line: 4,
					Declarations: []ast.Node{
						&ast.ExpressionStatement{
							Expression: &ast.Literal{
								Value: "abc",
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
			got, err := p.statement()
			if !errors.Is(err, tt.err) {
				t.Errorf("statement() error = %v, want err %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Statement() result")
				spew.Dump(got)
				t.Errorf("want")
				spew.Dump(tt.want)
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
		{
			name: "class 4",
			source: `
			class BostonCream < Doughnut {
			}
			`,
			want: &ast.Class{
				Line: 2,
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "BostonCream",
					Line:      2,
					Literal:   nil,
				},
				Methods: nil,
				SuperClass: &ast.Variable{
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "Doughnut",
						Line:      2,
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
