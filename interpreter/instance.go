package interpreter

import (
	"stmt/ast"
	"stmt/token"
)

type instance struct {
	Class  *ast.Class
	Fields map[string]any
}

func newInstance(class *ast.Class) *instance {
	return &instance{
		Class:  class,
		Fields: make(map[string]any),
	}
}

func (i *instance) get(name *token.Token) (any, error) {
	value, ok := i.Fields[name.Lexeme]
	if ok {
		return value, nil
	} else {
		print("Undefined property '" + name.Lexeme + "'.")
		return nil, ErrUndefinedProperty
	}
}

func (i *instance) set(name *token.Token, value any) {
	i.Fields[name.Lexeme] = value
}
