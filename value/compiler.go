package value

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Function struct {
	Code        []uint8
	NumParams   uint64
	NumUpvalues uint64
}

func NewFunction(code []uint8, numParams uint64, numUpvalues uint64) *Function {
	return &Function{
		Code:        code,
		NumParams:   numParams,
		NumUpvalues: numUpvalues,
	}
}

func (f *Function) String() string {
	return fmt.Sprintf("Function(%d, %d)%v", f.NumParams, f.NumUpvalues, f.Code)
}

func (f *Function) Print(w io.Writer) error {
	_, err := w.Write([]byte(f.String() + "\n"))
	return err
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
