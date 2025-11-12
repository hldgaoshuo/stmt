package ast

import "stmt/token"

type Expr interface {
	Node
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator *token.Token
	Right    Expr
}

type Variable struct {
	Name *token.Token
}

type Assign struct {
	Name  *token.Token
	Value Expr
}
