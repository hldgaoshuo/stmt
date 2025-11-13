package ast

import "stmt/token"

type Binary struct {
	Left     Node
	Operator *token.Token
	Right    Node
}

type Grouping struct {
	Expression Node
}

type Literal struct {
	Value any
}

type Logical struct {
	Left     Node
	Operator *token.Token
	Right    Node
}

type Unary struct {
	Operator *token.Token
	Right    Node
}

type Variable struct {
	Name *token.Token
}

type Assign struct {
	Name  *token.Token
	Value Node
}
