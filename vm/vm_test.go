package vm

import (
	"errors"
	"stmt/ast"
	"stmt/compiler"
	"stmt/parser"
	"stmt/scanner"
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
			result:   int64(1),
			stackLen: 1,
		},
		{
			name:     "1.5",
			source:   "1.5",
			err:      nil,
			result:   1.5,
			stackLen: 1,
		},
		{
			name:     `"abc"`,
			source:   `"abc"`,
			err:      nil,
			result:   "abc",
			stackLen: 1,
		},
		{
			name:     "true",
			source:   "true",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "false",
			source:   "false",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "nil",
			source:   "nil",
			err:      nil,
			result:   nil,
			stackLen: 1,
		},
		{
			name:     "-1",
			source:   "-1",
			err:      nil,
			result:   int64(-1),
			stackLen: 1,
		},
		{
			name:     "1+2",
			source:   "1+2",
			err:      nil,
			result:   int64(3),
			stackLen: 1,
		},
		{
			name:     "1-2",
			source:   "1-2",
			err:      nil,
			result:   int64(-1),
			stackLen: 1,
		},
		{
			name:     "1*2",
			source:   "1*2",
			err:      nil,
			result:   int64(2),
			stackLen: 1,
		},
		{
			name:     "1_/_2",
			source:   "1/2",
			err:      nil,
			result:   int64(0),
			stackLen: 1,
		},
		{
			name:     "!true",
			source:   "!true",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "!false",
			source:   "!false",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1==1",
			source:   "1==1",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1==1.0",
			source:   "1==1.0",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1.0==1",
			source:   "1.0==1",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1.0==1.0",
			source:   "1.0==1.0",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1==2",
			source:   "1==2",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "1.0==2",
			source:   "1.0==2",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "1==2.0",
			source:   "1==2.0",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "1.0==2.0",
			source:   "1.0==2.0",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "true==true",
			source:   "true==true",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "false==false",
			source:   "false==false",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "true==false",
			source:   "true==false",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "false==true",
			source:   "false==true",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "nil==nil",
			source:   "nil==nil",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1>2",
			source:   "1>2",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "1<2",
			source:   "1<2",
			err:      nil,
			result:   true,
			stackLen: 1,
		},
		{
			name:     "1>=2",
			source:   "1>=2",
			err:      nil,
			result:   false,
			stackLen: 1,
		},
		{
			name:     "1<=2",
			source:   "1<=2",
			err:      nil,
			result:   true,
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
			vm := New(code, constants)
			err = vm.Run()
			if !errors.Is(err, tt.err) {
				t.Errorf("Run() err = %v, want %v", err, tt.err)
			}
			result := vm.StackPeek(0)
			if result != tt.result {
				t.Errorf("result = %v, want %v", result, tt.result)
			}
			stackLen := vm.StackLen()
			if stackLen != tt.stackLen {
				t.Errorf("stackLen = %v, want %v", stackLen, tt.stackLen)
			}
		})
	}
}
