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
		switch value := _node.Value.(type) {
		case int64:
			obj := &object.Object{
				Literal: &object.Object_LiteralInt{
					LiteralInt: value,
				},
			}
			index := c.constantAdd(obj)
			scope.CodeEmit(OP_CONSTANT, index)
			return nil
		case float64:
			obj := &object.Object{
				Literal: &object.Object_LiteralFloat{
					LiteralFloat: value,
				},
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
		for _, statement := range _node.Body.Declarations {
			err = c.compile(statement, symbolTable, scope)
			if err != nil {
				return err
			}
		}
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
			_scope.CodeEmit(OP_NIL)
			_scope.CodeEmit(OP_RETURN)
		}
		obj := &object.Object{
			Literal: &object.Object_LiteralFunction{
				LiteralFunction: &object.Function{
					Code:        _scope.Code,
					NumParams:   uint64(len(_node.Params)),
					NumUpvalues: uint64(len(_symbolTable.UpValues)),
				},
			},
		}
		index := c.constantAdd(obj)
		scope.CodeEmit(OP_CLOSURE, index)
		for _, upInfo := range _symbolTable.UpValues {
			if upInfo.IsLocal {
				scope.CodeEmitClosureMeta(1)
			} else {
				scope.CodeEmitClosureMeta(0)
			}
			scope.CodeEmitClosureMeta(uint8(upInfo.LocalIndex))
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
		scope.CodeEmit(OP_CALL, len(_node.Arguments))
		return nil
	case *ast.Return:
		scope.HaveReturn = true
		if _node.Expression != nil {
			err := c.compile(_node.Expression, symbolTable, scope)
			if err != nil {
				return err
			}
		} else {
			scope.CodeEmit(OP_NIL)
		}
		scope.CodeEmit(OP_RETURN)
		return nil
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
