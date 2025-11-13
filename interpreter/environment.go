package interpreter

import (
	"stmt/token"
)

type environment struct {
	enclosing *environment
	values    map[string]any
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *environment) define(name string, value any) error {
	e.values[name] = value
	return nil
	// todo 目前 define 是支持重复声明的。error 返回是为了之后升级到不支持重复声明可以返回错误
}

func (e *environment) get(name *token.Token) (any, error) {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value, nil
	} else {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		} else {
			return nil, ErrUndefinedVariable
		}
	}
}

func (e *environment) assign(name *token.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil
	} else {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		} else {
			return ErrUndefinedVariable
		}
	}
}
