package interpreter

import (
	"stmt/ast"
)

type closure struct {
	Function *ast.Function
	Env      *environment
}
