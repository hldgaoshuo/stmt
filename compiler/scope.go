package compiler

import "encoding/binary"

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
