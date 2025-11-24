package ast

import (
	"stmt/token"
)

type Var struct {
	Line        int
	Name        *token.Token
	Initializer Node
}

type Function struct {
	Line   int
	Name   *token.Token
	Params []*token.Token
	Body   Node
}

type Class struct {
	Line    int
	Name    *token.Token
	Methods []*Function
}
