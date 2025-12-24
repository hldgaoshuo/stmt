package vm

import (
	"bytes"
	"errors"
	"reflect"
	"stmt/ast"
	"stmt/compiler"
	"stmt/parser"
	"stmt/scanner"
	"stmt/value"
	"testing"
)

func TestVM_RunExpr(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		err      error
		result   any
		stackLen int
	}{
		{
			name:     "1",
			source:   "1",
			err:      nil,
			result:   value.NewInt(1),
			stackLen: 1,
		},
		{
			name:     "1.5",
			source:   "1.5",
			err:      nil,
			result:   value.NewFloat(1.5),
			stackLen: 1,
		},
		{
			name:     `"abc"`,
			source:   `"abc"`,
			err:      nil,
			result:   value.NewString("abc"),
			stackLen: 1,
		},
		{
			name:     "true",
			source:   "true",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "false",
			source:   "false",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "nil",
			source:   "nil",
			err:      nil,
			result:   value.NewNil(),
			stackLen: 1,
		},
		{
			name:     "-1",
			source:   "-1",
			err:      nil,
			result:   value.NewInt(-1),
			stackLen: 1,
		},
		{
			name:     "1+2",
			source:   "1+2",
			err:      nil,
			result:   value.NewInt(3),
			stackLen: 1,
		},
		{
			name:     "1-2",
			source:   "1-2",
			err:      nil,
			result:   value.NewInt(-1),
			stackLen: 1,
		},
		{
			name:     "1*2",
			source:   "1*2",
			err:      nil,
			result:   value.NewInt(2),
			stackLen: 1,
		},
		{
			name:     "1_/_2",
			source:   "1/2",
			err:      nil,
			result:   value.NewInt(0),
			stackLen: 1,
		},
		{
			name:     "!true",
			source:   "!true",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "!false",
			source:   "!false",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1==1",
			source:   "1==1",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1==1.0",
			source:   "1==1.0",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1.0==1",
			source:   "1.0==1",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1.0==1.0",
			source:   "1.0==1.0",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1==2",
			source:   "1==2",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "1.0==2",
			source:   "1.0==2",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "1==2.0",
			source:   "1==2.0",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "1.0==2.0",
			source:   "1.0==2.0",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "true==true",
			source:   "true==true",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "false==false",
			source:   "false==false",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "true==false",
			source:   "true==false",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "false==true",
			source:   "false==true",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "nil==nil",
			source:   "nil==nil",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1>2",
			source:   "1>2",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "1<2",
			source:   "1<2",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "1>=2",
			source:   "1>=2",
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "1<=2",
			source:   "1<=2",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     `"a"+"b"`,
			source:   `"a"+"b"`,
			err:      nil,
			result:   value.NewString("ab"),
			stackLen: 1,
		},
		{
			name:     "true_and_true",
			source:   "true and true",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "true_and_false",
			source:   `true and false`,
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "false_and_true",
			source:   `false and true`,
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "false_and_false",
			source:   `false and false`,
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
		},
		{
			name:     "true_or_true",
			source:   "true or true",
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "true_or_false",
			source:   `true or false`,
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "false_or_true",
			source:   `false or true`,
			err:      nil,
			result:   value.NewBool(true),
			stackLen: 1,
		},
		{
			name:     "false_or_false",
			source:   `false or false`,
			err:      nil,
			result:   value.NewBool(false),
			stackLen: 1,
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
			compiler_ := compiler.New([]ast.Node{node})
			code, constants, err := compiler_.Compile()
			if err != nil {
				t.Errorf("Compile() err = %v", err)
				return
			}
			vm := New(code, constants, len(compiler.Global.LocalValues))
			err = vm.Run()
			if !errors.Is(err, tt.err) {
				t.Errorf("Run() err = %v, want %v", err, tt.err)
			}
			result := vm.StackPeek(0)
			if !reflect.DeepEqual(result, tt.result) {
				t.Errorf("result = %v, want %v", result, tt.result)
			}
			stackLen := vm.StackLen()
			if stackLen != uint64(tt.stackLen) {
				t.Errorf("stackLen = %v, want %v", stackLen, tt.stackLen)
			}
		})
	}
}

