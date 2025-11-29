package compiler

import (
	"encoding/binary"
	"errors"
	"stmt/ast"
	"stmt/token"
)

var (
	ErrInvalidNodeType     = errors.New("invalid node type")
	ErrInvalidOperandType  = errors.New("invalid operand type")
	ErrInvalidOperatorType = errors.New("invalid operator type")
)

type Compiler struct {
	// in
	ast []ast.Node
	// out
	code      []uint8
	constants []int64
}

func New(ast []ast.Node) *Compiler {
	return &Compiler{
		ast:       ast,
		code:      []uint8{},
		constants: []int64{},
	}
}

func (c *Compiler) Compile() ([]uint8, []int64, error) {
	for _, node := range c.ast {
		err := c.compile(node)
		if err != nil {
			return nil, nil, err
		}
	}
	return c.code, c.constants, nil
}

func (c *Compiler) compile(node ast.Node) error {
	switch _node := node.(type) {
	case *ast.Literal:
		switch value := _node.Value.(type) {
		case int64:
			index := c.constantAdd(value)
			c.codeEmit(OP_CONSTANT, index)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *ast.Grouping:
		err := c.compile(_node.Expression)
		if err != nil {
			return err
		}
		return nil
	case *ast.Unary:
		err := c.compile(_node.Right)
		if err != nil {
			return err
		}
		switch _node.Operator.TokenType {
		case token.MINUS:
			c.codeEmit(OP_NEGATE)
			return nil
		default:
			return ErrInvalidOperatorType
		}
	case *ast.Binary:
		err := c.compile(_node.Left)
		if err != nil {
			return err
		}
		err = c.compile(_node.Right)
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
		default:
			return ErrInvalidOperatorType
		}
	default:
		return ErrInvalidNodeType
	}
}

// code
func (c *Compiler) codeEmit(op uint8, operands ...int) {
	code := c.codeMake(op, operands...)
	c.code = append(c.code, code...)
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
func (c *Compiler) constantAdd(obj int64) int {
	c.constants = append(c.constants, obj)
	index := len(c.constants) - 1
	return index
}
