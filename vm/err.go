package vm

import "errors"

var (
	ErrInvalidOpcodeType   = errors.New("invalid opcode type")
	ErrOpcodeHaveNoOperand = errors.New("opcode have no operand")
	ErrInvalidOperandWidth = errors.New("invalid operand width")
	ErrInvalidOperandType  = errors.New("invalid operand type")
)
