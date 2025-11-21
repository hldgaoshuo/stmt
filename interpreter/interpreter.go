package interpreter

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"stmt/ast"
	"stmt/token"

	"github.com/davecgh/go-spew/spew"
)

// Output 是一个可自定义的输出接口，默认为 os.Stdout
var Output io.Writer = os.Stdout

func Interpreter(decls []ast.Node) error {
	var err error
	env := newEnvironment(nil)
	for funName, fun := range builtins {
		err = env.define(funName, fun)
		if err != nil {
			return err
		}
	}
	for _, decl := range decls {
		_, err = interpreter(decl, env)
		if err != nil {
			return err
		}
	}
	return nil
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
		switch _node.Operator.TokenType {
		case token.BANG:
			_right, ok := right.(bool)
			if !ok {
				slog.Error("operand must be a bool", "line", _node.Operator.Line)
				return nil, ErrOperandMustBeBool
			}
			return !_right, nil
		case token.MINUS:
			_right, ok := right.(float64)
			if !ok {
				slog.Error("operand must be a float64", "line", _node.Operator.Line)
				return nil, ErrOperandMustBeFloat64
			}
			return -_right, nil
		default:
			slog.Error("operator not support in unary", "operator", _node.Operator.TokenType, "line", _node.Operator.Line)
			return nil, ErrOperatorNotSupportInUnary
		}
	case *ast.Call:
		fun, err := interpreter(_node.Callee, env)
		if err != nil {
			return nil, err
		}
		switch _fun := fun.(type) {
		case *ast.Function:
			lenParams := len(_fun.Params)
			lenArgs := len(_node.Arguments)
			if lenParams != lenArgs {
				slog.Error("function parameters num should equ to call arguments num", "function parameters num", lenParams, "call arguments num", lenArgs)
				return nil, ErrNumParamsArgsNotMatch
			}
			_env := newEnvironment(env)
			for i := 0; i < lenParams; i++ {
				param := _fun.Params[i]
				arg := _node.Arguments[i]
				_arg, err := interpreter(arg, _env)
				if err != nil {
					return nil, err
				}
				err = _env.define(param.Lexeme, _arg)
				if err != nil {
					return nil, err
				}
			}
			result, err := interpreter(_fun.Body, _env)
			if err != nil {
				return nil, err
			}
			return result, nil
		case builtin:
			var args []any
			for _, arg := range _node.Arguments {
				_arg, err := interpreter(arg, env)
				if err != nil {
					return nil, err
				}
				args = append(args, _arg)
			}
			result, err := _fun(args...)
			if err != nil {
				return nil, err
			}
			return result, nil
		default:
			slog.Error("function not declare", "callee", _node.Callee)
			return nil, ErrFunctionNotDeclare
		}
	case *ast.Logical:
		left, err := interpreter(_node.Left, env)
		if err != nil {
			return nil, err
		}
		_left, ok := left.(bool)
		if !ok {
			slog.Error("operand must be a bool", "line", _node.Operator.Line)
			return nil, ErrOperandMustBeBool
		}
		switch _node.Operator.TokenType {
		case token.AND:
			if !_left {
				return false, nil
			} else {
				return interpreter(_node.Right, env)
			}
		case token.OR:
			if _left {
				return true, nil
			} else {
				return interpreter(_node.Right, env)
			}
		default:
			slog.Error("operator not support in logical", "operator", _node.Operator.TokenType, "line", _node.Operator.Line)
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
		switch _node.Operator.TokenType {
		case token.GREATER:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left > _right, nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.GREATER_EQUAL:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left >= _right, nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.LESS:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left < _right, nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.LESS_EQUAL:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left <= _right, nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.PLUS:
			_leftString, isLeftString := left.(string)
			_rightString, isRightString := right.(string)
			_leftFloat64, isLeftFloat64 := left.(float64)
			_rightFloat64, isRightFloat64 := right.(float64)
			if isLeftString && isRightString {
				return _leftString + _rightString, nil
			} else if isLeftFloat64 && isRightFloat64 {
				return _leftFloat64 + _rightFloat64, nil
			} else {
				slog.Error("operand must be two float64 or two string", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64OrTwoString
			}
		case token.MINUS:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left - _right, nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.STAR:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left * _right, nil
			} else {
				slog.Error("operand must be two float64", "left", left, "right", right)
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.SLASH:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left / _right, nil
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
		_, err = fmt.Fprintf(Output, "%#v\n", value)
		return nil, err
	case *ast.Block:
		_env := newEnvironment(env)
		for _, decl := range _node.Declarations {
			value, err := interpreter(decl, _env)
			if errors.Is(err, ErrReturn) {
				return value, nil
			}
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.If:
		condition, err := interpreter(_node.Condition, env)
		if err != nil {
			return nil, err
		}
		_condition, ok := condition.(bool)
		if !ok {
			slog.Error("condition result must be a bool", "line", _node.Line)
			return nil, ErrOperandMustBeBool
		}
		if _condition {
			return interpreter(_node.ThenBranch, env)
		} else {
			if _node.ElseBranch == nil {
				return nil, nil
			} else {
				return interpreter(_node.ElseBranch, env)
			}
		}
	case *ast.While:
		for {
			condition, err := interpreter(_node.Condition, env)
			if err != nil {
				return nil, err
			}
			_condition, ok := condition.(bool)
			if !ok {
				slog.Error("condition result must be a bool", "line", _node.Line)
				return nil, ErrOperandMustBeBool
			}
			if !_condition {
				break
			}
			_, err = interpreter(_node.Body, env)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.Return:
		value, err := interpreter(_node.Expression, env)
		if err != nil {
			return nil, err
		}
		return value, ErrReturn
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
	case *ast.Function:
		functionName := _node.Name.Lexeme
		err := env.define(functionName, _node)
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
