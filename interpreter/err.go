package interpreter

import "errors"

var (
	ErrReturn                              = errors.New("return")
	ErrUndefinedVariable                   = errors.New("undefined variable")
	ErrExpressionTypeNotSupport            = errors.New("expression type not support")
	ErrOperatorNotSupportInUnary           = errors.New("operator not support in unary")
	ErrOperatorNotSupportInBinary          = errors.New("operator not support in binary")
	ErrOperandMustBeBool                   = errors.New("operand must be a bool")
	ErrOperandMustBeFloat64                = errors.New("operand must be a float64")
	ErrOperandsMustBeTwoFloat64            = errors.New("operand must be two float64")
	ErrOperandsMustBeTwoFloat64OrTwoString = errors.New("operand must be two float64 or two string")
	ErrFunctionNotDeclare                  = errors.New("function not declare")
	ErrNumParamsArgsNotMatch               = errors.New("function parameters num should equ to call arguments num")
	ErrNotInstance                         = errors.New("only instances have properties")
	ErrOnlyInstanceHaveFields              = errors.New("only instances have fields")
	ErrUndefinedProperty                   = errors.New("undefined property")
)
