package interpreter

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"stmt/ast"
	"stmt/token"

	"github.com/davecgh/go-spew/spew"
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

// Output 是一个可自定义的输出接口，默认为 os.Stdout
var Output io.Writer = os.Stdout

func Interpreter(decls []ast.Decl) (any, error) {
	env := newEnvironment()
	var result any
	var err error
	for _, decl := range decls {
		result, err = interpreter(decl, env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func interpreter(node ast.Node, env *environment) (any, error) {
	switch _node := node.(type) {
	case *ast.Literal:
		return _node.Value, nil
	case *ast.Variable:
		return env.get(_node.Name)
	case *ast.Grouping:
		return interpreter(_node.Expression, env)
	case *ast.Unary:
		right, err := interpreter(_node.Right, env)
		if err != nil {
			return nil, err
		}
		rightType := reflect.TypeOf(right)
		switch _node.Operator.TokenType {
		case token.BANG:
			if rightType.Kind() != reflect.Bool {
				slog.Error("operand must be a bool", "right type", rightType, "line", _node.Operator.Line)
				return nil, ErrOperandMustBeBool
			}
			return !right.(bool), nil
		case token.MINUS:
			if rightType.Kind() != reflect.Float64 {
				slog.Error("operand must be a float64", "right type", rightType, "line", _node.Operator.Line)
				return nil, ErrOperandMustBeFloat64
			}
			return -right.(float64), nil
		default:
			slog.Error("operator not support in unary", "operator", _node.Operator.TokenType, "line", _node.Operator.Line)
			return nil, ErrOperatorNotSupportInUnary
		}
	case *ast.Binary:
		left, err := interpreter(_node.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := interpreter(_node.Right, env)
		if err != nil {
			return nil, err
		}
		leftType := reflect.TypeOf(left)
		rightType := reflect.TypeOf(right)
		switch _node.Operator.TokenType {
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
			slog.Error("operator not support in binary", "expr", _node)
			return nil, ErrOperatorNotSupportInBinary
		}
	case *ast.Assign:
		value, err := interpreter(_node.Value, env)
		if err != nil {
			return nil, err
		}
		err = env.assign(_node.Name, value)
		return value, err
	case *ast.ExpressionStatement:
		_, err := interpreter(_node.Expression, env)
		return nil, err
	case *ast.Print:
		value, err := interpreter(_node.Expression, env)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(Output, "%#v\n", value)
		return nil, nil
	case *ast.Var:
		var value any = nil
		var err error
		if _node.Initializer != nil {
			value, err = interpreter(_node.Initializer, env)
			if err != nil {
				return nil, err
			}
		}
		err = env.define(_node.Name.Lexeme, value)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		slog.Error("node type not support")
		spew.Dump(_node)
		return nil, ErrExpressionTypeNotSupport
	}
}
