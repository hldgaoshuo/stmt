package compiler

import "errors"

var (
	ErrOpcodeMismatch         = errors.New("opcode mismatch")
	ErrInvalidNodeType        = errors.New("invalid node type")
	ErrInvalidConstantIndex   = errors.New("invalid constant index")
	ErrInvalidClosureIndex    = errors.New("invalid closure index")
	ErrInvalidOperandType     = errors.New("invalid operand type")
	ErrInvalidOperatorType    = errors.New("invalid operator type")
	ErrInvalidSymbolScope     = errors.New("invalid symbol scope")
	ErrVariableNotDefined     = errors.New("variable not defined")
	ErrVariableAlreadyDefined = errors.New("variable already defined")
)
