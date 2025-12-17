package ast

import "stmt/token"

type Set struct {
	Line   int
	Object Expr
	Name   *token.Token
	Value  Expr
}

func (s *Set) node()    {}
func (s *Set) expr()    {}
func (s *Set) Pos() int { return s.Line }

type Assign struct {
	Line  int
	Name  *token.Token
	Value Expr
}

func (a *Assign) node()    {}
func (a *Assign) expr()    {}
func (a *Assign) Pos() int { return a.Line }

type Logical struct {
	Line     int
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (l *Logical) node()    {}
func (l *Logical) expr()    {}
func (l *Logical) Pos() int { return l.Line }

type Binary struct {
	Line     int
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (b *Binary) node()    {}
func (b *Binary) expr()    {}
func (b *Binary) Pos() int { return b.Line }

type Unary struct {
	Line     int
	Operator *token.Token
	Right    Expr
}

func (u *Unary) node()    {}
func (u *Unary) expr()    {}
func (u *Unary) Pos() int { return u.Line }

type Get struct {
	Line   int
	Object Expr
	Name   *token.Token
}

func (g *Get) node()    {}
func (g *Get) expr()    {}
func (g *Get) Pos() int { return g.Line }

type Call struct {
	Line      int
	Callee    Expr
	Arguments []Expr
}

func (c *Call) node()    {}
func (c *Call) expr()    {}
func (c *Call) Pos() int { return c.Line }

type This struct {
	Line    int
	Keyword *token.Token
}

func (t *This) node()    {}
func (t *This) expr()    {}
func (t *This) Pos() int { return t.Line }

type Super struct {
	Line    int
	Keyword *token.Token
	Method  *token.Token
}

func (s *Super) node()    {}
func (s *Super) expr()    {}
func (s *Super) Pos() int { return s.Line }

type Grouping struct {
	Line       int
	Expression Expr
}

func (g *Grouping) node()    {}
func (g *Grouping) expr()    {}
func (g *Grouping) Pos() int { return g.Line }

type Variable struct {
	Line int
	Name *token.Token
}

func (v *Variable) node()    {}
func (v *Variable) expr()    {}
func (v *Variable) Pos() int { return v.Line }

type Literal struct {
	Line  int
	Value any
}

func (l *Literal) node()    {}
func (l *Literal) expr()    {}
func (l *Literal) Pos() int { return l.Line }
