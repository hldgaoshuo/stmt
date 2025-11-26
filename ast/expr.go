package ast

import "stmt/token"

type Set struct {
	Object Node
	Name   *token.Token
	Value  Node
}

type Assign struct {
	Name  *token.Token
	Value Node
}

type Logical struct {
	Left     Node
	Operator *token.Token
	Right    Node
}

type Binary struct {
	Left     Node
	Operator *token.Token
	Right    Node
}

type Unary struct {
	Operator *token.Token
	Right    Node
}

type Get struct {
	Object Node
	Name   *token.Token
}

type Call struct {
	Callee    Node
	Arguments []Node
}

type This struct {
	Keyword *token.Token
}

type Super struct {
	Keyword *token.Token
	Method  *token.Token
}

type Grouping struct {
	Expression Node
}

type Variable struct {
	Name *token.Token
}

type Literal struct {
	Value any
}
