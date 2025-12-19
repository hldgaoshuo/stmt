package compiler

import (
	"stmt/ast"
	"stmt/opcode"
	"stmt/token"
	"stmt/value"
)

type Compiler struct {
	ast       []ast.Node
	constants []value.Value
}

func New(ast []ast.Node) *Compiler {
	return &Compiler{
		ast:       ast,
		constants: []value.Value{},
	}
}

func (c *Compiler) Compile() ([]uint8, []value.Value, error) {
	symbolTable := NewSymbolTable(nil)
	for _, node := range c.ast {
		err := c.collectGlobal(node, symbolTable)
		if err != nil {
			return nil, nil, err
		}
	}
	mainScope := NewScope(false)
	for _, node := range c.ast {
		err := c.compile(node, symbolTable, mainScope)
		if err != nil {
			return nil, nil, err
		}
	}
	return mainScope.Code, c.constants, nil
}

func (c *Compiler) collectGlobal(node ast.Node, symbolTable *SymbolTable) error {
	switch _node := node.(type) {
	case *ast.Var:
		return symbolTable.DefineGlobal(_node.Name.Lexeme)
	case *ast.Function:
		return symbolTable.DefineGlobal(_node.Name.Lexeme)
	default:
		return nil
	}
}

func (c *Compiler) compile(node ast.Node, symbolTable *SymbolTable, scope *Scope) error {
	switch _node := node.(type) {
	case *ast.Literal:
		switch value_ := _node.Value.(type) {
		case int64:
			obj := value.NewInt(value_)
			index := c.constantAdd(obj)
			err := scope.ConstantEmit(index)
			if err != nil {
				return err
			}
			return nil
		case float64:
			obj := value.NewFloat(value_)
			index := c.constantAdd(obj)
			err := scope.ConstantEmit(index)
			if err != nil {
				return err
			}
			return nil
		case bool:
			if value_ {
				scope.Emit(opcode.OP_TRUE)
			} else {
				scope.Emit(opcode.OP_FALSE)
			}
			return nil
		case nil:
			scope.Emit(opcode.OP_NIL)
			return nil
		case string:
			obj := value.NewString(value_)
			index := c.constantAdd(obj)
			err := scope.ConstantEmit(index)
			if err != nil {
				return err
			}
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *ast.Grouping:
		err := c.compile(_node.Expression, symbolTable, scope)
		if err != nil {
			return err
		}
		return nil
	case *ast.Unary:
		err := c.compile(_node.Right, symbolTable, scope)
		if err != nil {
			return err
		}
		err = scope.UnaryEmit(_node)
		if err != nil {
			return err
		}
		return nil
	case *ast.Binary:
		err := c.compile(_node.Left, symbolTable, scope)
		if err != nil {
			return err
		}
		err = c.compile(_node.Right, symbolTable, scope)
		if err != nil {
			return err
		}
		err = scope.BinaryEmit(_node)
		if err != nil {
			return err
		}
		return nil
	case *ast.ExpressionStatement:
		err := c.compile(_node.Expression, symbolTable, scope)
		if err != nil {
			return err
		}
		if _, ok := _node.Expression.(*ast.Assign); !ok {
			scope.Emit(opcode.OP_POP)
		}
		return nil
	case *ast.Print:
		err := c.compile(_node.Expression, symbolTable, scope)
		if err != nil {
			return err
		}
		scope.Emit(opcode.OP_PRINT)
		return nil
	case *ast.Var:
		if _node.Initializer == nil {
			scope.Emit(opcode.OP_NIL)
		} else {
			err := c.compile(_node.Initializer, symbolTable, scope)
			if err != nil {
				return err
			}
		}
		symbolIndex, symbolScope, err := symbolTable.Define(_node.Name.Lexeme)
		if err != nil {
			return err
		}
		err = scope.SymbolSetEmit(symbolIndex, symbolScope)
		if err != nil {
			return err
		}
		return nil
	case *ast.Variable:
		symbolIndex, symbolScope, ex := symbolTable.Get(_node.Name.Lexeme)
		if !ex {
			return ErrVariableNotDefined
		}
		err := scope.SymbolGetEmit(symbolIndex, symbolScope)
		if err != nil {
			return err
		}
		return nil
	case *ast.Assign:
		err := c.compile(_node.Value, symbolTable, scope)
		if err != nil {
			return err
		}
		symbolIndex, symbolScope, ex := symbolTable.Get(_node.Name.Lexeme)
		if !ex {
			return ErrVariableNotDefined
		}
		err = scope.SymbolSetEmit(symbolIndex, symbolScope)
		if err != nil {
			return err
		}
		return nil
	case *ast.Block:
		_symbolTable := NewSymbolTable(symbolTable)
		for _, statement := range _node.Declarations {
			err := c.compile(statement, _symbolTable, scope)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.If:
		err := c.compile(_node.Condition, symbolTable, scope)
		if err != nil {
			return err
		}
		offsetFalse := scope.EmitWithOperand(opcode.OP_JUMP_FALSE, 0)
		scope.Emit(opcode.OP_POP)
		err = c.compile(_node.ThenBranch, symbolTable, scope)
		if err != nil {
			return err
		}
		offset := scope.EmitWithOperand(opcode.OP_JUMP, 0)
		err = scope.Patch(offsetFalse, opcode.OP_JUMP_FALSE)
		if err != nil {
			return err
		}
		scope.Emit(opcode.OP_POP)
		if _node.ElseBranch != nil {
			err = c.compile(_node.ElseBranch, symbolTable, scope)
			if err != nil {
				return err
			}
		}
		err = scope.Patch(offset, opcode.OP_JUMP)
		if err != nil {
			return err
		}
		return nil
	case *ast.Logical:
		err := c.compile(_node.Left, symbolTable, scope)
		if err != nil {
			return err
		}
		switch _node.Operator.TokenType {
		case token.AND:
			offsetFalse := scope.EmitWithOperand(opcode.OP_JUMP_FALSE, 0)
			scope.Emit(opcode.OP_POP)
			err = c.compile(_node.Right, symbolTable, scope)
			if err != nil {
				return err
			}
			err = scope.Patch(offsetFalse, opcode.OP_JUMP_FALSE)
			if err != nil {
				return err
			}
			return nil
		case token.OR:
			offsetFalse := scope.EmitWithOperand(opcode.OP_JUMP_FALSE, 0)
			offset := scope.EmitWithOperand(opcode.OP_JUMP, 0)
			err = scope.Patch(offsetFalse, opcode.OP_JUMP_FALSE)
			if err != nil {
				return err
			}
			scope.Emit(opcode.OP_POP)
			err = c.compile(_node.Right, symbolTable, scope)
			if err != nil {
				return err
			}
			err = scope.Patch(offset, opcode.OP_JUMP)
			if err != nil {
				return err
			}
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.While:
		loop := scope.Offset()
		err := c.compile(_node.Condition, symbolTable, scope)
		if err != nil {
			return err
		}
		offsetFalse := scope.EmitWithOperand(opcode.OP_JUMP_FALSE, 0)
		scope.Emit(opcode.OP_POP)
		for _, statement := range _node.Body.Declarations {
			err = c.compile(statement, symbolTable, scope)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
		scope.EmitWithOperand(opcode.OP_LOOP, loop)
		err = scope.Patch(offsetFalse, opcode.OP_JUMP_FALSE)
		if err != nil {
			return err
		}
		scope.Emit(opcode.OP_POP)
		return nil
	case *ast.Function:
		symbolIndex, symbolScope, err := symbolTable.Define(_node.Name.Lexeme)
		if err != nil {
			return err
		}
		_symbolTable := NewSymbolTable(symbolTable)
		for _, param := range _node.Params {
			_, _, err = _symbolTable.Define(param.Lexeme)
			if err != nil {
				return err
			}
		}
		_scope := NewScope(false)
		for _, statement := range _node.Body.Declarations {
			err = c.compile(statement, _symbolTable, _scope)
			if err != nil {
				return err
			}
		}
		if !_scope.HaveReturn {
			_scope.Emit(opcode.OP_NIL)
			_scope.Emit(opcode.OP_RETURN)
		}
		obj := value.NewFunction(_scope.Code, uint64(len(_node.Params)), uint64(len(_symbolTable.UpValues)))
		index := c.constantAdd(obj)
		err = scope.ClosureEmit(index, _symbolTable.UpValues)
		if err != nil {
			return err
		}
		err = scope.SymbolSetEmit(symbolIndex, symbolScope)
		if err != nil {
			return err
		}
		return nil
	case *ast.Call:
		err := c.compile(_node.Callee, symbolTable, scope)
		if err != nil {
			return err
		}
		for _, argument := range _node.Arguments {
			err = c.compile(argument, symbolTable, scope)
			if err != nil {
				return err
			}
		}
		scope.EmitWithOperand(opcode.OP_CALL, uint64(len(_node.Arguments)))
		return nil
	case *ast.Return:
		scope.HaveReturn = true
		if _node.Expression != nil {
			err := c.compile(_node.Expression, symbolTable, scope)
			if err != nil {
				return err
			}
		} else {
			scope.Emit(opcode.OP_NIL)
		}
		scope.Emit(opcode.OP_RETURN)
		return nil
	default:
		return ErrInvalidNodeType
	}
}

// constants
func (c *Compiler) constantAdd(obj value.Value) uint64 {
	c.constants = append(c.constants, obj)
	index := len(c.constants) - 1
	return uint64(index)
}
