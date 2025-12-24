package value

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Value interface {
	String() string
	Print(w io.Writer) error
	ValueType() uint8
	WriteTo(w io.Writer) error
	GetLiteral() any
	SetLiteral(literal any)
}

const (
	TypeInt uint8 = iota
	TypeFloat
	TypeString
	TypeFunction
	TypeBool
	TypeNil
	TypeClosure
)

type Int struct {
	Literal int64
}

func NewInt(literal int64) *Int {
	return &Int{
		Literal: literal,
	}
}

func (i *Int) String() string {
	return fmt.Sprintf("Int(%d)", i.Literal)
}

func (i *Int) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%d\n", i.Literal)
	return err
}

func (i *Int) ValueType() uint8 {
	return TypeInt
}

func (i *Int) WriteTo(w io.Writer) error {
	// 格式: [type:1byte][value:8bytes]
	if err := binary.Write(w, binary.BigEndian, i.ValueType()); err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, i.Literal)
}

func (i *Int) GetLiteral() any {
	return i.Literal
}

func (i *Int) SetLiteral(literal any) {
	i.Literal = literal.(int64)
}

type Float struct {
	Literal float64
}

func NewFloat(literal float64) *Float {
	return &Float{
		Literal: literal,
	}
}

func (f *Float) String() string {
	return fmt.Sprintf("Float(%f)", f.Literal)
}

func (f *Float) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%f\n", f.Literal)
	return err
}

func (f *Float) ValueType() uint8 {
	return TypeFloat
}

func (f *Float) WriteTo(w io.Writer) error {
	// 格式: [type:1byte][value:8bytes]
	if err := binary.Write(w, binary.BigEndian, f.ValueType()); err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, f.Literal)
}

func (f *Float) GetLiteral() any {
	return f.Literal
}

func (f *Float) SetLiteral(literal any) {
	f.Literal = literal.(float64)
}

type String struct {
	Literal string
}

func NewString(literal string) *String {
	return &String{
		Literal: literal,
	}
}

func (s *String) String() string {
	return fmt.Sprintf("String(%s)", s.Literal)
}

func (s *String) Print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s\n", s.Literal)
	return err
}

func (s *String) ValueType() uint8 {
	return TypeString
}

func (s *String) WriteTo(w io.Writer) error {
	// 格式: [type:1byte][length:8bytes][data:length bytes]
	if err := binary.Write(w, binary.BigEndian, s.ValueType()); err != nil {
		return err
	}
	length := int64(len(s.Literal))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := w.Write([]byte(s.Literal))
	return err
}

func (s *String) GetLiteral() any {
	return s.Literal
}

func (s *String) SetLiteral(literal any) {
	s.Literal = literal.(string)
}
