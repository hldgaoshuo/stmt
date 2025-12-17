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
				Line:  1,
				Value: false,
			},
			err: nil,
		},
		{
			name:   "true",
			source: "true",
			want: &ast.Literal{
				Line:  1,
				Value: true,
			},
			err: nil,
		},
		{
			name:   "nil",
			source: "nil",
			want: &ast.Literal{
				Line:  1,
				Value: nil,
			},
			err: nil,
		},
		{
			name:   "123",
			source: "123",
			want: &ast.Literal{
				Line:  1,
				Value: int64(123),
			},
			err: nil,
		},
		{
			name:   "12.3",
			source: "12.3",
			want: &ast.Literal{
				Line:  1,
				Value: 12.3,
			},
			err: nil,
		},
		{
			name:   `"abc"`,
			source: `"abc"`,
			want: &ast.Literal{
				Line:  1,
				Value: "abc",
			},
			err: nil,
		},
		{
			name:   "abc",
			source: "abc",
			want: &ast.Variable{
				Line: 1,
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
				Line: 1,
				Expression: &ast.Literal{
					Line:  1,
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name:   "lucky()",
			source: "lucky()",
			want: &ast.Call{
				Line: 1,
				Callee: &ast.Variable{
					Line: 1,
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
				Line: 1,
				Callee: &ast.Variable{
					Line: 1,
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "lucky",
						Literal:   nil,
						Line:      1,
					},
				},
				Arguments: []ast.Expr{
					&ast.Literal{
						Line:  1,
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
				Line: 1,
				Callee: &ast.Variable{
					Line: 1,
					Name: &token.Token{
						TokenType: token.IDENTIFIER,
						Lexeme:    "lucky",
						Literal:   nil,
						Line:      1,
					},
				},
				Arguments: []ast.Expr{
					&ast.Literal{
						Line:  1,
						Value: int64(123),
					},
					&ast.Literal{
						Line:  1,
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
				Line: 1,
				Operator: &token.Token{
					TokenType: token.MINUS,
					Lexeme:    "-",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name:   "4*2",
			source: "4*2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.STAR,
					Lexeme:    "*",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4/2",
			source: "4/2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.SLASH,
					Lexeme:    "/",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4+2",
			source: "4+2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.PLUS,
					Lexeme:    "+",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4-2",
			source: "4-2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.MINUS,
					Lexeme:    "-",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4>2",
			source: "4>2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.GREATER,
					Lexeme:    ">",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4>=2",
			source: "4>=2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.GREATER_EQUAL,
					Lexeme:    ">=",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4<2",
			source: "4<2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.LESS,
					Lexeme:    "<",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4<=2",
			source: "4<=2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.LESS_EQUAL,
					Lexeme:    "<=",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4==2",
			source: "4==2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.EQUAL_EQUAL,
					Lexeme:    "==",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "4!=2",
			source: "4!=2",
			want: &ast.Binary{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: int64(4),
				},
				Operator: &token.Token{
					TokenType: token.BANG_EQUAL,
					Lexeme:    "!=",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: int64(2),
				},
			},
			err: nil,
		},
		{
			name:   "true and true",
			source: "true and true",
			want: &ast.Logical{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: true,
				},
				Operator: &token.Token{
					TokenType: token.AND,
					Lexeme:    "and",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: true,
				},
			},
			err: nil,
		},
		{
			name:   "true or true",
			source: "true or true",
			want: &ast.Logical{
				Line: 1,
				Left: &ast.Literal{
					Line:  1,
					Value: true,
				},
				Operator: &token.Token{
					TokenType: token.OR,
					Lexeme:    "or",
					Literal:   nil,
					Line:      1,
				},
				Right: &ast.Literal{
					Line:  1,
					Value: true,
				},
			},
			err: nil,
		},
		{
			name:   "abc=123",
			source: "abc=123",
			want: &ast.Assign{
				Line: 1,
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "abc",
					Literal:   nil,
					Line:      1,
				},
				Value: &ast.Literal{
					Line:  1,
					Value: int64(123),
				},
			},
			err: nil,
		},
		{
			name:   "get",
			source: "someObject.someProperty",
			want: &ast.Get{
				Line: 1,
				Object: &ast.Variable{
					Line: 1,
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
				Line: 1,
				Object: &ast.Variable{
					Line: 1,
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
					Line: 1,
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
			tokens := s.Scan()
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
					Line:  1,
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
					Line:  1,
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
					Line: 2,
					Left: &ast.Variable{
						Line: 2,
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
						Line:  2,
						Value: int64(2),
					},
				},
				Body: &ast.Block{
					Line: 2,
					Declarations: []ast.Stmt{
						&ast.ExpressionStatement{
							Line: 3,
							Expression: &ast.Literal{
								Line:  3,
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
					Line: 2,
					Left: &ast.Variable{
						Line: 2,
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
						Line:  2,
						Value: int64(2),
					},
				},
				ThenBranch: &ast.Block{
					Line: 2,
					Declarations: []ast.Stmt{
						&ast.ExpressionStatement{
							Line: 3,
							Expression: &ast.Literal{
								Line:  3,
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
					Line: 2,
					Left: &ast.Variable{
						Line: 2,
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
						Line:  2,
						Value: int64(2),
					},
				},
				ThenBranch: &ast.Block{
					Line: 2,
					Declarations: []ast.Stmt{
						&ast.ExpressionStatement{
							Line: 3,
							Expression: &ast.Literal{
								Line:  3,
								Value: int64(123),
							},
						},
					},
				},
				ElseBranch: &ast.Block{
					Line: 4,
					Declarations: []ast.Stmt{
						&ast.ExpressionStatement{
							Line: 5,
							Expression: &ast.Literal{
								Line:  5,
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
			tokens := s.Scan()
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
			name: "function",
			source: `
			fun breakfast() {
				print "Eggs a-fryin'!";
			}
			`,
			want: &ast.Function{
				Line: 2,
				Name: &token.Token{
					TokenType: token.IDENTIFIER,
					Lexeme:    "breakfast",
					Line:      2,
					Literal:   nil,
				},
				Params: nil,
				Body: &ast.Block{
					Line: 2,
					Declarations: []ast.Stmt{
						&ast.Print{
							Line: 3,
							Expression: &ast.Literal{
								Line:  3,
								Value: "Eggs a-fryin'!",
							},
						},
					},
				},
			},
			err: nil,
		},
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
							Declarations: []ast.Stmt{
								&ast.Print{
									Line: 4,
									Expression: &ast.Literal{
										Line:  4,
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
							Declarations: []ast.Stmt{
								&ast.Print{
									Line: 4,
									Expression: &ast.Literal{
										Line:  4,
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
							Declarations: []ast.Stmt{
								&ast.Print{
									Line: 7,
									Expression: &ast.Binary{
										Line: 7,
										Left: &ast.Binary{
											Line: 7,
											Left: &ast.Literal{
												Line:  7,
												Value: "Enjoy your breakfast, ",
											},
											Operator: &token.Token{
												TokenType: token.PLUS,
												Lexeme:    "+",
												Line:      7,
											},
											Right: &ast.Variable{
												Line: 7,
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
											Line:  7,
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
			tokens := s.Scan()
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
