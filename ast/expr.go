package ast

import "stmt/token"

type Expr interface {
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
