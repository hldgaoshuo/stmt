package interpreter

import (
	"errors"
	"log/slog"
	"reflect"
	"stmt/ast"
	"stmt/token"
)

var (
	ErrExpressionTypeNotSupport            = errors.New("expression type not support")
	ErrOperatorNotSupportInUnary           = errors.New("operator not support in unary")
	ErrOperatorNotSupportInBinary          = errors.New("operator not support in binary")
	ErrOperandMustBeBool                   = errors.New("operand must be a bool")
	ErrOperandMustBeFloat64                = errors.New("operand must be a float64")
	ErrOperandsMustBeTwoFloat64            = errors.New("operand must be two float64")
	ErrOperandsMustBeTwoFloat64OrTwoString = errors.New("operand must be two float64 or two string")
)

func Interpreter(expression ast.Expr) (any, error) {
	switch expr := expression.(type) {
	case *ast.Literal:
		return expr.Value, nil
	case *ast.Grouping:
		return Interpreter(expr.Expression)
	case *ast.Unary:
		right, err := Interpreter(expr.Right)
		if err != nil {
			return nil, err
		}
		rightType := reflect.TypeOf(right)
		switch expr.Operator.TokenType {
		case token.BANG:
			if rightType.Kind() != reflect.Bool {
				slog.Error("operand must be a bool", "right type", rightType, "line", expr.Operator.Line)
				return nil, ErrOperandMustBeBool
			}
			return !right.(bool), nil
		case token.MINUS:
			if rightType.Kind() != reflect.Float64 {
				slog.Error("operand must be a float64", "right type", rightType, "line", expr.Operator.Line)
				return nil, ErrOperandMustBeFloat64
			}
			return -right.(float64), nil
		default:
			slog.Error("operator not support in unary", "operator", expr.Operator.TokenType, "line", expr.Operator.Line)
			return nil, ErrOperatorNotSupportInUnary
		}
	case *ast.Binary:
		left, err := Interpreter(expr.Left)
		if err != nil {
			return nil, err
		}
		right, err := Interpreter(expr.Right)
		if err != nil {
			return nil, err
		}
		leftType := reflect.TypeOf(left)
		rightType := reflect.TypeOf(right)
		switch expr.Operator.TokenType {
		case token.GREATER:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) > right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.GREATER_EQUAL:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) >= right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.LESS:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) < right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.LESS_EQUAL:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) <= right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.PLUS:
			if leftType.Kind() == reflect.String && rightType.Kind() == reflect.String {
				return left.(string) + right.(string), nil
			} else if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) + right.(float64), nil
			} else {
				slog.Error("operand must be two float64 or two string", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64OrTwoString
			}
		case token.MINUS:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) - right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.STAR:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) * right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.SLASH:
			if leftType.Kind() == reflect.Float64 && rightType.Kind() == reflect.Float64 {
				return left.(float64) / right.(float64), nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		default:
			slog.Error("operator not support in binary", "expr", expr)
			return nil, ErrOperatorNotSupportInBinary
		}
	default:
		slog.Error("expression type not support", "expression", expression)
		return nil, ErrExpressionTypeNotSupport
	}
}
