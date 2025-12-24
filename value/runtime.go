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

func (b *Bool) GetLiteral() any {
	return b.Literal
}

func (b *Bool) SetLiteral(literal any) {
	b.Literal = literal.(bool)
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

func (n *Nil) GetLiteral() any {
	return nil
}

func (n *Nil) SetLiteral(literal any) {
}

type Closure struct {
	Function *Function
	Upvalues []Value
}

func NewClosure(function *Function) *Closure {
	return &Closure{
		Function: function,
		Upvalues: make([]Value, function.NumUpvalues),
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

func (c *Closure) GetLiteral() any {
	panic("closure have no literal")
}

func (c *Closure) SetLiteral(literal any) {
	panic("closure have no literal")
}
