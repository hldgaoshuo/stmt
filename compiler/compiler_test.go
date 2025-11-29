package compiler

import (
	"reflect"
	"stmt/ast"
	"stmt/parser"
	"stmt/scanner"
	"testing"
)

func TestCompiler_Compile(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		err       error
		code      []uint8
		constants []int64
	}{
		{
			name:      "1",
			source:    "1",
			code:      []uint8{OP_CONSTANT, 0},
			constants: []int64{1},
		},
		{
			name:      "(1)",
			source:    "(1)",
			code:      []uint8{OP_CONSTANT, 0},
			constants: []int64{1},
		},
		{
			name:   "-1",
			source: "-1",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_NEGATE,
			},
			constants: []int64{1},
		},
		{
			name:   "1+2",
			source: "1+2",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_CONSTANT, 1,
				OP_ADD,
			},
			constants: []int64{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner_ := scanner.New(tt.source)
			tokens := scanner_.Scan()
			parser_ := parser.New(tokens)
			node, err := parser_.Expression()
			if err != nil {
				t.Errorf("Parse() err = %v", err)
				return
			}
			compiler_ := New([]ast.Node{node})
			code, constants, err := compiler_.Compile()
			if !reflect.DeepEqual(code, tt.code) {
				t.Errorf("Compile() got = %v, want %v", code, tt.code)
			}
			if !reflect.DeepEqual(constants, tt.constants) {
				t.Errorf("Compile() got1 = %v, want %v", constants, tt.constants)
			}
		})
	}
}
