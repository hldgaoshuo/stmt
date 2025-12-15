package compiler

import (
	"reflect"
	"stmt/ast"
	object "stmt/object"
	"stmt/parser"
	"stmt/scanner"
	"testing"
)

func TestCompiler_CompileExpr(t *testing.T) {
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
		{
			name:      "true",
			source:    "true",
			code:      []uint8{OP_TRUE},
			constants: []*object.Object{},
		},
		{
			name:      "false",
			source:    "false",
			code:      []uint8{OP_FALSE},
			constants: []*object.Object{},
		},
		{
			name:      "nil",
			source:    "nil",
			code:      []uint8{OP_NIL},
			constants: []*object.Object{},
		},
		{
			name:   "!true",
			source: "!true",
			code: []uint8{
				OP_TRUE,
				OP_NOT,
			},
			constants: []*object.Object{},
		},
		{
			name:   "1<2",
			source: "1<2",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_CONSTANT, 1,
				OP_LT,
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
		{
			name:   `"abc"`,
			source: `"abc"`,
			code:   []uint8{OP_CONSTANT, 0},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralString{
						LiteralString: "abc",
					},
					ObjectType: object.ObjectType_OBJ_STRING,
				},
			},
		},
		{
			name:      "a",
			source:    "a",
			code:      nil,
			constants: nil,
			err:       ErrVariableNotDefined,
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
			if err != tt.err {
				t.Errorf("Compile() err = %v, want %v", code, tt.code)
			}
			if !reflect.DeepEqual(code, tt.code) {
				t.Errorf("Compile() code = %v, want %v", code, tt.code)
			}
			if !reflect.DeepEqual(constants, tt.constants) {
				t.Errorf("Compile() constants = %v, want %v", constants, tt.constants)
			}
		})
	}
}

