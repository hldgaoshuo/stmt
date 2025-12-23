package vm

import "errors"

var (
	ErrInvalidOpcodeType   = errors.New("invalid opcode type")
	ErrOpcodeHaveNoOperand = errors.New("opcode have no operand")
	ErrInvalidOperandWidth = errors.New("invalid operand width")
	ErrInvalidOperandType  = errors.New("invalid operand type")
	ErrInvalidCondType     = errors.New("invalid cond type")
	ErrInvalidCallType     = errors.New("invalid call type")
	ErrInvalidClosureType  = errors.New("invalid closure type")
	ErrZeroInDivide        = errors.New("zero in divide")
	ErrZeroInModulo        = errors.New("zero in modulo")
)
