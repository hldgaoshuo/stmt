package compiler

import (
	"reflect"
	"stmt/ast"
	object "stmt/object"
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
		constants []*object.Object
	}{
		{
			name:   "1",
			source: "1",
			code:   []uint8{OP_CONSTANT, 0},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
		},
		{
			name:   "1.2",
			source: "1.2",
			code:   []uint8{OP_CONSTANT, 0},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralFloat{
						LiteralFloat: 1.2,
					},
					ObjectType: object.ObjectType_OBJ_FLOAT,
				},
			},
		},
		{
			name:   "(1)",
			source: "(1)",
			code:   []uint8{OP_CONSTANT, 0},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
		},
		{
			name:   "-1",
			source: "-1",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_NEGATE,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
		},
		{
			name:   "1+2",
			source: "1+2",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_CONSTANT, 1,
				OP_ADD,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
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
			if err != nil {
				t.Errorf("Compile() err = %v", err)
				return
			}
			if !reflect.DeepEqual(code, tt.code) {
				t.Errorf("Compile() got = %v, want %v", code, tt.code)
				return
			}
			if !reflect.DeepEqual(constants, tt.constants) {
				t.Errorf("Compile() got1 = %v, want %v", constants, tt.constants)
				return
			}
			err = compiler_.chunk(tt.name)
			if err != nil {
				t.Errorf("chunk() err = %v", err)
				return
			}
		})
	}
}
