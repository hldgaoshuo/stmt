package interpreter

import "stmt/token"

type class struct {
	SuperClass *class
	Closures   []*closure
}

func (c *class) get(name *token.Token) *closure {
	for _, clo := range c.Closures {
		if clo.Function.Name.Lexeme == name.Lexeme {
			return clo
		}
	}
	if c.SuperClass != nil {
		return c.SuperClass.get(name)
	} else {
		return nil
	}
}
