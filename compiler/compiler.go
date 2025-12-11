package compiler

import (
	"stmt/ast"
	object "stmt/object"
	"stmt/token"
)

type Compiler struct {
	ast       []ast.Node
	constants []*object.Object
}

func New(ast []ast.Node) *Compiler {
	return &Compiler{
		ast:       ast,
		constants: []*object.Object{},
	}
}

func (c *Compiler) Compile() ([]uint8, []*object.Object, error) {
	symbolTable := NewSymbolTable(nil)
	for _, node := range c.ast {
		err := c.collectGlobal(node, symbolTable)
		if err != nil {
			return nil, nil, err
		}
	}
	mainScope := NewScope()
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
		return symbolTable.SetGlobal(_node.Name.Lexeme)
	case *ast.Function:
		return symbolTable.SetGlobal(_node.Name.Lexeme)
	default:
		return nil
	}
}

func (c *Compiler) compile(node ast.Node, symbolTable *SymbolTable, scope *Scope) error {
	switch _node := node.(type) {
	case *ast.Literal:
		switch value := _node.Value.(type) {
		case int64:
			obj := &object.Object{
				Literal: &object.Object_LiteralInt{
					LiteralInt: value,
				},
				ObjectType: object.ObjectType_OBJ_INT,
			}
			index := c.constantAdd(obj)
			scope.CodeEmit(OP_CONSTANT, index)
			return nil
		case float64:
			obj := &object.Object{
				Literal: &object.Object_LiteralFloat{
					LiteralFloat: value,
				},
				ObjectType: object.ObjectType_OBJ_FLOAT,
			}
			index := c.constantAdd(obj)
			scope.CodeEmit(OP_CONSTANT, index)
			return nil
		case bool:
			if value {
				scope.CodeEmit(OP_TRUE)
			} else {
				scope.CodeEmit(OP_FALSE)
			}
			return nil
		case nil:
			scope.CodeEmit(OP_NIL)
			return nil
		case string:
			obj := &object.Object{
				Literal: &object.Object_LiteralString{
					LiteralString: value,
				},
				ObjectType: object.ObjectType_OBJ_STRING,
			}
			index := c.constantAdd(obj)
			scope.CodeEmit(OP_CONSTANT, index)
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
		switch _node.Operator.TokenType {
		case token.MINUS:
			scope.CodeEmit(OP_NEGATE)
			return nil
		case token.BANG:
			scope.CodeEmit(OP_NOT)
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.Binary:
		err := c.compile(_node.Left, symbolTable, scope)
		if err != nil {
			return err
		}
		err = c.compile(_node.Right, symbolTable, scope)
		if err != nil {
			return err
		}
		switch _node.Operator.TokenType {
		case token.PLUS:
			scope.CodeEmit(OP_ADD)
			return nil
		case token.MINUS:
			scope.CodeEmit(OP_SUBTRACT)
			return nil
		case token.STAR:
			scope.CodeEmit(OP_MULTIPLY)
			return nil
		case token.SLASH:
			scope.CodeEmit(OP_DIVIDE)
			return nil
		case token.PERCENTAGE:
			scope.CodeEmit(OP_MODULO)
			return nil
		case token.GREATER:
			scope.CodeEmit(OP_GT)
			return nil
		case token.LESS:
			scope.CodeEmit(OP_LT)
			return nil
		case token.EQUAL_EQUAL:
			scope.CodeEmit(OP_EQ)
			return nil
		case token.GREATER_EQUAL:
			scope.CodeEmit(OP_GE)
			return nil
		case token.LESS_EQUAL:
			scope.CodeEmit(OP_LE)
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.ExpressionStatement:
		err := c.compile(_node.Expression, symbolTable, scope)
		if err != nil {
			return err
		}
		if _, ok := _node.Expression.(*ast.Assign); !ok {
			scope.CodeEmit(OP_POP)
		}
		return nil
	case *ast.Print:
		err := c.compile(_node.Expression, symbolTable, scope)
		if err != nil {
			return err
		}
		scope.CodeEmit(OP_PRINT)
		return nil
	case *ast.Var:
		if _node.Initializer == nil {
			scope.CodeEmit(OP_NIL)
		} else {
			err := c.compile(_node.Initializer, symbolTable, scope)
			if err != nil {
				return err
			}
		}
		symbolInfo, err := symbolTable.Set(_node.Name.Lexeme)
		if err != nil {
			return err
		}
		switch symbolInfo.Scope {
		case LocalScope:
			scope.CodeEmit(OP_SET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			scope.CodeEmit(OP_SET_GLOBAL, int(symbolInfo.Index))
			return nil
		default:
			return ErrInvalidSymbolScope
		}
	case *ast.Variable:
		symbolInfo, ok := symbolTable.Get(_node.Name.Lexeme)
		if !ok {
			return ErrVariableNotDefined
		}
		switch symbolInfo.Scope {
		case LocalScope:
			scope.CodeEmit(OP_GET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			scope.CodeEmit(OP_GET_GLOBAL, int(symbolInfo.Index))
			return nil
		default:
			return ErrInvalidSymbolScope
		}
	case *ast.Assign:
		err := c.compile(_node.Value, symbolTable, scope)
		if err != nil {
			return err
		}
		symbolInfo, ok := symbolTable.Get(_node.Name.Lexeme)
		if !ok {
			return ErrVariableNotDefined
		}
		switch symbolInfo.Scope {
		case LocalScope:
			scope.CodeEmit(OP_SET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			scope.CodeEmit(OP_SET_GLOBAL, int(symbolInfo.Index))
			return nil
		default:
			return ErrInvalidSymbolScope
		}
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
		offsetFalse := scope.CodeEmit(OP_JUMP_FALSE, 0)
		scope.CodeEmit(OP_POP)
		err = c.compile(_node.ThenBranch, symbolTable, scope)
		if err != nil {
			return err
		}
		offset := scope.CodeEmit(OP_JUMP, 0)
		err = scope.CodePatch(offsetFalse, OP_JUMP_FALSE)
		if err != nil {
			return err
		}
		scope.CodeEmit(OP_POP)
		if _node.ElseBranch != nil {
			err = c.compile(_node.ElseBranch, symbolTable, scope)
			if err != nil {
				return err
			}
		}
		err = scope.CodePatch(offset, OP_JUMP)
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
			offsetFalse := scope.CodeEmit(OP_JUMP_FALSE, 0)
			scope.CodeEmit(OP_POP)
			err = c.compile(_node.Right, symbolTable, scope)
			if err != nil {
				return err
			}
			err = scope.CodePatch(offsetFalse, OP_JUMP_FALSE)
			if err != nil {
				return err
			}
			return nil
		case token.OR:
			offsetFalse := scope.CodeEmit(OP_JUMP_FALSE, 0)
			offset := scope.CodeEmit(OP_JUMP, 0)
			err = scope.CodePatch(offsetFalse, OP_JUMP_FALSE)
			if err != nil {
				return err
			}
			scope.CodeEmit(OP_POP)
			err = c.compile(_node.Right, symbolTable, scope)
			if err != nil {
				return err
			}
			err = scope.CodePatch(offset, OP_JUMP)
			if err != nil {
				return err
			}
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.While:
		loop := scope.CodeOffset()
		err := c.compile(_node.Condition, symbolTable, scope)
		if err != nil {
			return err
		}
		offsetFalse := scope.CodeEmit(OP_JUMP_FALSE, 0)
		scope.CodeEmit(OP_POP)
		err = c.compile(_node.Body, symbolTable, scope)
		if err != nil {
			return err
		}
		scope.CodeEmit(OP_LOOP, loop)
		err = scope.CodePatch(offsetFalse, OP_JUMP_FALSE)
		if err != nil {
			return err
		}
		scope.CodeEmit(OP_POP)
		return nil
	case *ast.Function:
		symbolInfo, err := symbolTable.Set(_node.Name.Lexeme)
		if err != nil {
			return err
		}
		_symbolTable := NewSymbolTable(symbolTable)
		for _, param := range _node.Params {
			_, err = _symbolTable.Set(param.Lexeme)
			if err != nil {
				return err
			}
		}
		_scope := NewScope()
		err = c.compile(_node.Body, _symbolTable, _scope)
		if err != nil {
			return err
		}
		obj := &object.Object{
			ObjectType: object.ObjectType_OBJ_FUNCTION,
			Literal: &object.Object_LiteralFunction{
				LiteralFunction: &object.Function{
					Code:      _scope.Code,
					NumParams: uint64(len(_node.Params)),
				},
			},
		}
		index := c.constantAdd(obj)
		scope.CodeEmit(OP_CONSTANT, index)
		switch symbolInfo.Scope {
		case LocalScope:
			scope.CodeEmit(OP_SET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			scope.CodeEmit(OP_SET_GLOBAL, int(symbolInfo.Index))
			return nil
		default:
			return ErrInvalidSymbolScope
		}
	default:
		return ErrInvalidNodeType
	}
}

// constants
func (c *Compiler) constantAdd(obj *object.Object) int {
	c.constants = append(c.constants, obj)
	index := len(c.constants) - 1
	return index
}
