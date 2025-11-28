package compiler

const (
	OP_RETURN uint8 = iota
	OP_CONSTANT
	OP_NEGATE
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
)

// op 可能有多个操作数
var operandWidths = map[uint8][]int{
	OP_RETURN:   []int{},
	OP_CONSTANT: []int{1},
	OP_NEGATE:   []int{},
	OP_ADD:      []int{},
	OP_SUBTRACT: []int{},
	OP_MULTIPLY: []int{},
	OP_DIVIDE:   []int{},
}
