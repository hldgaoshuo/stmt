package ast

import "stmt/token"

type Decl interface {
	Node
}

type Var struct {
	Line        int
	Name        *token.Token
	Initializer Expr
}
