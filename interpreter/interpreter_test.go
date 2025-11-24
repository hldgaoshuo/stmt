package interpreter

import (
	"bytes"
	"reflect"
	"stmt/parser"
	"stmt/scanner"
	"testing"
)

func TestExpr(t *testing.T) {
	tests := []struct {
		name   string
		source string
		err    error
		want   any
	}{
		{
			name:   "unary wrong type data",
			source: `2 * (3 / -"muffin")`,
			err:    ErrOperandMustBeFloat64,
		},
		{
			name:   "1 + 1",
			source: "1 + 1",
			err:    nil,
			want:   2.0,
		},
		{
			name:   "true and false",
			source: "true and false",
			err:    nil,
			want:   false,
		},
		{
			name:   "false and true",
			source: "false and true",
			err:    nil,
			want:   false,
		},
		{
			name:   "false and false",
			source: "false and false",
			err:    nil,
			want:   false,
		},
		{
			name:   "true and true",
			source: "true and true",
			err:    nil,
			want:   true,
		},
		{
			name:   "true or false",
			source: "true or false",
			err:    nil,
			want:   true,
		},
		{
			name:   "false or true",
			source: "false or true",
			err:    nil,
			want:   true,
		},
		{
			name:   "false or false",
			source: "false or false",
			err:    nil,
			want:   false,
		},
		{
			name:   "true or true",
			source: "true or true",
			err:    nil,
			want:   true,
		},
		{
			name:   "1 > 2 and true",
			source: "1 > 2 and true",
			err:    nil,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := scanner.New(tt.source)
			tokens := s.ScanTokens()
			p := parser.New(tokens)
			tree, err := p.Expression()
			if err != nil {
				t.Errorf("Parse() err = %v", err)
				return
			}
			env := newEnvironment(nil)
			got, err := interpreter(tree, env)
			if err != tt.err {
				t.Errorf("Interpreter() got err = %v, want err = %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interpreter() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestStmtAndDecl(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		err        error
		wantOutput string
	}{
		{
			name:       "print string",
			source:     `print "one";`,
			err:        nil,
			wantOutput: `"one"` + "\n",
		},
		{
			name:       "print bool",
			source:     "print true;",
			err:        nil,
			wantOutput: `true` + "\n",
		},
		{
			name:       "print expr",
			source:     "print 2 + 1;",
			err:        nil,
			wantOutput: `3` + "\n",
		},
		{
			name: "var after print",
			source: `
			print a;
			var a = "too late!";
			`,
			err: ErrUndefinedVariable,
		},
		{
			name: "var nil",
			source: `
			var a;
			print a; // "nil".
			`,
			err:        nil,
			wantOutput: `<nil>` + "\n",
		},
		{
			name: "var",
			source: `
			var a = 1;
			var b = 2;
			print a + b;
			`,
			err:        nil,
			wantOutput: `3` + "\n",
		},
		{
			name: "assign",
			source: `
			var a = 1;
			print a = 2; // "2".
			`,
			err:        nil,
			wantOutput: `2` + "\n",
		},
		{
			name: "block",
			source: `
			{
				var a = "first";
				print a; // "first".
			}
			{
				var a = "second";
				print a; // "second".
			}
			`,
			err:        nil,
			wantOutput: `"first"` + "\n" + `"second"` + "\n",
		},
		{
			name: "block 2",
			source: `
			{
				var a = "in block";
			}
			print a;
			`,
			err: ErrUndefinedVariable,
		},
		{
			name: "block 3",
			source: `
			var volume = 11;
			volume = 0;
			{
				var volume = 3 * 4 * 5;
				print volume;
			}
			`,
			err:        nil,
			wantOutput: `60` + "\n",
		},
		{
			name: "block 4",
			source: `
			var global = "outside";
			{
				var local = "inside";
				print global + local;
			}
			`,
			err:        nil,
			wantOutput: `"outsideinside"` + "\n",
		},
		{
			name: "block 5",
			source: `
			var a = 1;
			{
				var a = a + 2;
				print a;
			}
			print a;
			`,
			err:        nil,
			wantOutput: `3` + "\n" + `1` + "\n",
		},
		{
			name: "for",
			source: `
			var a = 0;
			var temp;
			for (var b = 1; a < 10; b = temp + b) {
				print a;
				temp = a;
				a = b;
			}
			`,
			err:        nil,
			wantOutput: `0` + "\n" + `1` + "\n" + `1` + "\n" + `2` + "\n" + `3` + "\n" + `5` + "\n" + `8` + "\n",
		},
		{
			name: "call",
			source: `
			fun count(n) {
				if (n > 1) count(n - 1);
				print n;
			}
			count(2);
			`,
			err:        nil,
			wantOutput: `1` + "\n" + `2` + "\n",
		},
		{
			name: "call 2",
			source: `
			fun add(a, b, c) {
				print a + b + c;
			}
			add(1, 2, 3);
			`,
			err:        nil,
			wantOutput: `6` + "\n",
		},
		// {
		// 	name: "builtin",
		// 	source: `
		// 	fun time() {
		// 		print clock();
		// 	}
		// 	time();
		// 	`,
		// 	err:        nil,
		// 	wantOutput: `6` + "\n",
		// },
		{
			name: "return",
			source: `
			fun fib(n) {
				if (n <= 1) return n;
				return fib(n - 2) + fib(n - 1);
			}
			for (var i = 0; i < 5; i = i + 1) {
				print fib(i);
			}
			`,
			err:        nil,
			wantOutput: `0` + "\n" + `1` + "\n" + `1` + "\n" + `2` + "\n" + `3` + "\n",
		},
		{
			name: "closure",
			source: `
			fun makeCounter() {
				var i = 0;
				fun count() {
					i = i + 1;
					print i;
				}
				return count;
			}
			var counter = makeCounter();
			counter();
			counter();
			`,
			err:        nil,
			wantOutput: `1` + "\n" + `2` + "\n",
		},
		{
			name: "closure 2",
			source: `
			var a = "global";
			{
				fun showA() {
					print a;
				}
				showA();
				var a = "block";
				showA();
			}
			`,
			err:        nil,
			wantOutput: `"global"` + "\n" + `"global"` + "\n",
		},
		{
			name: "class field",
			source: `
			class SomeObject {
			}
			var someObject = SomeObject();
			someObject.someProperty = "value";
			print someObject.someProperty;
			`,
			err:        nil,
			wantOutput: `"value"` + "\n",
		},
		{
			name: "class method",
			source: `
			class Bacon {
				eat() {
					print "Crunch crunch crunch!";
				}
			}
			var bacon = Bacon();
			bacon.eat();
			`,
			err:        nil,
			wantOutput: `"Crunch crunch crunch!"` + "\n",
		},
		{
			name: "class this",
			source: `
			class Cake {
				taste() {
					var adjective = "delicious";
					print "The " + this.flavor + " cake is " + adjective + "!";
				}
			}
			var cake = Cake();
			cake.flavor = "German chocolate";
			cake.taste(); // Prints "The German chocolate cake is delicious!".
			`,
			err:        nil,
			wantOutput: `"The German chocolate cake is delicious!"` + "\n",
		},
		{
			name:       "class this 2",
			source:     "print this;",
			err:        ErrUndefinedVariable,
			wantOutput: "",
		},
		{
			name: "class superclass",
			source: `
			class Doughnut {
				cook() {
					print "Fry until golden brown.";
				}
			}
			class BostonCream < Doughnut {
			}
			BostonCream().cook();
			`,
			err:        nil,
			wantOutput: `"Fry until golden brown."` + "\n",
		},
		{
			name: "class super",
			source: `
			class A {
				method() {
					print "A method";
				}
			}
			
			class B < A {
				method() {
					print "B method";
				}
				test() {
					super.method();
				}
			}
			
			class C < B {
			}
			
			C().test();
			`,
			err:        nil,
			wantOutput: `"A method"` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个缓冲区来捕获输出
			var buf bytes.Buffer
			Output = &buf

			s := scanner.New(tt.source)
			tokens := s.ScanTokens()
			p := parser.New(tokens)
			tree, err := p.Parse()
			if err != nil {
				t.Errorf("Parse() err = %v", err)
				return
			}
			err = Interpreter(tree)
			if err != tt.err {
				t.Errorf("Interpreter() got err = %v, want err = %v", err, tt.err)
				return
			}
			if tt.wantOutput != "" && buf.String() != tt.wantOutput {
				t.Errorf("Interpreter() output = %q, want %q", buf.String(), tt.wantOutput)
				return
			}
		})
	}
}
