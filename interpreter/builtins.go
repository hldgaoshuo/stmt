package interpreter

import "time"

type builtin func(args ...any) (any, error)

var builtins = map[string]builtin{
	"clock": clock,
}

func clock(args ...any) (any, error) {
	return time.Now().Unix(), nil
}