func TestCompiler_CompileStmtDecl(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		err       error
		code      []uint8
		constants []*object.Object
	}{
		{
			name:   "expr",
			source: "1+2;",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_CONSTANT, 1,
				OP_ADD,
				OP_POP,
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
		{
			name:   "print",
			source: "print 1;",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_PRINT,
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
			name:   "var",
			source: "var a;",
			code: []uint8{
				OP_NIL,
				OP_SET_GLOBAL, 0,
			},
			constants: []*object.Object{},
		},
		{
			name:   "var 2",
			source: "var a = 1;",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
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
			name: "var 3",
			source: `
			var a = 1;
			a;
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_POP,
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
			name: "var 4",
			source: `
			var a = 1;
			print a;
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_PRINT,
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
			name: "var 5",
			source: `
			var a = 1;
			print a;
			{
				var a = 2;
				print a;
			}
			print a;
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_PRINT,
				OP_CONSTANT, 1,
				OP_SET_LOCAL, 0,
				OP_GET_LOCAL, 0,
				OP_PRINT,
				OP_GET_GLOBAL, 0,
				OP_PRINT,
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
		{
			name: "assign",
			source: `
			var a = 1;
			a = 2;
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_CONSTANT, 1,
				OP_SET_GLOBAL, 0,
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
		{
			name: "assign 2",
			source: `
			var a = 1;
			a = 2;
			print a;
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_CONSTANT, 1,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_PRINT,
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
		{
			name: "block",
			source: `
			var a = 1;
			{
				var a = 2;
			}
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_CONSTANT, 1,
				OP_SET_LOCAL, 0,
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
		{
			name: "if",
			source: `
			if (true)
			{
				print 10;
			}
			print 20;
			`,
			code: []uint8{
				OP_TRUE,
				OP_JUMP_FALSE, 9,
				OP_POP,
				OP_CONSTANT, 0,
				OP_PRINT,
				OP_JUMP, 10,
				OP_POP,
				OP_CONSTANT, 1,
				OP_PRINT,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 10,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 20,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
		},
		{
			name: "if else",
			source: `
			if (false)
			{
				print 10;
			}
			else
			{
				print 20;
			}
			`,
			code: []uint8{
				OP_FALSE,
				OP_JUMP_FALSE, 9,
				OP_POP,
				OP_CONSTANT, 0,
				OP_PRINT,
				OP_JUMP, 13,
				OP_POP,
				OP_CONSTANT, 1,
				OP_PRINT,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 10,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 20,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
		},
		{
			name: "and",
			source: `
			true and true;
			`,
			code: []uint8{
				OP_TRUE,
				OP_JUMP_FALSE, 5,
				OP_POP,
				OP_TRUE,
				OP_POP,
			},
			constants: []*object.Object{},
		},
		{
			name: "or",
			source: `
			true or true;
			`,
			code: []uint8{
				OP_TRUE,
				OP_JUMP_FALSE, 5,
				OP_JUMP, 7,
				OP_POP,
				OP_TRUE,
				OP_POP,
			},
			constants: []*object.Object{},
		},
		{
			name: "while",
			source: `
			var i = 0;
			while (i < 5)
			{
				print i;
				i = i + 1;
			}
			`,
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_CONSTANT, 1,
				OP_LT,
				OP_JUMP_FALSE, 24,
				OP_POP,
				OP_GET_GLOBAL, 0,
				OP_PRINT,
				OP_GET_GLOBAL, 0,
				OP_CONSTANT, 2,
				OP_ADD,
				OP_SET_GLOBAL, 0,
				OP_LOOP, 4,
				OP_POP,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 0,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 5,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
					ObjectType: object.ObjectType_OBJ_INT,
				},
			},
		},
		{
			name: "function",
			source: `
			fun pt() {
				print 1;
			}
			`,
			code: []uint8{
				OP_CLOSURE, 1,
				OP_SET_GLOBAL, 0,
			},
			constants: []*object.Object{
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_FUNCTION,
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 0,
								OP_PRINT,
								OP_NIL,
								OP_RETURN,
							},
							NumParams: 0,
						},
					},
				},
			},
		},
		{
			name: "function return nil",
			source: `
			fun pt() {
				print 1;
				return;
			}
			`,
			code: []uint8{
				OP_CLOSURE, 1,
				OP_SET_GLOBAL, 0,
			},
			constants: []*object.Object{
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_FUNCTION,
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 0,
								OP_PRINT,
								OP_NIL,
								OP_RETURN,
							},
							NumParams: 0,
						},
					},
				},
			},
		},
		{
			name: "function return value",
			source: `
			fun pt() {
				print 1;
				return 2;
			}
			`,
			code: []uint8{
				OP_CLOSURE, 2,
				OP_SET_GLOBAL, 0,
			},
			constants: []*object.Object{
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_FUNCTION,
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 0,
								OP_PRINT,
								OP_CONSTANT, 1,
								OP_RETURN,
							},
							NumParams: 0,
						},
					},
				},
			},
		},
		{
			name: "call",
			source: `
			fun pt() {
				print 1;
			}
			pt();
			`,
			code: []uint8{
				OP_CLOSURE, 1,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_CALL, 0,
				OP_POP,
			},
			constants: []*object.Object{
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_FUNCTION,
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 0,
								OP_PRINT,
								OP_NIL,
								OP_RETURN,
							},
							NumParams: 0,
						},
					},
				},
			},
		},
		{
			name: "call arg",
			source: `
			fun pt(a, b) {
				print a + b;
			}
			pt(1, 2);
			`,
			code: []uint8{
				OP_CLOSURE, 0,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_CONSTANT, 1,
				OP_CONSTANT, 2,
				OP_CALL, 2,
				OP_POP,
			},
			constants: []*object.Object{
				{
					ObjectType: object.ObjectType_OBJ_FUNCTION,
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_GET_LOCAL, 0,
								OP_GET_LOCAL, 1,
								OP_ADD,
								OP_PRINT,
								OP_NIL,
								OP_RETURN,
							},
							NumParams: 2,
						},
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
				},
			},
		},
		{
			name: "call arg return",
			source: `
			fun add(a, b) {
				return a + b;
			}
			print add(1, 2);
			`,
			code: []uint8{
				OP_CLOSURE, 0,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_CONSTANT, 1,
				OP_CONSTANT, 2,
				OP_CALL, 2,
				OP_PRINT,
			},
			constants: []*object.Object{
				{
					ObjectType: object.ObjectType_OBJ_FUNCTION,
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_GET_LOCAL, 0,
								OP_GET_LOCAL, 1,
								OP_ADD,
								OP_RETURN,
							},
							NumParams: 2,
						},
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					ObjectType: object.ObjectType_OBJ_INT,
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner_ := scanner.New(tt.source)
			tokens := scanner_.Scan()
			parser_ := parser.New(tokens)
			node, err := parser_.Parse()
			if err != nil {
				t.Errorf("Parse() err = %v", err)
				return
			}
			compiler_ := New(node)
			code, constants, err := compiler_.Compile()
			if err != nil {
				t.Errorf("Compile() err = %v", err)
				return
			}
			if !reflect.DeepEqual(code, tt.code) {
				t.Errorf("Compile() \n code: %v \n want: %v", code, tt.code)
			}
			if !reflect.DeepEqual(constants, tt.constants) {
				t.Errorf("Compile() constants = %v, want %v", constants, tt.constants)
			}
		})
	}
}
