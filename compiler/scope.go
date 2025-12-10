package compiler

import "encoding/binary"

type Scope struct {
	Code []uint8
}

func NewScope() *Scope {
	return &Scope{
		Code: []uint8{},
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
