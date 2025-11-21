package ast

import (
	"log/slog"
	"stmt/token"
)

type Var struct {
	Line        int
	Name        *token.Token
	Initializer Node
}

type Function struct {
	Line   int
	Name   *token.Token
	Params []*token.Token
	Body   Node
}

func (f *Function) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("Line", f.Line),
	)
}
