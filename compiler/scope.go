package compiler

import (
	"encoding/binary"
	"stmt/ast"
	"stmt/token"
)

// todo 需要一个新的字段，标识当前作用域是否为 main，如果 main 中包含 return，需要在编译阶段报错
type Scope struct {
	Code       []uint8
	HaveReturn bool
}

func NewScope(haveReturn bool) *Scope {
	return &Scope{
		Code:       []uint8{},
		HaveReturn: haveReturn,
	}
}

func (s *Scope) UnaryEmit(node *ast.Unary) error {
	switch node.Operator.TokenType {
	case token.MINUS:
		s.CodeEmit(OP_NEGATE)
		return nil
	case token.BANG:
		s.CodeEmit(OP_NOT)
		return nil
	default:
		return ErrInvalidOperatorType
	}
}

func (s *Scope) BinaryEmit(node *ast.Binary) error {
	switch node.Operator.TokenType {
	case token.PLUS:
		s.CodeEmit(OP_ADD)
		return nil
	case token.MINUS:
		s.CodeEmit(OP_SUBTRACT)
		return nil
	case token.STAR:
		s.CodeEmit(OP_MULTIPLY)
		return nil
	case token.SLASH:
		s.CodeEmit(OP_DIVIDE)
		return nil
	case token.PERCENTAGE:
		s.CodeEmit(OP_MODULO)
		return nil
	case token.GREATER:
		s.CodeEmit(OP_GT)
		return nil
	case token.LESS:
		s.CodeEmit(OP_LT)
		return nil
	case token.EQUAL_EQUAL:
		s.CodeEmit(OP_EQ)
		return nil
	case token.GREATER_EQUAL:
		s.CodeEmit(OP_GE)
		return nil
	case token.LESS_EQUAL:
		s.CodeEmit(OP_LE)
		return nil
	default:
		return ErrInvalidOperatorType
	}
}

func (s *Scope) SymbolGetEmit(symbolInfo *SymbolInfo, symbolScope string) error {
	switch symbolScope {
	case LocalScope:
		s.CodeEmit(OP_GET_LOCAL, int(symbolInfo.Index))
		return nil
	case CloserScope:
		s.CodeEmit(OP_GET_UPVALUE, int(symbolInfo.Index))
		return nil
	case GlobalScope:
		s.CodeEmit(OP_GET_GLOBAL, int(symbolInfo.Index))
		return nil
	default:
		return ErrInvalidSymbolScope
	}
}

func (s *Scope) SymbolSetEmit(symbolInfo *SymbolInfo, symbolScope string) error {
	switch symbolScope {
	case LocalScope:
		s.CodeEmit(OP_SET_LOCAL, int(symbolInfo.Index))
		return nil
	case CloserScope:
		s.CodeEmit(OP_SET_UPVALUE, int(symbolInfo.Index))
		return nil
	case GlobalScope:
		s.CodeEmit(OP_SET_GLOBAL, int(symbolInfo.Index))
		return nil
	default:
		return ErrInvalidSymbolScope
	}
}

func (s *Scope) CodePatch(offset int, op uint8) error {
	_op := s.Code[offset]
	if _op != op {
		return ErrOpCodeMismatch
	}
	operand := s.CodeOffset()
	Code := s.CodeMake(op, operand)
	copy(s.Code[offset:], Code)
	return nil
}

func (s *Scope) CodeEmit(op uint8, operands ...int) int {
	offset := s.CodeOffset()
	Code := s.CodeMake(op, operands...)
	s.Code = append(s.Code, Code...)
	return offset
}

func (s *Scope) CodeOffset() int {
	offset := len(s.Code)
	return offset
}

func (s *Scope) CodeMake(op byte, operands ...int) []byte {
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
