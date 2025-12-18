package value

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Value interface {
	String() string
	WriteTo(w io.Writer) error
	ValueType() uint8
}

const (
	TypeInt uint8 = iota
	TypeFloat
	TypeString
	TypeFunction
)

type Int struct {
	Literal int64
}

func (i *Int) String() string {
	return fmt.Sprintf("Int(%d)", i.Literal)
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

type Float struct {
	Literal float64
}

func (f *Float) String() string {
	return fmt.Sprintf("Float(%f)", f.Literal)
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

type String struct {
	Literal string
}

func (s *String) String() string {
	return fmt.Sprintf("String(%s)", s.Literal)
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

type Function struct {
	Code        []byte
	NumParams   uint64
	NumUpvalues uint64
}

func (f *Function) String() string {
	return fmt.Sprintf("Function(%d, %d)", f.NumParams, f.NumUpvalues)
}

func (f *Function) ValueType() uint8 {
	return TypeFunction
}

func (f *Function) WriteTo(w io.Writer) error {
	// 格式: [type:1byte][numParams:8bytes][numUpvalues:8bytes][codeLength:8bytes][code:codeLength bytes]
	if err := binary.Write(w, binary.BigEndian, f.ValueType()); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, f.NumParams); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, f.NumUpvalues); err != nil {
		return err
	}
	codeLength := int64(len(f.Code))
	if err := binary.Write(w, binary.BigEndian, codeLength); err != nil {
		return err
	}
	_, err := w.Write(f.Code)
	return err
}
