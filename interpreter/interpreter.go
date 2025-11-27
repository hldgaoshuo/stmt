package interpreter

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
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
	case *ast.Super:
		// 特殊的 *ast.Get
		super, err := env.get(_node.Keyword)
		if err != nil {
			return nil, err
		}
		superClass, ok := super.(*class)
		if !ok {
			return nil, ErrNotClass
		}
		clo := superClass.get(_node.Method)
		if clo == nil {
			return nil, ErrUndefinedProperty
		}
		ins, err := env.this()
		if err != nil {
			return nil, err
		}
		_ins, ok := ins.(*instance)
		if !ok {
			return nil, ErrNotInstance
		}
		_clo, err := clo.bind(_ins)
		if err != nil {
			return nil, err
		}
		return _clo, nil
	case *ast.This:
		return env.get(_node.Keyword)
	case *ast.Call:
		callable, err := interpreter(_node.Callee, env)
		if err != nil {
			return nil, err
		}
		switch _callable := callable.(type) {
		case *class:
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
		property, err := ins.get(_node.Name)
		if err != nil {
			return nil, err
		}
		return property, nil
	case *ast.Unary:
		right, err := interpreter(_node.Right, env)
		if err != nil {
			return nil, err
		}
		switch rightValue := right.(type) {
		case int64:
			switch _node.Operator.TokenType {
			case token.MINUS:
				return -rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		case float64:
			switch _node.Operator.TokenType {
			case token.MINUS:
				return -rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		case bool:
			switch _node.Operator.TokenType {
			case token.BANG:
				return !rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		default:
			return nil, ErrInvalidOperandType
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
		leftType := reflect.TypeOf(left).Kind()
		rightType := reflect.TypeOf(right).Kind()
		if leftType == reflect.Int64 && rightType == reflect.Int64 {
			leftValue := left.(int64)
			rightValue := right.(int64)
			switch _node.Operator.TokenType {
			case token.PLUS:
				return leftValue + rightValue, nil
			case token.MINUS:
				return leftValue - rightValue, nil
			case token.STAR:
				return leftValue * rightValue, nil
			case token.SLASH:
				return leftValue / rightValue, nil
			case token.PERCENTAGE:
				return leftValue % rightValue, nil
			case token.EQUAL_EQUAL:
				return leftValue == rightValue, nil
			case token.BANG_EQUAL:
				return leftValue != rightValue, nil
			case token.GREATER:
				return leftValue > rightValue, nil
			case token.GREATER_EQUAL:
				return leftValue >= rightValue, nil
			case token.LESS:
				return leftValue < rightValue, nil
			case token.LESS_EQUAL:
				return leftValue <= rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		} else if leftType == reflect.Float64 && rightType == reflect.Float64 {
			leftValue := left.(float64)
			rightValue := right.(float64)
			switch _node.Operator.TokenType {
			case token.PLUS:
				return leftValue + rightValue, nil
			case token.MINUS:
				return leftValue - rightValue, nil
			case token.STAR:
				return leftValue * rightValue, nil
			case token.SLASH:
				return leftValue / rightValue, nil
			case token.PERCENTAGE:
				return math.Mod(leftValue, rightValue), nil
			case token.EQUAL_EQUAL:
				return leftValue == rightValue, nil
			case token.BANG_EQUAL:
				return leftValue != rightValue, nil
			case token.GREATER:
				return leftValue > rightValue, nil
			case token.GREATER_EQUAL:
				return leftValue >= rightValue, nil
			case token.LESS:
				return leftValue < rightValue, nil
			case token.LESS_EQUAL:
				return leftValue <= rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		} else if leftType == reflect.Int64 && rightType == reflect.Float64 {
			leftValue := float64(left.(int64))
			rightValue := right.(float64)
			switch _node.Operator.TokenType {
			case token.PLUS:
				return leftValue + rightValue, nil
			case token.MINUS:
				return leftValue - rightValue, nil
			case token.STAR:
				return leftValue * rightValue, nil
			case token.SLASH:
				return leftValue / rightValue, nil
			case token.PERCENTAGE:
				return math.Mod(leftValue, rightValue), nil
			case token.EQUAL_EQUAL:
				return leftValue == rightValue, nil
			case token.BANG_EQUAL:
				return leftValue != rightValue, nil
			case token.GREATER:
				return leftValue > rightValue, nil
			case token.GREATER_EQUAL:
				return leftValue >= rightValue, nil
			case token.LESS:
				return leftValue < rightValue, nil
			case token.LESS_EQUAL:
				return leftValue <= rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		} else if leftType == reflect.Float64 && rightType == reflect.Int64 {
			leftValue := left.(float64)
			rightValue := float64(right.(int64))
			switch _node.Operator.TokenType {
			case token.PLUS:
				return leftValue + rightValue, nil
			case token.MINUS:
				return leftValue - rightValue, nil
			case token.STAR:
				return leftValue * rightValue, nil
			case token.SLASH:
				return leftValue / rightValue, nil
			case token.PERCENTAGE:
				return math.Mod(leftValue, rightValue), nil
			case token.EQUAL_EQUAL:
				return leftValue == rightValue, nil
			case token.BANG_EQUAL:
				return leftValue != rightValue, nil
			case token.GREATER:
				return leftValue > rightValue, nil
			case token.GREATER_EQUAL:
				return leftValue >= rightValue, nil
			case token.LESS:
				return leftValue < rightValue, nil
			case token.LESS_EQUAL:
				return leftValue <= rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		} else if leftType == reflect.Bool && rightType == reflect.Bool {
			leftValue := left.(bool)
			rightValue := right.(bool)
			switch _node.Operator.TokenType {
			case token.EQUAL_EQUAL:
				return leftValue == rightValue, nil
			case token.BANG_EQUAL:
				return leftValue != rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		} else if leftType == reflect.String && rightType == reflect.String {
			leftValue := left.(string)
			rightValue := right.(string)
			switch _node.Operator.TokenType {
			case token.PLUS:
				return leftValue + rightValue, nil
			case token.EQUAL_EQUAL:
				return leftValue == rightValue, nil
			case token.BANG_EQUAL:
				return leftValue != rightValue, nil
			case token.GREATER:
				return leftValue > rightValue, nil
			case token.GREATER_EQUAL:
				return leftValue >= rightValue, nil
			case token.LESS:
				return leftValue < rightValue, nil
			case token.LESS_EQUAL:
				return leftValue <= rightValue, nil
			default:
				return nil, ErrInvalidOperatorType
			}
		} else {
			return nil, ErrInvalidOperandUnion
		}
	case *ast.Logical:
		left, err := interpreter(_node.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := interpreter(_node.Right, env)
		if err != nil {
			return nil, err
		}
		leftValue, ok := left.(bool)
		if !ok {
			return nil, ErrInvalidOperandType
		}
		rightValue, ok := right.(bool)
		if !ok {
			return nil, ErrInvalidOperandType
		}
		switch _node.Operator.TokenType {
		case token.AND:
			return leftValue && rightValue, nil
		case token.OR:
			return leftValue || rightValue, nil
		default:
			return nil, ErrInvalidOperatorType
		}
	case *ast.Assign:
		value, err := interpreter(_node.Value, env)
		if err != nil {
			return nil, err
		}
		err = env.assign(_node.Name, value)
		return value, err
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
	case *ast.ExpressionStatement:
		_, err := interpreter(_node.Expression, env)
		return nil, err
	case *ast.Break:
		return nil, ErrBreak
	case *ast.Continue:
		return nil, ErrContinue
	case *ast.Return:
		value, err := interpreter(_node.Expression, env)
		if err != nil {
			return nil, err
		}
		return value, ErrReturn
	case *ast.While:
		for {
			condition, err := interpreter(_node.Condition, env)
			if err != nil {
				return nil, err
			}
			_condition, ok := condition.(bool)
			if !ok {
				return nil, ErrInvalidOperandType
			}
			if !_condition {
				break
			}
			_, err = interpreter(_node.Body, env)
			if errors.Is(err, ErrBreak) {
				return nil, nil
			}
			if errors.Is(err, ErrContinue) {
				continue
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
			return nil, ErrInvalidOperandType
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
	case *ast.Print:
		value, err := interpreter(_node.Expression, env)
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprintf(Output, "%#v\n", value)
		return nil, err
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
		var superClass *class
		if _node.SuperClass != nil {
			super, err := interpreter(_node.SuperClass, env)
			if err != nil {
				return nil, err
			}
			var ok bool
			superClass, ok = super.(*class)
			if !ok {
				return nil, ErrNotClass
			}
		}
		cls := &class{
			SuperClass: superClass,
			Closures:   []*closure{},
		}
		_env := env.copy()
		if superClass != nil {
			_env = newEnvironment(_env)
			err := _env.define("super", superClass)
			if err != nil {
				return nil, err
			}
		}
		for _, method := range _node.Methods {
			clo := &closure{
				Function: method,
				Env:      _env,
			}
			cls.Closures = append(cls.Closures, clo)
		}
		err := env.define(className, cls)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, ErrExpressionTypeNotSupport
	}
}
