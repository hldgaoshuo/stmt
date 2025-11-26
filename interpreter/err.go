package interpreter

import "errors"

var (
	ErrReturn                   = errors.New("return")   // 用于 return 语句
	ErrBreak                    = errors.New("break")    // 用于 break 语句
	ErrContinue                 = errors.New("continue") // 用于 continue 语句
	ErrUndefinedVariable        = errors.New("undefined variable")
	ErrInvalidNodeType          = errors.New("invalid node type")
	ErrInvalidOperatorType      = errors.New("invalid operator type")
	ErrInvalidOperandType       = errors.New("invalid operand type")
	ErrInvalidOperandUnion      = errors.New("invalid operand union")
	ErrExpressionTypeNotSupport = errors.New("expression type not support")
	ErrFunctionNotDeclare       = errors.New("function not declare")
	ErrNumParamsArgsNotMatch    = errors.New("function parameters num should equ to call arguments num")
	ErrNotClass                 = errors.New("superclass must be a class")
	ErrNotInstance              = errors.New("only instances have properties")
	ErrOnlyInstanceHaveFields   = errors.New("only instances have fields")
	ErrUndefinedProperty        = errors.New("undefined property")
)
