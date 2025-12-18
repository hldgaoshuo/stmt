package compiler2

import (
	"fmt"
	"reflect"
	"stmt/ast"
	object "stmt/object"
	"stmt/parser"
	"stmt/scanner"
	"strings"
	"testing"
)

func formatObjects(objects []*object.Object) string {
	if len(objects) == 0 {
		return "[]"
	}

	result := make([]string, len(objects))
	for i, obj := range objects {
		result[i] = formatObject(obj)
	}
	return "[" + strings.Join(result, ", ") + "]"
}

func formatObject(obj *object.Object) string {
	if obj == nil {
		return "nil"
	}

	switch literal := obj.Literal.(type) {
	case *object.Object_LiteralInt:
		return fmt.Sprintf("int(%d)", literal.LiteralInt)
	case *object.Object_LiteralFloat:
		return fmt.Sprintf("float(%g)", literal.LiteralFloat)
	case *object.Object_LiteralString:
		return fmt.Sprintf("string(%q)", literal.LiteralString)
	case *object.Object_LiteralFunction:
		fn := literal.LiteralFunction
		return fmt.Sprintf("function(params=%d, upvalues=%d, code=%v)", fn.NumParams, fn.NumUpvalues, fn.Code)
	case *object.Object_LiteralBool:
		if literal.LiteralBool {
			return "bool(true)"
		}
		return "bool(false)"
	case *object.Object_LiteralNil:
		return "nil"
	default:
		return "unknown"
	}
}

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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 20,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 20,
					},
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
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 5,
					},
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
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
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
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
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
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
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
				},
				{
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
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
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
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
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
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
				},
			},
		},
		{
			name: "closure",
			source: `
			fun outer() {
				var x = "outside";
				fun inner() {
					print x;
				}
				inner();
			}
			outer();
			`,
			code: []uint8{
				OP_CLOSURE, 2,
				OP_SET_GLOBAL, 0,
				OP_GET_GLOBAL, 0,
				OP_CALL, 0,
				OP_POP,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralString{
						LiteralString: "outside",
					},
				},
				{
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_GET_UPVALUE, 0,
								OP_PRINT,
								OP_NIL,
								OP_RETURN,
							},
							NumParams:   0,
							NumUpvalues: 1,
						},
					},
				},
				{
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 0,
								OP_SET_LOCAL, 0,
								OP_CLOSURE, 1,
								1, 0,
								OP_SET_LOCAL, 1,
								OP_GET_LOCAL, 1,
								OP_CALL, 0,
								OP_POP,
								OP_NIL,
								OP_RETURN,
							},
							NumParams:   0,
							NumUpvalues: 0,
						},
					},
				},
			},
		},
		{
			name: "closure 1",
			source: `
			fun outer() {
				var a = 1;
				var b = 2;
				fun middle() {
					var c = 3;
					var d = 4;
					fun inner() {
						print a + c + b + d;
					}
				}
			}
			`,
			code: []uint8{
				OP_CLOSURE, 6,
				OP_SET_GLOBAL, 0,
			},
			constants: []*object.Object{
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 1,
					},
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 2,
					},
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 3,
					},
				},
				{
					Literal: &object.Object_LiteralInt{
						LiteralInt: 4,
					},
				},
				{
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_GET_UPVALUE, 0,
								OP_GET_UPVALUE, 1,
								OP_ADD,
								OP_GET_UPVALUE, 2,
								OP_ADD,
								OP_GET_UPVALUE, 3,
								OP_ADD,
								OP_PRINT,
								OP_NIL,
								OP_RETURN,
							},
							NumParams:   0,
							NumUpvalues: 4,
						},
					},
				},
				{
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 2,
								OP_SET_LOCAL, 0,
								OP_CONSTANT, 3,
								OP_SET_LOCAL, 1,
								OP_CLOSURE, 4,
								0, 0, 1, 0, 0, 1, 1, 1,
								OP_SET_LOCAL, 2,
								OP_NIL,
								OP_RETURN,
							},
							NumParams:   0,
							NumUpvalues: 2,
						},
					},
				},
				{
					Literal: &object.Object_LiteralFunction{
						LiteralFunction: &object.Function{
							Code: []uint8{
								OP_CONSTANT, 0,
								OP_SET_LOCAL, 0,
								OP_CONSTANT, 1,
								OP_SET_LOCAL, 1,
								OP_CLOSURE, 5,
								1, 0, 1, 1,
								OP_SET_LOCAL, 2,
								OP_NIL,
								OP_RETURN,
							},
							NumParams:   0,
							NumUpvalues: 0,
						},
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
				t.Errorf("Compile() \n code: \n %v \n want: \n %v", code, tt.code)
			}
			if !reflect.DeepEqual(constants, tt.constants) {
				t.Errorf("Compile() \n constants: \n %v \n want: \n %v", formatObjects(constants), formatObjects(tt.constants))
			}
		})
	}
}
