package ast

import "stmt/token"

type Decl interface {
	Node
}

type Var struct {
	Name        *token.Token
	Initializer Expr
}
