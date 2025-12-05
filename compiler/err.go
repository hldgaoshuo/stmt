package compiler

import "errors"

var (
	ErrInvalidNodeType        = errors.New("invalid node type")
	ErrInvalidOperandType     = errors.New("invalid operand type")
	ErrInvalidOperatorType    = errors.New("invalid operator type")
	ErrInvalidSymbolScope     = errors.New("invalid symbol scope")
	ErrVariableNotDefined     = errors.New("variable not defined")
	ErrVariableAlreadyDefined = errors.New("variable already defined")
)
