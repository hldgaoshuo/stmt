package vm

import "stmt/value"

type Closure struct {
	Function *value.Function
}