func TestVM_RunStmtDecl(t *testing.T) {
	tests := []struct {
		name   string
		source string
		err    error
		result string
	}{
		{
			name: "print",
			source: `
			print 123;
			`,
			err:    nil,
			result: "123" + "\n",
		},
		{
			name: "var",
			source: `
			var a = 123;
			print a;
			`,
			err:    nil,
			result: "123" + "\n",
		},
		{
			name: "var_2",
			source: `
			var a = 1;
			print a;
			{
				var a = 2;
				print a;
			}
			print a;
			`,
			err:    nil,
			result: "1" + "\n" + "2" + "\n" + "1" + "\n",
		},
		{
			name: "if",
			source: `
			if (true)
			{
				print 1;
			}
			print 2;
			`,
			err:    nil,
			result: "1" + "\n" + "2" + "\n",
		},
		{
			name: "if_else",
			source: `
			if (false)
			{
				print 1;
			}
			else
			{
				print 2;
			}
			`,
			err:    nil,
			result: "2" + "\n",
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
			err:    nil,
			result: "0" + "\n" + "1" + "\n" + "2" + "\n" + "3" + "\n" + "4" + "\n",
		},
		{
			name: "call",
			source: `
			fun pt() {
				print 1;
			}
			pt();
			`,
			err:    nil,
			result: "1" + "\n",
		},
		{
			name: "call_arg",
			source: `
			fun pt(a, b) {
				print a + b;
			}
			pt(1, 2);
			`,
			err:    nil,
			result: "3" + "\n",
		},
		{
			name: "call_arg_return",
			source: `
			fun add(a, b) {
				return a + b;
			}
			print add(1, 2);
			`,
			err:    nil,
			result: "3" + "\n",
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
			err:    nil,
			result: "outside" + "\n",
		},
		{
			name: "closure_2",
			source: `
			var x = "global";
			fun outer() {
				var x = "outer";
				fun inner() {
					print x;
				}
				inner();
			}
			outer();
			`,
			err:    nil,
			result: "outer" + "\n",
		},
		{
			name: "closure_3",
			source: `
			fun outer() {
				var x = 1;
				fun middle() {
					fun inner() {
						print x;
					}
					inner();
				}
				middle();
			}
			outer();
			`,
			err:    nil,
			result: "1" + "\n",
		},
		{
			name: "closure_4",
			source: `
			fun makeClosure() {
				var local = "local";
				fun closure() {
					print local;
				}
				return closure;
			}
			var closure = makeClosure();
			closure();
			`,
			err:    nil,
			result: "local" + "\n",
		},
		{
			name: "closure_5",
			source: `
			fun outer() {
				var x = "value";
				fun middle() {
					fun inner() {
						print x;
					}
					print "create inner closure";
					return inner;
				}
				print "return from outer";
				return middle;
			}
			var mid = outer();
			var in = mid();
			in();
			`,
			err:    nil,
			result: "return from outer" + "\n" + "create inner closure" + "\n" + "value" + "\n",
		},
		{
			name: "closure_6",
			source: `
			fun outer() {
				var x = "before";
				fun inner() {
					x = "assigned";
				}
				inner();
				print x;
			}
			outer();
			`,
			err:    nil,
			result: "assigned" + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个缓冲区来捕获输出
			var buf bytes.Buffer
			Output = &buf

			scanner_ := scanner.New(tt.source)
			tokens := scanner_.Scan()
			parser_ := parser.New(tokens)
			node, err := parser_.Parse()
			if err != nil {
				t.Errorf("Parse() err = %v", err)
				return
			}
			compiler_ := compiler.New(node)
			code, constants, err := compiler_.Compile()
			if err != nil {
				t.Errorf("Compile() err = %v", err)
				return
			}
			vm := New(code, constants, len(compiler.Global.LocalValues))
			err = vm.Run()
			if !errors.Is(err, tt.err) {
				t.Errorf("Run() err = %v, want %v", err, tt.err)
			}
			if tt.result != "" && buf.String() != tt.result {
				t.Errorf("Run() output = %q, want %q", buf.String(), tt.result)
			}
		})
	}
}
