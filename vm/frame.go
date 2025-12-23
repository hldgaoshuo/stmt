package vm

import (
	"encoding/binary"
	"stmt/opcode"
	"stmt/value"
)

type Frame struct {
	Closure     *value.Closure
	BasePointer uint64
	Ip          uint64
}

func NewFrame(closure *value.Closure, basePointer uint64) *Frame {
	return &Frame{
		Closure:     closure,
		BasePointer: basePointer,
		Ip:          0,
	}
}

func (f *Frame) MoveIp(offset uint64) {
	f.Ip += offset
}

func (f *Frame) CodeSize() uint64 {
	function := f.Closure.Function
	code := function.Code
	size := len(code)
	return uint64(size)
}

func (f *Frame) Opcode() uint8 {
	function := f.Closure.Function
	code := function.Code
	char := code[f.Ip]
	f.Ip++
	return char
}

func (f *Frame) Operand(op uint8) (uint64, error) {
	width, ok := opcode.OperandWidth[op]
	if !ok {
		return 0, ErrOpcodeHaveNoOperand
	}
	function := f.Closure.Function
	code := function.Code
	switch width {
	case 1:
		operand := uint64(code[f.Ip])
		f.Ip++
		return operand, nil
	case 2:
		operand := uint64(binary.BigEndian.Uint16(code[f.Ip:]))
		f.Ip += 2
		return operand, nil
	case 4:
		operand := uint64(binary.BigEndian.Uint32(code[f.Ip:]))
		f.Ip += 4
		return operand, nil
	case 8:
		operand := binary.BigEndian.Uint64(code[f.Ip:])
		f.Ip += 8
		return operand, nil
	default:
		return 0, ErrInvalidOperandWidth
	}
}

func (f *Frame) CodeNext() uint8 {
	code := f.Closure.Function.Code
	char := code[f.Ip]
	f.Ip++
	return char
}
