package compiler

import (
	"encoding/binary"
	"math"
	"stmt/ast"
	"stmt/opcode"
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

func (s *Scope) ConstantEmit(index uint64) error {
	if index <= math.MaxUint8 {
		s.EmitWithOperand(opcode.OP_CONSTANT, index)
		return nil
	} else if index <= math.MaxUint16 {
		s.EmitWithOperand(opcode.OP_CONSTANT_2, index)
		return nil
	} else if index <= math.MaxUint32 {
		s.EmitWithOperand(opcode.OP_CONSTANT_4, index)
		return nil
	} else if index <= math.MaxUint64 {
		s.EmitWithOperand(opcode.OP_CONSTANT_8, index)
		return nil
	}
	return ErrInvalidConstantIndex
}

func (s *Scope) UnaryEmit(node *ast.Unary) error {
	switch node.Operator.TokenType {
	case token.MINUS:
		s.Emit(opcode.OP_NEGATE)
		return nil
	case token.BANG:
		s.Emit(opcode.OP_NOT)
		return nil
	default:
		return ErrInvalidOperatorType
	}
}

func (s *Scope) BinaryEmit(node *ast.Binary) error {
	switch node.Operator.TokenType {
	case token.PLUS:
		s.Emit(opcode.OP_ADD)
		return nil
	case token.MINUS:
		s.Emit(opcode.OP_SUBTRACT)
		return nil
	case token.STAR:
		s.Emit(opcode.OP_MULTIPLY)
		return nil
	case token.SLASH:
		s.Emit(opcode.OP_DIVIDE)
		return nil
	case token.PERCENTAGE:
		s.Emit(opcode.OP_MODULO)
		return nil
	case token.GREATER:
		s.Emit(opcode.OP_GT)
		return nil
	case token.LESS:
		s.Emit(opcode.OP_LT)
		return nil
	case token.EQUAL_EQUAL:
		s.Emit(opcode.OP_EQ)
		return nil
	case token.GREATER_EQUAL:
		s.Emit(opcode.OP_GE)
		return nil
	case token.LESS_EQUAL:
		s.Emit(opcode.OP_LE)
		return nil
	default:
		return ErrInvalidOperatorType
	}
}

func (s *Scope) SymbolGetEmit(symbolIndex uint64, symbolScope string) error {
	switch symbolScope {
	case LocalScope:
		s.EmitWithOperand(opcode.OP_GET_LOCAL, symbolIndex)
		return nil
	case UpScope:
		s.EmitWithOperand(opcode.OP_GET_UPVALUE, symbolIndex)
		return nil
	case GlobalScope:
		s.EmitWithOperand(opcode.OP_GET_GLOBAL, symbolIndex)
		return nil
	default:
		return ErrInvalidSymbolScope
	}
}

func (s *Scope) SymbolSetEmit(symbolIndex uint64, symbolScope string) error {
	switch symbolScope {
	case LocalScope:
		s.EmitWithOperand(opcode.OP_SET_LOCAL, symbolIndex)
		return nil
	case UpScope:
		s.EmitWithOperand(opcode.OP_SET_UPVALUE, symbolIndex)
		return nil
	case GlobalScope:
		s.EmitWithOperand(opcode.OP_SET_GLOBAL, symbolIndex)
		return nil
	default:
		return ErrInvalidSymbolScope
	}
}

func (s *Scope) ClosureEmit(index uint64, upValues []*UpInfo) error {
	if index <= math.MaxUint8 {
		s.EmitWithOperand(opcode.OP_CLOSURE, index)
	} else if index <= math.MaxUint16 {
		s.EmitWithOperand(opcode.OP_CLOSURE_2, index)
	} else if index <= math.MaxUint32 {
		s.EmitWithOperand(opcode.OP_CLOSURE_4, index)
	} else if index <= math.MaxUint64 {
		s.EmitWithOperand(opcode.OP_CLOSURE_8, index)
	} else {
		return ErrInvalidClosureIndex
	}
	for _, upInfo := range upValues {
		if upInfo.IsLocal {
			s.EmitOther(1)
		} else {
			s.EmitOther(0)
		}
		s.EmitOther(uint8(upInfo.LocalIndex))
	}
	return nil
}

func (s *Scope) Emit(opcode uint8) uint64 {
	offset := s.Offset()
	s.Code = append(s.Code, opcode)
	return offset
}

func (s *Scope) EmitWithOperand(opcode uint8, operand uint64) uint64 {
	offset := s.Offset()
	Code := CodeMake(opcode, operand)
	s.Code = append(s.Code, Code...)
	return offset
}

func (s *Scope) EmitOther(other uint8) {
	s.Code = append(s.Code, other)
}

func (s *Scope) Patch(offset uint64, op uint8) error {
	_op := s.Code[offset]
	if _op != op {
		return ErrOpcodeMismatch
	}
	_offset := s.Offset()
	length := _offset - offset - 5 // 所有 jump 指令长度为 5
	Code := CodeMake(op, length)
	copy(s.Code[offset:], Code)
	return nil
}

func (s *Scope) Loop(init uint64) {
	offset := s.Offset()
	length := offset - init + 1
	s.EmitWithOperand(opcode.OP_LOOP, length)
}

func (s *Scope) Offset() uint64 {
	offset := len(s.Code)
	return uint64(offset)
}

func CodeMake(op uint8, operand uint64) []uint8 {
	width := opcode.OperandWidth[op]
	instructions := make([]uint8, 1+width)
	instructions[0] = op
	offset := 1
	switch width {
	case 1:
		instructions[offset] = uint8(operand)
	case 2:
		binary.BigEndian.PutUint16(instructions[offset:], uint16(operand))
	case 4:
		binary.BigEndian.PutUint32(instructions[offset:], uint32(operand))
	case 8:
		binary.BigEndian.PutUint64(instructions[offset:], operand)
	}

	return instructions
}
