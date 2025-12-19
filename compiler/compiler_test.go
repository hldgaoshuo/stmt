package compiler

import (
	"errors"
	"fmt"
	"reflect"
	"stmt/ast"
	"stmt/opcode"
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

func newCode(codeMatrix ...[]uint8) []uint8 {
	var result []uint8
	for _, codeList := range codeMatrix {
		result = append(result, codeList...)
	}
	return result
}

func toCode(opcode uint8, operand ...uint64) []uint8 {
	return []uint8{opcode}
}

func newClosureMeta(metas ...uint8) []uint8 {
	return metas
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
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
			),
			constants: []value.Value{
				value.NewInt(1),
			},
			err: nil,
		},
		{
			name:   "1.2",
			source: "1.2",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
			),
			constants: []value.Value{
				value.NewFloat(1.2),
			},
			err: nil,
		},
		{
			name:   "(1)",
			source: "(1)",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
			),
			constants: []value.Value{
				value.NewInt(1),
			},
			err: nil,
		},
		{
			name:   "-1",
			source: "-1",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				toCode(opcode.OP_NEGATE),
			),
			constants: []value.Value{
				value.NewInt(1),
			},
			err: nil,
		},
		{
			name:   "1+2",
			source: "1+2",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_CONSTANT, 1),
				toCode(opcode.OP_ADD),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
			},
			err: nil,
		},
		{
			name:   "true",
			source: "true",
			code: newCode(
				toCode(opcode.OP_TRUE),
			),
			constants: []value.Value{},
			err:       nil,
		},
		{
			name:   "false",
			source: "false",
			code: newCode(
				toCode(opcode.OP_FALSE),
			),
			constants: []value.Value{},
			err:       nil,
		},
		{
			name:   "nil",
			source: "nil",
			code: newCode(
				toCode(opcode.OP_NIL),
			),
			constants: []value.Value{},
			err:       nil,
		},
		{
			name:   "!true",
			source: "!true",
			code: newCode(
				toCode(opcode.OP_TRUE),
				toCode(opcode.OP_NOT),
			),
			constants: []value.Value{},
			err:       nil,
		},
		{
			name:   "1<2",
			source: "1<2",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_CONSTANT, 1),
				toCode(opcode.OP_LT),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
			},
			err: nil,
		},
		{
			name:   `"abc"`,
			source: `"abc"`,
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
			),
			constants: []value.Value{
				value.NewString("abc"),
			},
			err: nil,
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
		code      []uint8
		constants []value.Value
	}{
		{
			name:   "expr",
			source: "1+2;",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_CONSTANT, 1),
				toCode(opcode.OP_ADD),
				toCode(opcode.OP_POP),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
			},
		},
		{
			name:   "print",
			source: "print 1;",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				toCode(opcode.OP_PRINT),
			),
			constants: []value.Value{
				value.NewInt(1),
			},
		},
		{
			name:   "var",
			source: "var a;",
			code: newCode(
				toCode(opcode.OP_NIL),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
			),
			constants: []value.Value{},
		},
		{
			name:   "var 2",
			source: "var a = 1;",
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
			),
			constants: []value.Value{
				value.NewInt(1),
			},
		},
		{
			name: "var 3",
			source: `
			var a = 1;
			a;
			`,
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_GET_GLOBAL, 0),
				toCode(opcode.OP_POP),
			),
			constants: []value.Value{
				value.NewInt(1),
			},
		},
		{
			name: "var 4",
			source: `
			var a = 1;
			print a;
			`,
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_GET_GLOBAL, 0),
				toCode(opcode.OP_PRINT),
			),
			constants: []value.Value{
				value.NewInt(1),
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
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_GET_GLOBAL, 0),
				toCode(opcode.OP_PRINT),
				CodeMake(opcode.OP_CONSTANT, 1),
				CodeMake(opcode.OP_SET_LOCAL, 0),
				CodeMake(opcode.OP_GET_LOCAL, 0),
				toCode(opcode.OP_PRINT),
				CodeMake(opcode.OP_GET_GLOBAL, 0),
				toCode(opcode.OP_PRINT),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
			},
		},
		{
			name: "assign",
			source: `
			var a = 1;
			a = 2;
			`,
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_CONSTANT, 1),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
			},
		},
		{
			name: "assign 2",
			source: `
			var a = 1;
			a = 2;
			print a;
			`,
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_CONSTANT, 1),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_GET_GLOBAL, 0),
				toCode(opcode.OP_PRINT),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
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
			code: newCode(
				CodeMake(opcode.OP_CONSTANT, 0),
				CodeMake(opcode.OP_SET_GLOBAL, 0),
				CodeMake(opcode.OP_CONSTANT, 1),
				CodeMake(opcode.OP_SET_LOCAL, 0),
			),
			constants: []value.Value{
				value.NewInt(1),
				value.NewInt(2),
			},
		},
		// {
		// 	name: "if",
		// 	source: `
		// 	if (true)
		// 	{
		// 		print 10;
		// 	}
		// 	print 20;
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_TRUE),
		// 		CodeMake(OP_JUMP_FALSE, 9),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_CONSTANT, 0),
		// 		CodeMake(OP_PRINT),
		// 		CodeMake(OP_JUMP, 10),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_CONSTANT, 1),
		// 		CodeMake(OP_PRINT),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(10),
		// 		value.NewInt(20),
		// 	},
		// },
		// {
		// 	name: "if else",
		// 	source: `
		// 	if (false)
		// 	{
		// 		print 10;
		// 	}
		// 	else
		// 	{
		// 		print 20;
		// 	}
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_FALSE),
		// 		CodeMake(OP_JUMP_FALSE, 9),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_CONSTANT, 0),
		// 		CodeMake(OP_PRINT),
		// 		CodeMake(OP_JUMP, 13),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_CONSTANT, 1),
		// 		CodeMake(OP_PRINT),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(10),
		// 		value.NewInt(20),
		// 	},
		// },
		// {
		// 	name: "and",
		// 	source: `
		// 	true and true;
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_TRUE),
		// 		CodeMake(OP_JUMP_FALSE, 5),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_TRUE),
		// 		CodeMake(OP_POP),
		// 	),
		// 	constants: []value.Value{},
		// },
		// {
		// 	name: "or",
		// 	source: `
		// 	true or true;
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_TRUE),
		// 		CodeMake(OP_JUMP_FALSE, 5),
		// 		CodeMake(OP_JUMP, 7),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_TRUE),
		// 		CodeMake(OP_POP),
		// 	),
		// 	constants: []value.Value{},
		// },
		// {
		// 	name: "while",
		// 	source: `
		// 	var i = 0;
		// 	while (i < 5)
		// 	{
		// 		print i;
		// 		i = i + 1;
		// 	}
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CONSTANT, 0),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_CONSTANT, 1),
		// 		CodeMake(OP_LT),
		// 		CodeMake(OP_JUMP_FALSE, 24),
		// 		CodeMake(OP_POP),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_PRINT),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_CONSTANT, 2),
		// 		CodeMake(OP_ADD),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 		CodeMake(OP_LOOP, 4),
		// 		CodeMake(OP_POP),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(0),
		// 		value.NewInt(5),
		// 		value.NewInt(1),
		// 	},
		// },
		// {
		// 	name: "function",
		// 	source: `
		// 	fun pt() {
		// 		print 1;
		// 	}
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 1),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(1),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 0),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 0),
		// 	},
		// },
		// {
		// 	name: "function return nil",
		// 	source: `
		// 	fun pt() {
		// 		print 1;
		// 		return;
		// 	}
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 1),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(1),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 0),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 0),
		// 	},
		// },
		// {
		// 	name: "function return value",
		// 	source: `
		// 	fun pt() {
		// 		print 1;
		// 		return 2;
		// 	}
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 2),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(1),
		// 		value.NewInt(2),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 0),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_CONSTANT, 1),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 0),
		// 	},
		// },
		// {
		// 	name: "call",
		// 	source: `
		// 	fun pt() {
		// 		print 1;
		// 	}
		// 	pt();
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 1),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_CALL, 0),
		// 		CodeMake(OP_POP),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(1),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 0),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 0),
		// 	},
		// },
		// {
		// 	name: "call arg",
		// 	source: `
		// 	fun pt(a, b) {
		// 		print a + b;
		// 	}
		// 	pt(1, 2);
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 0),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_CONSTANT, 1),
		// 		CodeMake(OP_CONSTANT, 2),
		// 		CodeMake(OP_CALL, 2),
		// 		CodeMake(OP_POP),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_GET_LOCAL, 0),
		// 			CodeMake(OP_GET_LOCAL, 1),
		// 			CodeMake(OP_ADD),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 2, 0),
		// 		value.NewInt(1),
		// 		value.NewInt(2),
		// 	},
		// },
		// {
		// 	name: "call arg return",
		// 	source: `
		// 	fun add(a, b) {
		// 		return a + b;
		// 	}
		// 	print add(1, 2);
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 0),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_CONSTANT, 1),
		// 		CodeMake(OP_CONSTANT, 2),
		// 		CodeMake(OP_CALL, 2),
		// 		CodeMake(OP_PRINT),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_GET_LOCAL, 0),
		// 			CodeMake(OP_GET_LOCAL, 1),
		// 			CodeMake(OP_ADD),
		// 			CodeMake(OP_RETURN),
		// 		), 2, 0),
		// 		value.NewInt(1),
		// 		value.NewInt(2),
		// 	},
		// },
		// {
		// 	name: "closure",
		// 	source: `
		// 	fun outer() {
		// 		var x = "outside";
		// 		fun inner() {
		// 			print x;
		// 		}
		// 		inner();
		// 	}
		// 	outer();
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 2),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 		CodeMake(OP_GET_GLOBAL, 0),
		// 		CodeMake(OP_CALL, 0),
		// 		CodeMake(OP_POP),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewString("outside"),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_GET_UPVALUE, 0),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 1),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 0),
		// 			CodeMake(OP_SET_LOCAL, 0),
		// 			CodeMake(OP_CLOSURE, 1),
		// 			newClosureMeta(1, 0),
		// 			CodeMake(OP_SET_LOCAL, 1),
		// 			CodeMake(OP_GET_LOCAL, 1),
		// 			CodeMake(OP_CALL, 0),
		// 			CodeMake(OP_POP),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 0),
		// 	},
		// },
		// {
		// 	name: "closure 1",
		// 	source: `
		// 	fun outer() {
		// 		var a = 1;
		// 		var b = 2;
		// 		fun middle() {
		// 			var c = 3;
		// 			var d = 4;
		// 			fun inner() {
		// 				print a + c + b + d;
		// 			}
		// 		}
		// 	}
		// 	`,
		// 	code: newCode(
		// 		CodeMake(OP_CLOSURE, 6),
		// 		CodeMake(OP_SET_GLOBAL, 0),
		// 	),
		// 	constants: []value.Value{
		// 		value.NewInt(1),
		// 		value.NewInt(2),
		// 		value.NewInt(3),
		// 		value.NewInt(4),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_GET_UPVALUE, 0),
		// 			CodeMake(OP_GET_UPVALUE, 1),
		// 			CodeMake(OP_ADD),
		// 			CodeMake(OP_GET_UPVALUE, 2),
		// 			CodeMake(OP_ADD),
		// 			CodeMake(OP_GET_UPVALUE, 3),
		// 			CodeMake(OP_ADD),
		// 			CodeMake(OP_PRINT),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 4),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 2),
		// 			CodeMake(OP_SET_LOCAL, 0),
		// 			CodeMake(OP_CONSTANT, 3),
		// 			CodeMake(OP_SET_LOCAL, 1),
		// 			CodeMake(OP_CLOSURE, 4),
		// 			newClosureMeta(0, 0, 1, 0, 0, 1, 1, 1),
		// 			CodeMake(OP_SET_LOCAL, 2),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 2),
		// 		value.NewFunction(newCode(
		// 			CodeMake(OP_CONSTANT, 0),
		// 			CodeMake(OP_SET_LOCAL, 0),
		// 			CodeMake(OP_CONSTANT, 1),
		// 			CodeMake(OP_SET_LOCAL, 1),
		// 			CodeMake(OP_CLOSURE, 5),
		// 			newClosureMeta(1, 0, 1, 1),
		// 			CodeMake(OP_SET_LOCAL, 2),
		// 			CodeMake(OP_NIL),
		// 			CodeMake(OP_RETURN),
		// 		), 0, 0),
		// 	},
		// },
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
