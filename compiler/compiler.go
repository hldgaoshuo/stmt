package compiler

import (
	"encoding/binary"
	"os"
	"stmt/ast"
	object "stmt/object"
	"stmt/token"

	"google.golang.org/protobuf/proto"
)

type Compiler struct {
	// in
	ast []ast.Node
	// out
	code      []uint8
	constants []*object.Object
}

func New(ast []ast.Node) *Compiler {
	return &Compiler{
		ast:       ast,
		code:      []uint8{},
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
	for _, node := range c.ast {
		err := c.compile(node, symbolTable)
		if err != nil {
			return nil, nil, err
		}
	}
	return c.code, c.constants, nil
}

func (c *Compiler) collectGlobal(node ast.Node, symbolTable *SymbolTable) error {
	switch _node := node.(type) {
	case *ast.Var:
		return symbolTable.SetGlobal(_node.Name.Lexeme)
	default:
		return nil
	}
}

func (c *Compiler) compile(node ast.Node, symbolTable *SymbolTable) error {
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
			c.codeEmit(OP_CONSTANT, index)
			return nil
		case float64:
			obj := &object.Object{
				Literal: &object.Object_LiteralFloat{
					LiteralFloat: value,
				},
				ObjectType: object.ObjectType_OBJ_FLOAT,
			}
			index := c.constantAdd(obj)
			c.codeEmit(OP_CONSTANT, index)
			return nil
		case bool:
			if value {
				c.codeEmit(OP_TRUE)
			} else {
				c.codeEmit(OP_FALSE)
			}
			return nil
		case nil:
			c.codeEmit(OP_NIL)
			return nil
		case string:
			obj := &object.Object{
				Literal: &object.Object_LiteralString{
					LiteralString: value,
				},
				ObjectType: object.ObjectType_OBJ_STRING,
			}
			index := c.constantAdd(obj)
			c.codeEmit(OP_CONSTANT, index)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *ast.Grouping:
		err := c.compile(_node.Expression, symbolTable)
		if err != nil {
			return err
		}
		return nil
	case *ast.Unary:
		err := c.compile(_node.Right, symbolTable)
		if err != nil {
			return err
		}
		switch _node.Operator.TokenType {
		case token.MINUS:
			c.codeEmit(OP_NEGATE)
			return nil
		case token.BANG:
			c.codeEmit(OP_NOT)
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.Binary:
		err := c.compile(_node.Left, symbolTable)
		if err != nil {
			return err
		}
		err = c.compile(_node.Right, symbolTable)
		if err != nil {
			return err
		}
		switch _node.Operator.TokenType {
		case token.PLUS:
			c.codeEmit(OP_ADD)
			return nil
		case token.MINUS:
			c.codeEmit(OP_SUBTRACT)
			return nil
		case token.STAR:
			c.codeEmit(OP_MULTIPLY)
			return nil
		case token.SLASH:
			c.codeEmit(OP_DIVIDE)
			return nil
		case token.PERCENTAGE:
			c.codeEmit(OP_MODULO)
			return nil
		case token.GREATER:
			c.codeEmit(OP_GT)
			return nil
		case token.LESS:
			c.codeEmit(OP_LT)
			return nil
		case token.EQUAL_EQUAL:
			c.codeEmit(OP_EQ)
			return nil
		case token.GREATER_EQUAL:
			c.codeEmit(OP_GE)
			return nil
		case token.LESS_EQUAL:
			c.codeEmit(OP_LE)
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.ExpressionStatement:
		err := c.compile(_node.Expression, symbolTable)
		if err != nil {
			return err
		}
		if _, ok := _node.Expression.(*ast.Assign); !ok {
			c.codeEmit(OP_POP)
		}
		return nil
	case *ast.Print:
		err := c.compile(_node.Expression, symbolTable)
		if err != nil {
			return err
		}
		c.codeEmit(OP_PRINT)
		return nil
	case *ast.Var:
		if _node.Initializer == nil {
			c.codeEmit(OP_NIL)
		} else {
			err := c.compile(_node.Initializer, symbolTable)
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
			c.codeEmit(OP_SET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			c.codeEmit(OP_SET_GLOBAL, int(symbolInfo.Index))
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
			c.codeEmit(OP_GET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			c.codeEmit(OP_GET_GLOBAL, int(symbolInfo.Index))
			return nil
		default:
			return ErrInvalidSymbolScope
		}
	case *ast.Assign:
		err := c.compile(_node.Value, symbolTable)
		if err != nil {
			return err
		}
		symbolInfo, ok := symbolTable.Get(_node.Name.Lexeme)
		if !ok {
			return ErrVariableNotDefined
		}
		switch symbolInfo.Scope {
		case LocalScope:
			c.codeEmit(OP_SET_LOCAL, int(symbolInfo.Index))
			return nil
		case GlobalScope:
			c.codeEmit(OP_SET_GLOBAL, int(symbolInfo.Index))
			return nil
		default:
			return ErrInvalidSymbolScope
		}
	case *ast.Block:
		_symbolTable := NewSymbolTable(symbolTable)
		for _, statement := range _node.Declarations {
			err := c.compile(statement, _symbolTable)
			if err != nil {
				return err
			}
		}
		return nil
	case *ast.If:
		err := c.compile(_node.Condition, symbolTable)
		if err != nil {
			return err
		}
		offsetFalse := c.codeEmit(OP_JUMP_FALSE, 0)
		c.codeEmit(OP_POP)
		err = c.compile(_node.ThenBranch, symbolTable)
		if err != nil {
			return err
		}
		offset := c.codeEmit(OP_JUMP, 0)
		err = c.codePatch(offsetFalse, OP_JUMP_FALSE)
		if err != nil {
			return err
		}
		c.codeEmit(OP_POP)
		if _node.ElseBranch != nil {
			err = c.compile(_node.ElseBranch, symbolTable)
			if err != nil {
				return err
			}
		}
		err = c.codePatch(offset, OP_JUMP)
		if err != nil {
			return err
		}
		return nil
	case *ast.Logical:
		err := c.compile(_node.Left, symbolTable)
		if err != nil {
			return err
		}
		switch _node.Operator.TokenType {
		case token.AND:
			offsetFalse := c.codeEmit(OP_JUMP_FALSE, 0)
			c.codeEmit(OP_POP)
			err = c.compile(_node.Right, symbolTable)
			if err != nil {
				return err
			}
			err = c.codePatch(offsetFalse, OP_JUMP_FALSE)
			if err != nil {
				return err
			}
			return nil
		case token.OR:
			offsetFalse := c.codeEmit(OP_JUMP_FALSE, 0)
			offset := c.codeEmit(OP_JUMP, 0)
			err = c.codePatch(offsetFalse, OP_JUMP_FALSE)
			if err != nil {
				return err
			}
			c.codeEmit(OP_POP)
			err = c.compile(_node.Right, symbolTable)
			if err != nil {
				return err
			}
			err = c.codePatch(offset, OP_JUMP)
			if err != nil {
				return err
			}
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.While:
		loop := c.codeOffset()
		err := c.compile(_node.Condition, symbolTable)
		if err != nil {
			return err
		}
		offsetFalse := c.codeEmit(OP_JUMP_FALSE, 0)
		c.codeEmit(OP_POP)
		err = c.compile(_node.Body, symbolTable)
		if err != nil {
			return err
		}
		c.codeEmit(OP_LOOP, loop)
		err = c.codePatch(offsetFalse, OP_JUMP_FALSE)
		if err != nil {
			return err
		}
		c.codeEmit(OP_POP)
		return nil
	default:
		return ErrInvalidNodeType
	}
}

// code
func (c *Compiler) codePatch(offset int, op uint8) error {
	_op := c.code[offset]
	if _op != op {
		return ErrOpCodeMismatch
	}
	operand := c.codeOffset()
	code := c.codeMake(op, operand)
	copy(c.code[offset:], code)
	return nil
}

func (c *Compiler) codeEmit(op uint8, operands ...int) int {
	offset := c.codeOffset()
	code := c.codeMake(op, operands...)
	c.code = append(c.code, code...)
	return offset
}

func (c *Compiler) codeOffset() int {
	offset := len(c.code)
	return offset
}

func (c *Compiler) codeMake(op byte, operands ...int) []byte {
	widths, ok := operandWidths[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range widths {
		instructionLen += w
	}

	instructions := make([]byte, instructionLen)
	instructions[0] = op
	offset := 1
	for i, o := range operands {
		width := widths[i]
		switch width {
		case 1:
			instructions[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instructions[offset:], uint16(o))
		}
		offset += width
	}

	return instructions
}

// constants
func (c *Compiler) constantAdd(obj *object.Object) int {
	c.constants = append(c.constants, obj)
	index := len(c.constants) - 1
	return index
}

// chunk
func (c *Compiler) chunk(name string) error {
	chunk := &object.Chunk{
		Code:         c.code,
		Constants:    c.constants,
		GlobalsCount: Global.NumDefinitions,
	}

	row, err := proto.Marshal(chunk)
	if err != nil {
		return err
	}

	path := "../vm2" + name + ".bin"

	err = os.WriteFile(path, row, 0644)
	if err != nil {
		return err
	}

	return nil
}
