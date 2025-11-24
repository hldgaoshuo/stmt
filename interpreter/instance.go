package interpreter

import (
	"stmt/token"
)

type instance struct {
	Class  *class
	Fields map[string]any
}

func newInstance(class *class) *instance {
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
		clo := i.Class.get(name)
		if clo != nil {
			clo, err := clo.bind(i)
			if err != nil {
				return nil, err
			}
			return clo, nil
		} else {
			print("Undefined property '" + name.Lexeme + "'.")
			return nil, ErrUndefinedProperty
		}
	}
}

func (i *instance) set(name *token.Token, value any) {
	i.Fields[name.Lexeme] = value
}
