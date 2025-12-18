package compiler

const (
	OP_CONSTANT byte = iota
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
	OP_JUMP_FALSE
	OP_JUMP
	OP_AND
	OP_OR
	OP_LOOP
	OP_CALL
	OP_RETURN
	OP_CLOSURE
	OP_GET_UPVALUE
	OP_SET_UPVALUE
)

// op 可能有多个操作数
var operandWidths = map[byte][]int{
	OP_CONSTANT:    {2},
	OP_NEGATE:      {},
	OP_ADD:         {},
	OP_SUBTRACT:    {},
	OP_MULTIPLY:    {},
	OP_DIVIDE:      {},
	OP_TRUE:        {},
	OP_FALSE:       {},
	OP_NIL:         {},
	OP_NOT:         {},
	OP_EQ:          {},
	OP_GT:          {},
	OP_LT:          {},
	OP_GE:          {},
	OP_LE:          {},
	OP_POP:         {},
	OP_PRINT:       {},
	OP_SET_GLOBAL:  {2},
	OP_GET_GLOBAL:  {2},
	OP_SET_LOCAL:   {2},
	OP_GET_LOCAL:   {2},
	OP_JUMP_FALSE:  {4},
	OP_JUMP:        {4},
	OP_AND:         {},
	OP_OR:          {},
	OP_LOOP:        {4},
	OP_CALL:        {2},
	OP_RETURN:      {},
	OP_CLOSURE:     {2},
	OP_GET_UPVALUE: {2},
	OP_SET_UPVALUE: {2},
}
