package interpreter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"stmt/ast"
	"stmt/token"
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
				return nil, ErrOperandMustBeBool
			}
			return !_right, nil
		case token.MINUS:
			_right, ok := right.(float64)
			if !ok {
				return nil, ErrOperandMustBeFloat64
			}
			return -_right, nil
		default:
			return nil, ErrOperatorNotSupportInUnary
		}
	case *ast.Call:
		callable, err := interpreter(_node.Callee, env)
		if err != nil {
			return nil, err
		}
		switch _callable := callable.(type) {
		case *ast.Class:
			ins := newInstance(_callable)
			return ins, nil
		case *closure:
			fun := _callable.Function
			lenParams := len(fun.Params)
			lenArgs := len(_node.Arguments)
			if lenParams != lenArgs {
				return nil, ErrNumParamsArgsNotMatch
			}
			_env := newEnvironment(_callable.Env)
			for i := 0; i < lenParams; i++ {
				param := fun.Params[i]
				arg := _node.Arguments[i]
				_arg, err := interpreter(arg, env)
				if err != nil {
					return nil, err
				}
				err = _env.define(param.Lexeme, _arg)
				if err != nil {
					return nil, err
				}
			}
			result, err := interpreter(fun.Body, _env)
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
			result, err := _callable(args...)
			if err != nil {
				return nil, err
			}
			return result, nil
		default:
			return nil, ErrFunctionNotDeclare
		}
	case *ast.Get:
		object, err := interpreter(_node.Object, env)
		if err != nil {
			return nil, err
		}
		ins, ok := object.(*instance)
		if !ok {
			return nil, ErrNotInstance
		}
		field, err := ins.get(_node.Name)
		if err != nil {
			return nil, err
		}
		return field, nil
	case *ast.Set:
		object, err := interpreter(_node.Object, env)
		if err != nil {
			return nil, err
		}
		ins, ok := object.(*instance)
		if !ok {
			print("Only instances have fields.")
			return nil, ErrOnlyInstanceHaveFields
		}
		value, err := interpreter(_node.Value, env)
		if err != nil {
			return nil, err
		}
		ins.set(_node.Name, value)
		return value, nil
	case *ast.Logical:
		left, err := interpreter(_node.Left, env)
		if err != nil {
			return nil, err
		}
		_left, ok := left.(bool)
		if !ok {
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
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.GREATER_EQUAL:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left >= _right, nil
			} else {
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.LESS:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left < _right, nil
			} else {
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.LESS_EQUAL:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left <= _right, nil
			} else {
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
				return nil, ErrOperandsMustBeTwoFloat64OrTwoString
			}
		case token.MINUS:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left - _right, nil
			} else {
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.STAR:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left * _right, nil
			} else {
				return nil, ErrOperandsMustBeTwoFloat64
			}
		case token.SLASH:
			_left, isLeftFloat64 := left.(float64)
			_right, isRightFloat64 := right.(float64)
			if isLeftFloat64 && isRightFloat64 {
				return _left / _right, nil
			} else {
				return nil, ErrOperandsMustBeTwoFloat64
			}
		default:
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
		_env := env.copy()
		clo := &closure{
			Function: _node,
			Env:      _env,
		}
		err := env.define(functionName, clo)
		if err != nil {
			return nil, err
		}
		err = _env.define(functionName, clo)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case *ast.Class:
		className := _node.Name.Lexeme
		err := env.define(className, _node)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, ErrExpressionTypeNotSupport
	}
}
