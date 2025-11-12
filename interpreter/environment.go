package interpreter

import (
	"errors"
	"stmt/token"
)

var (
	ErrUndefinedVariable = errors.New("undefined variable")
)

type environment struct {
	values map[string]any
}

func newEnvironment() *environment {
	return &environment{
		values: make(map[string]any),
	}
}

func (e *environment) define(name string, value any) error {
	e.values[name] = value
	return nil
	// todo 目前 define 是支持重复声明的。error 返回是为了之后升级到不支持重复声明可以返回错误
}

func (e *environment) get(name *token.Token) (any, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, ErrUndefinedVariable
}

func (e *environment) assign(name *token.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	return ErrUndefinedVariable
}
