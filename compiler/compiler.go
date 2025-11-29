package compiler

import "stmt/ast"

type Compiler struct {
	// in
	ast []ast.Node
	// out
	code      []uint8
	constants []int64
}
