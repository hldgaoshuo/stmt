package interpreter

import (
	"stmt/token"
)

type environment struct {
	Enclosing *environment
	Values    map[string]any
}

func newEnvironment(enclosing *environment) *environment {
	return &environment{
		Enclosing: enclosing,
		Values:    make(map[string]any),
	}
}

func (e *environment) define(name string, value any) error {
	e.Values[name] = value
	return nil
	// todo 目前 define 是支持重复声明的。error 返回是为了之后升级到不支持重复声明可以返回错误
}

func (e *environment) get(name *token.Token) (any, error) {
	value, ok := e.Values[name.Lexeme]
	if ok {
		return value, nil
	} else {
		if e.Enclosing != nil {
			return e.Enclosing.get(name)
		} else {
			return nil, ErrUndefinedVariable
		}
	}
}

func (e *environment) assign(name *token.Token, value any) error {
	_, ok := e.Values[name.Lexeme]
	if ok {
		e.Values[name.Lexeme] = value
		return nil
	} else {
		if e.Enclosing != nil {
			return e.Enclosing.assign(name, value)
		} else {
			return ErrUndefinedVariable
		}
	}
}
