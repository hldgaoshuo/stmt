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
)

// op 可能有多个操作数
var operandWidths = map[uint8][]int{
	OP_RETURN:   {},
	OP_CONSTANT: {1},
	OP_NEGATE:   {},
	OP_ADD:      {},
	OP_SUBTRACT: {},
	OP_MULTIPLY: {},
	OP_DIVIDE:   {},
	OP_TRUE:     {},
	OP_FALSE:    {},
	OP_NIL:      {},
	OP_NOT:      {},
	OP_EQ:       {},
	OP_GT:       {},
	OP_LT:       {},
	OP_GE:       {},
	OP_LE:       {},
}
