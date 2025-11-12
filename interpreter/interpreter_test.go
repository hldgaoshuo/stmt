package interpreter

import (
	"bytes"
	"log/slog"
	"os"
	"reflect"
	"stmt/parser"
	"stmt/scanner"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		want       any
		err        error
		wantOutput string
	}{
		{
			name:   "unary wrong type data",
			source: `2 * (3 / -"muffin");`,
			want:   nil,
			err:    ErrOperandMustBeFloat64,
		},
		{
			name:   "1 + 1",
			source: "1 + 1;",
			want:   2.0,
			err:    nil,
		},
		{
			name:       "print string",
			source:     `print "one";`,
			want:       nil,
			err:        nil,
			wantOutput: `"one"` + "\n",
		},
		{
			name:       "print bool",
			source:     "print true;",
			want:       nil,
			err:        nil,
			wantOutput: `true` + "\n",
		},
		{
			name:       "print expr",
			source:     "print 2 + 1;",
			want:       nil,
			err:        nil,
			wantOutput: `3` + "\n",
		},
		{
			name: "var after print",
			source: `
			print a;
			var a = "too late!";
			`,
			want: nil,
			err:  ErrUndefinedVariable,
		},
		{
			name: "var nil",
			source: `
			var a;
			print a; // "nil".
			`,
			want:       nil,
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
			want:       nil,
			err:        nil,
			wantOutput: `3` + "\n",
		},
		{
			name: "assign",
			source: `
			var a = 1;
			print a = 2; // "2".
			`,
			want:       nil,
			err:        nil,
			wantOutput: `2` + "\n",
		},
	}

	// 自定义 ReplaceAttr 去除 time 字段
	removeTime := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{} // 返回空表示不输出
		}
		return a
	}

	// 创建 handler
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: removeTime,
	})

	// 设置全局 logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

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
			got, err := Interpreter(tree)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interpreter() = %v, want %v", got, tt.want)
				return
			}
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
