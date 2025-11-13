package ast

import "stmt/token"

type Var struct {
	Line        int
	Name        *token.Token
	Initializer Node
}
