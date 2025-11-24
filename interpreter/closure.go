package interpreter

import (
	"stmt/ast"
)

type closure struct {
	Function *ast.Function
	Env      *environment
}

func (c *closure) bind(ins *instance) (*closure, error) {
	env := newEnvironment(c.Env)
	err := env.define("this", ins)
	if err != nil {
		return nil, err
	}
	return &closure{
		Function: c.Function,
		Env:      env,
	}, nil
}
