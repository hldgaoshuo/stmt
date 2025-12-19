package vm

import (
	"math"
	"stmt/opcode"
	"stmt/value"
)

type VM struct {
	Stack     []any
	Globals   []any
	Frames    []*Frame
	Constants []value.Value
}

func New(code []uint8, constants []value.Value) *VM {
	mainFunction := &value.Function{
		Code:        code,
		NumParams:   0,
		NumUpvalues: 0,
	}
	mainClosure := &Closure{
		Function: mainFunction,
	}
	mainFrame := NewFrame(mainClosure, 0)
	return &VM{
		Stack:     []any{},
		Globals:   []any{},
		Frames:    []*Frame{mainFrame},
		Constants: constants,
	}
}

func (vm *VM) FramesBack() *Frame {
	return vm.Frames[len(vm.Frames)-1]
}

func (vm *VM) FramesPush(frame *Frame) {
	vm.Frames = append(vm.Frames, frame)
}

func (vm *VM) FramesPop() *Frame {
	vm.Frames = vm.Frames[:len(vm.Frames)-1]
	return vm.FramesBack()
}

func (vm *VM) StackPush(value any) {
	vm.Stack = append(vm.Stack, value)
}

func (vm *VM) StackPop() any {
	offset := len(vm.Stack) - 1
	value_ := vm.Stack[offset]
	vm.Stack = vm.Stack[:offset]
	return value_
}

func (vm *VM) StackPeek(num int) any {
	return vm.Stack[len(vm.Stack)-1-num]
}

func (vm *VM) StackSet(index int, value any) {
	vm.Stack[index] = value
}

func (vm *VM) StackGet(index int) any {
	return vm.Stack[index]
}

func (vm *VM) StackLen() int {
	return len(vm.Stack)
}

func (vm *VM) StackPushConstant(global value.Value) error {
	switch _global := global.(type) {
	case *value.Int:
		vm.StackPush(_global.Literal)
		return nil
	case *value.Float:
		vm.StackPush(_global.Literal)
		return nil
	case *value.String:
		vm.StackPush(_global.Literal)
		return nil
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushNegate(a any) error {
	switch _a := a.(type) {
	case int64:
		vm.StackPush(-_a)
		return nil
	case float64:
		vm.StackPush(-_a)
		return nil
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushAdd(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a + _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) + _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a + float64(_b))
			return nil
		case float64:
			vm.StackPush(_a + _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case string:
		switch _b := b.(type) {
		case string:
			vm.StackPush(_a + _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushSubtract(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a - _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) - _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a - float64(_b))
			return nil
		case float64:
			vm.StackPush(_a - _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushMultiply(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a * _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) * _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a * float64(_b))
			return nil
		case float64:
			vm.StackPush(_a * _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushDivide(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a / _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) / _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a / float64(_b))
			return nil
		case float64:
			vm.StackPush(_a / _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushModulo(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a % _b)
			return nil
		case float64:
			vm.StackPush(math.Mod(float64(_a), _b))
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(math.Mod(_a, float64(_b)))
			return nil
		case float64:
			vm.StackPush(math.Mod(_a, _b))
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushNot(a any) error {
	switch _a := a.(type) {
	case bool:
		vm.StackPush(!_a)
		return nil
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushEQ(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a == _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) == _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a == float64(_b))
			return nil
		case float64:
			vm.StackPush(_a == _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case bool:
		switch _b := b.(type) {
		case bool:
			vm.StackPush(_a == _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case nil:
		switch b.(type) {
		case nil:
			vm.StackPush(true)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushGT(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a > _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) > _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a > float64(_b))
			return nil
		case float64:
			vm.StackPush(_a > _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushLT(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a < _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) < _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a < float64(_b))
			return nil
		case float64:
			vm.StackPush(_a < _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushGE(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a >= _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) >= _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a >= float64(_b))
			return nil
		case float64:
			vm.StackPush(_a >= _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushLE(a any, b any) error {
	switch _a := a.(type) {
	case int64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a <= _b)
			return nil
		case float64:
			vm.StackPush(float64(_a) <= _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case float64:
		switch _b := b.(type) {
		case int64:
			vm.StackPush(_a <= float64(_b))
			return nil
		case float64:
			vm.StackPush(_a <= _b)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) Run() error {
	frame := vm.FramesBack()
	for frame.Ip < frame.CodeSize() {
		op := frame.Opcode()
		switch op {
		case opcode.OP_CONSTANT, opcode.OP_CONSTANT_2, opcode.OP_CONSTANT_4, opcode.OP_CONSTANT_8:
			globalIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			global := vm.Constants[globalIndex]
			err = vm.StackPushConstant(global)
			if err != nil {
				return err
			}
		case opcode.OP_TRUE:
			vm.StackPush(true)
		case opcode.OP_FALSE:
			vm.StackPush(false)
		case opcode.OP_NIL:
			vm.StackPush(nil)
		case opcode.OP_NEGATE:
			a := vm.StackPop()
			err := vm.StackPushNegate(a)
			if err != nil {
				return err
			}
		case opcode.OP_ADD:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushAdd(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_SUBTRACT:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushSubtract(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_MULTIPLY:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushMultiply(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_DIVIDE:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushDivide(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_MODULO:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushModulo(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_NOT:
			a := vm.StackPop()
			err := vm.StackPushNot(a)
			if err != nil {
				return err
			}
		case opcode.OP_EQ:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushEQ(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_GT:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushGT(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_LT:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushLT(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_GE:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushGE(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_LE:
			b := vm.StackPop()
			a := vm.StackPop()
			err := vm.StackPushLE(a, b)
			if err != nil {
				return err
			}
		case opcode.OP_POP:
			vm.StackPop()
		default:
			return ErrInvalidOpcodeType
		}
	}
	return nil
}
