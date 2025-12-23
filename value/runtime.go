package value

import (
	"fmt"
	"io"
)

type Bool struct {
	Literal bool
}

func NewBool(literal bool) *Bool {
	return &Bool{
		Literal: literal,
	}
}

func (b *Bool) String() string {
	return fmt.Sprintf("Bool(%t)", b.Literal)
}

func (b *Bool) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%t\n", b.Literal)
	return err
}

func (b *Bool) ValueType() uint8 {
	return TypeBool
}

func (b *Bool) WriteTo(w io.Writer) error {
	return nil
}

type Nil struct {
}

func NewNil() *Nil {
	return &Nil{}
}

func (n *Nil) String() string {
	return "Nil"
}

func (n *Nil) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "nil\n")
	return err
}

func (n *Nil) ValueType() uint8 {
	return TypeNil
}

func (n *Nil) WriteTo(w io.Writer) error {
	return nil
}

type Closure struct {
	Function *Function
}

func NewClosure(function *Function) *Closure {
	return &Closure{
		Function: function,
	}
}

func (c *Closure) String() string {
	return fmt.Sprintf("Closure(%s)", c.Function.String())
}

func (c *Closure) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "closure\n")
	return err
}

func (c *Closure) ValueType() uint8 {
	return TypeClosure
}

func (c *Closure) WriteTo(w io.Writer) error {
	return nil
}
