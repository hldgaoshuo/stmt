package ast

import (
	"stmt/token"
)

type Var struct {
	Line        int
	Name        *token.Token
	Initializer Expr
}

func (v *Var) node()    {}
func (v *Var) stmt()    {}
func (v *Var) Pos() int { return v.Line }

type Function struct {
	Line   int
	Name   *token.Token
	Params []*token.Token
	Body   *Block
}

func (f *Function) node()    {}
func (f *Function) stmt()    {}
func (f *Function) Pos() int { return f.Line }

type Class struct {
	Line       int
	Name       *token.Token
	Methods    []*Function
	SuperClass *Variable
}

func (c *Class) node()    {}
func (c *Class) stmt()    {}
func (c *Class) Pos() int { return c.Line }
