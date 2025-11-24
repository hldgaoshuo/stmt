package interpreter

import "stmt/token"

type class struct {
	Closures []*closure
}

func (c *class) get(name *token.Token) *closure {
	for _, clo := range c.Closures {
		if clo.Function.Name.Lexeme == name.Lexeme {
			return clo
		}
	}
	return nil
}
