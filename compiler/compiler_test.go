package compiler

import (
	"errors"
	"fmt"
	"reflect"
	"stmt/ast"
	"stmt/parser"
	"stmt/scanner"
	"stmt/value"
	"testing"
)

func formatConstants(constants []value.Value) string {
	if len(constants) == 0 {
		return "[]"
	}
	result := "[\n"
	for i, constant := range constants {
		result += fmt.Sprintf("   [%d] %s", i, constant.String())
		if i < len(constants)-1 {
			result += ",\n"
		} else {
			result += "\n"
		}
	}
	result += " ]"
	return result
}

func TestCompiler_CompileExpr(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		err       error
		code      []uint8
		constants []value.Value
	}{
		{
			name:   "1",
			source: "1",
			code: []uint8{
				OP_CONSTANT, 0,
			},
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
			},
		},
		{
			name:   "1.2",
			source: "1.2",
			code: []uint8{
				OP_CONSTANT, 0,
			},
			constants: []value.Value{
				&value.Float{
					Literal: 1.2,
				},
			},
		},
		{
			name:   "(1)",
			source: "(1)",
			code: []uint8{
				OP_CONSTANT, 0,
			},
			constants: []value.Value{
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
				},
			},
		},
		{
			name:   "true",
			source: "true",
			code: []uint8{
				OP_TRUE,
			},
			constants: []value.Value{},
		},
		{
			name:   "false",
			source: "false",
			code: []uint8{
				OP_FALSE,
			},
			constants: []value.Value{},
		},
		{
			name:   "nil",
			source: "nil",
			code: []uint8{
				OP_NIL,
			},
			constants: []value.Value{},
		},
		{
			name:   "!true",
			source: "!true",
			code: []uint8{
				OP_TRUE,
				OP_NOT,
			},
			constants: []value.Value{},
		},
		{
			name:   "1<2",
			source: "1<2",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_CONSTANT, 1,
				OP_LT,
			},
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
				},
			},
		},
		{
			name:   `"abc"`,
			source: `"abc"`,
			code: []uint8{
				OP_CONSTANT, 0,
			},
			constants: []value.Value{
				&value.String{
					Literal: "abc",
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
			if !errors.Is(err, tt.err) {
				t.Errorf("Compile() err = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(code, tt.code) {
				t.Errorf("Compile() code = %v, want %v", code, tt.code)
			}
			if !reflect.DeepEqual(constants, tt.constants) {
				t.Errorf("\n Compile() constants: \n %v \n want: \n %v", formatConstants(constants), formatConstants(tt.constants))
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
		constants []value.Value
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{},
		},
		{
			name:   "var 2",
			source: "var a = 1;",
			code: []uint8{
				OP_CONSTANT, 0,
				OP_SET_GLOBAL, 0,
			},
			constants: []value.Value{
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.Int{
					Literal: 10,
				},
				&value.Int{
					Literal: 20,
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
			constants: []value.Value{
				&value.Int{
					Literal: 10,
				},
				&value.Int{
					Literal: 20,
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
			constants: []value.Value{},
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
			constants: []value.Value{},
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
			constants: []value.Value{
				&value.Int{
					Literal: 0,
				},
				&value.Int{
					Literal: 5,
				},
				&value.Int{
					Literal: 1,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Function{
					Code: []uint8{
						OP_CONSTANT, 0,
						OP_PRINT,
						OP_NIL,
						OP_RETURN,
					},
					NumParams:   0,
					NumUpvalues: 0,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Function{
					Code: []uint8{
						OP_CONSTANT, 0,
						OP_PRINT,
						OP_NIL,
						OP_RETURN,
					},
					NumParams:   0,
					NumUpvalues: 0,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
				},
				&value.Function{
					Code: []uint8{
						OP_CONSTANT, 0,
						OP_PRINT,
						OP_CONSTANT, 1,
						OP_RETURN,
					},
					NumParams:   0,
					NumUpvalues: 0,
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Function{
					Code: []uint8{
						OP_CONSTANT, 0,
						OP_PRINT,
						OP_NIL,
						OP_RETURN,
					},
					NumParams:   0,
					NumUpvalues: 0,
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
			constants: []value.Value{
				&value.Function{
					Code: []uint8{
						OP_GET_LOCAL, 0,
						OP_GET_LOCAL, 1,
						OP_ADD,
						OP_PRINT,
						OP_NIL,
						OP_RETURN,
					},
					NumParams:   2,
					NumUpvalues: 0,
				},
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.Function{
					Code: []uint8{
						OP_GET_LOCAL, 0,
						OP_GET_LOCAL, 1,
						OP_ADD,
						OP_RETURN,
					},
					NumParams:   2,
					NumUpvalues: 0,
				},
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
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
			constants: []value.Value{
				&value.String{
					Literal: "outside",
				},
				&value.Function{
					Code: []uint8{
						OP_GET_UPVALUE, 0,
						OP_PRINT,
						OP_NIL,
						OP_RETURN,
					},
					NumParams:   0,
					NumUpvalues: 1,
				},
				&value.Function{
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
			constants: []value.Value{
				&value.Int{
					Literal: 1,
				},
				&value.Int{
					Literal: 2,
				},
				&value.Int{
					Literal: 3,
				},
				&value.Int{
					Literal: 4,
				},
				&value.Function{
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
				&value.Function{
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
				&value.Function{
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
				t.Errorf("\n Compile() constants: \n %v \n want: \n %v", formatConstants(constants), formatConstants(tt.constants))
			}
		})
	}
}
