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

type Call struct {
	Callee    Node
	Paren     *token.Token
	Arguments []Node
}

type Get struct {
	Object Node
	Name   *token.Token
}

type Set struct {
	Object Node
	Name   *token.Token
	Value  Node
}

type This struct {
	Keyword *token.Token
}

type Super struct {
	Keyword *token.Token
	Method  *token.Token
}
