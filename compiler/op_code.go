package compiler

const (
	OP_RETURN uint8 = iota
	OP_CONSTANT
	OP_NEGATE
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_MODULO
	OP_TRUE
	OP_FALSE
	OP_NIL
	OP_NOT
	OP_EQ
	OP_GT
	OP_LT
	OP_GE
	OP_LE
	OP_POP
	OP_PRINT
	OP_SET_GLOBAL
	OP_GET_GLOBAL
	OP_SET_LOCAL
	OP_GET_LOCAL
)

// op 可能有多个操作数
var operandWidths = map[uint8][]int{
	OP_RETURN:     {},
	OP_CONSTANT:   {1},
	OP_NEGATE:     {},
	OP_ADD:        {},
	OP_SUBTRACT:   {},
	OP_MULTIPLY:   {},
	OP_DIVIDE:     {},
	OP_TRUE:       {},
	OP_FALSE:      {},
	OP_NIL:        {},
	OP_NOT:        {},
	OP_EQ:         {},
	OP_GT:         {},
	OP_LT:         {},
	OP_GE:         {},
	OP_LE:         {},
	OP_POP:        {},
	OP_PRINT:      {},
	OP_SET_GLOBAL: {1},
	OP_GET_GLOBAL: {1},
	OP_SET_LOCAL:  {1},
	OP_GET_LOCAL:  {1},
}
