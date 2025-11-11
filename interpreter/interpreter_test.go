package interpreter

import (
	"log/slog"
	"os"
	"reflect"
	"stmt/parser"
	"stmt/scanner"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
		err    error
	}{
		{
			name:   `2 * (3 / -"muffin")`,
			source: `2 * (3 / -"muffin")`,
			want:   nil,
			err:    ErrOperandMustBeFloat64,
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
		})
	}
}
