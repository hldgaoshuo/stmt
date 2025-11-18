package interpreter

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
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
		leftType := reflect.TypeOf(left)
		if leftType.Kind() != reflect.Bool {
			slog.Error("operand must be a bool", "left type", leftType, "line", _node.Operator.Line)
			return nil, ErrOperandMustBeBool
		}
		switch _node.Operator.TokenType {
		case token.AND:
			if !left.(bool) {
				return false, nil
			} else {
				return interpreter(_node.Right, env)
			}
		case token.OR:
			if left.(bool) {
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
	case *ast.Block:
		_env := newEnvironment(env)
		for _, decl := range _node.Declarations {
			value, err := interpreter(decl, _env)
			if err == ErrReturn {
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
		conditionType := reflect.TypeOf(condition)
		if conditionType.Kind() != reflect.Bool {
			slog.Error("condition result must be a bool", "condition type", conditionType, "line", _node.Line)
			return nil, ErrOperandMustBeBool
		}
		if condition.(bool) {
			return interpreter(_node.ThenBranch, env)
		} else {
			if _node.ElseBranch == nil {
				return nil, nil
			} else {
				return interpreter(_node.ElseBranch, env)
			}
		}
	case *ast.While:
		condition, err := interpreter(_node.Condition, env)
		if err != nil {
			return nil, err
		}
		conditionType := reflect.TypeOf(condition)
		if conditionType.Kind() != reflect.Bool {
			slog.Error("condition result must be a bool", "condition type", conditionType, "line", _node.Line)
			return nil, ErrOperandMustBeBool
		}
		for condition.(bool) {
			_, err = interpreter(_node.Body, env)
			if err != nil {
				return nil, err
			}
			condition, err = interpreter(_node.Condition, env)
			if err != nil {
				return nil, err
			}
			conditionType = reflect.TypeOf(condition)
			if conditionType.Kind() != reflect.Bool {
				slog.Error("condition result must be a bool", "condition type", conditionType, "line", _node.Line)
				return nil, ErrOperandMustBeBool
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
