package vm

import (
	"io"
	"math"
	"os"
	"stmt/opcode"
	"stmt/value"
)

var Output io.Writer = os.Stdout

type VM struct {
	Stack     []value.Value
	Globals   []value.Value
	Frames    []*Frame
	Constants []value.Value
}

func New(code []uint8, constants []value.Value, globalCount int) *VM {
	mainFunction := &value.Function{
		Code:        code,
		NumParams:   0,
		NumUpvalues: 0,
	}
	mainClosure := &value.Closure{
		Function: mainFunction,
	}
	mainFrame := NewFrame(mainClosure, 0)
	return &VM{
		Stack:     []value.Value{},
		Globals:   make([]value.Value, globalCount),
		Frames:    []*Frame{mainFrame},
		Constants: constants,
	}
}

func (vm *VM) Run() error {
	frame := vm.FramesTop()
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
			r := value.NewBool(true)
			vm.StackPush(r)
		case opcode.OP_FALSE:
			r := value.NewBool(false)
			vm.StackPush(r)
		case opcode.OP_NIL:
			r := value.NewNil()
			vm.StackPush(r)
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
		case opcode.OP_PRINT:
			a := vm.StackPop()
			err := a.Print(Output)
			if err != nil {
				return err
			}
		case opcode.OP_SET_GLOBAL:
			globalIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			globalValue := vm.StackPop()
			vm.Globals[globalIndex] = globalValue
		case opcode.OP_GET_GLOBAL:
			globalIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			globalValue := vm.Globals[globalIndex]
			vm.StackPush(globalValue)
		case opcode.OP_SET_LOCAL:
			localIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			stackIndex := frame.BasePointer + localIndex
			value_ := vm.StackPop()
			vm.StackSet(stackIndex, value_)
		case opcode.OP_GET_LOCAL:
			localIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			stackIndex := frame.BasePointer + localIndex
			value_ := vm.StackGet(stackIndex)
			vm.StackPush(value_)
		case opcode.OP_JUMP_FALSE:
			offset, err := frame.Operand(op)
			if err != nil {
				return err
			}
			cond := vm.StackPeek(0)
			switch _cond := cond.(type) {
			case *value.Bool:
				if !_cond.Literal {
					frame.MoveIp(offset)
				}
			default:
				return ErrInvalidCondType
			}
		case opcode.OP_JUMP:
			offset, err := frame.Operand(op)
			if err != nil {
				return err
			}
			frame.MoveIp(offset)
		case opcode.OP_LOOP:
			offset, err := frame.Operand(op)
			if err != nil {
				return err
			}
			frame.MoveIp(-offset)
		case opcode.OP_CALL:
			argCount, err := frame.Operand(op)
			if err != nil {
				return err
			}
			closure := vm.StackPeek(argCount)
			_closure, ok := closure.(*value.Closure)
			if !ok {
				return ErrInvalidCallType
			}
			basePointer := vm.StackLen() - argCount
			frame = NewFrame(_closure, basePointer)
			vm.FramesPush(frame)
		case opcode.OP_RETURN:
			result := vm.StackPop()
			vm.StackResize(frame.BasePointer)
			vm.StackPush(result)
			frame = vm.FramesPop()
		case opcode.OP_CLOSURE, opcode.OP_CLOSURE_2, opcode.OP_CLOSURE_4, opcode.OP_CLOSURE_8:
			functionIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			function := vm.Constants[functionIndex]
			_function, ok := function.(*value.Function)
			if !ok {
				return ErrInvalidClosureType
			}
			closure := value.NewClosure(_function)
			for i := uint64(0); i < _function.NumUpvalues; i++ {
				isLocal := frame.CodeNext()
				index := frame.CodeNext()
				if isLocal == 1 {
					localIndex := index
					stackIndex := frame.BasePointer + uint64(localIndex)
					upvalue := vm.StackGet(stackIndex)
					closure.Upvalues[i] = upvalue
				} else {
					upvalueIndex := index
					upvalue := frame.Closure.Upvalues[upvalueIndex]
					closure.Upvalues[i] = upvalue
				}
			}
			vm.StackPush(closure)
		case opcode.OP_SET_UPVALUE:
			upvalueIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			value_ := vm.StackPop()
			upvalue := frame.Closure.Upvalues[upvalueIndex]
			upvalue.SetLiteral(value_.GetLiteral())
		case opcode.OP_GET_UPVALUE:
			upvalueIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			value_ := frame.Closure.Upvalues[upvalueIndex]
			vm.StackPush(value_)
		default:
			return ErrInvalidOpcodeType
		}
	}
	return nil
}

func (vm *VM) FramesTop() *Frame {
	return vm.Frames[len(vm.Frames)-1]
}

func (vm *VM) FramesPush(frame *Frame) {
	vm.Frames = append(vm.Frames, frame)
}

func (vm *VM) FramesPop() *Frame {
	vm.Frames = vm.Frames[:len(vm.Frames)-1]
	return vm.FramesTop()
}

func (vm *VM) StackPush(value_ value.Value) {
	vm.Stack = append(vm.Stack, value_)
}

func (vm *VM) StackPop() value.Value {
	offset := len(vm.Stack) - 1
	value_ := vm.Stack[offset]
	vm.Stack = vm.Stack[:offset]
	return value_
}

func (vm *VM) StackPeek(num uint64) value.Value {
	return vm.Stack[uint64(len(vm.Stack))-1-num]
}

func (vm *VM) StackSet(index uint64, value_ value.Value) {
	if vm.StackLen() == index {
		vm.StackPush(value_)
	} else {
		vm.Stack[index] = value_
	}
}

func (vm *VM) StackGet(index uint64) value.Value {
	return vm.Stack[index]
}

func (vm *VM) StackLen() uint64 {
	return uint64(len(vm.Stack))
}

func (vm *VM) StackResize(basePointer uint64) {
	vm.Stack = vm.Stack[:basePointer]
}

func (vm *VM) StackPushConstant(global value.Value) error {
	switch global.(type) {
	case *value.Int, *value.Float, *value.String:
		vm.StackPush(global)
		return nil
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushNegate(a value.Value) error {
	switch _a := a.(type) {
	case *value.Int:
		r := value.NewInt(-_a.Literal)
		vm.StackPush(r)
		return nil
	case *value.Float:
		r := value.NewFloat(-_a.Literal)
		vm.StackPush(r)
		return nil
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushAdd(a value.Value, b value.Value) error {
	switch _a := a.(type) {
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewInt(_a.Literal + _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewFloat(float64(_a.Literal) + _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewFloat(_a.Literal + float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewFloat(_a.Literal + _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.String:
		switch _b := b.(type) {
		case *value.String:
			r := value.NewString(_a.Literal + _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushSubtract(a value.Value, b value.Value) error {
	switch _a := a.(type) {
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewInt(_a.Literal - _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewFloat(float64(_a.Literal) - _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewFloat(_a.Literal - float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewFloat(_a.Literal - _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushMultiply(a value.Value, b value.Value) error {
	switch _a := a.(type) {
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewInt(_a.Literal * _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewFloat(float64(_a.Literal) * _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewFloat(_a.Literal * float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewFloat(_a.Literal * _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushDivide(a value.Value, b value.Value) error {
	switch _a := a.(type) {
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			if _b.Literal == 0 {
				return ErrZeroInDivide
			}
			r := value.NewInt(_a.Literal / _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			if _b.Literal == 0 {
				return ErrZeroInDivide
			}
			r := value.NewFloat(float64(_a.Literal) / _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			if _b.Literal == 0 {
				return ErrZeroInDivide
			}
			r := value.NewFloat(_a.Literal / float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			if _b.Literal == 0 {
				return ErrZeroInDivide
			}
			r := value.NewFloat(_a.Literal / _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushModulo(a value.Value, b value.Value) error {
	switch _a := a.(type) {
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			if _b.Literal == 0 {
				return ErrZeroInModulo
			}
			r := value.NewInt(_a.Literal % _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			if _b.Literal == 0 {
				return ErrZeroInModulo
			}
			r := value.NewFloat(math.Mod(float64(_a.Literal), _b.Literal))
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			if _b.Literal == 0 {
				return ErrZeroInModulo
			}
			r := value.NewFloat(math.Mod(_a.Literal, float64(_b.Literal)))
			vm.StackPush(r)
			return nil
		case *value.Float:
			if _b.Literal == 0 {
				return ErrZeroInModulo
			}
			r := value.NewFloat(math.Mod(_a.Literal, _b.Literal))
			vm.StackPush(r)
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
	case *value.Bool:
		r := value.NewBool(!_a.Literal)
		vm.StackPush(r)
		return nil
	default:
		return ErrInvalidOperandType
	}
}

func (vm *VM) StackPushEQ(a any, b any) error {
	switch _a := a.(type) {
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal == _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(float64(_a.Literal) == _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal == float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(_a.Literal == _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Bool:
		switch _b := b.(type) {
		case *value.Bool:
			r := value.NewBool(_a.Literal == _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Nil:
		switch b.(type) {
		case *value.Nil:
			r := value.NewBool(true)
			vm.StackPush(r)
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
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal > _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(float64(_a.Literal) > _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal > float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(_a.Literal > _b.Literal)
			vm.StackPush(r)
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
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal < _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(float64(_a.Literal) < _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal < float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(_a.Literal < _b.Literal)
			vm.StackPush(r)
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
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal >= _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(float64(_a.Literal) >= _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal >= float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(_a.Literal >= _b.Literal)
			vm.StackPush(r)
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
	case *value.Int:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal <= _b.Literal)
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(float64(_a.Literal) <= _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	case *value.Float:
		switch _b := b.(type) {
		case *value.Int:
			r := value.NewBool(_a.Literal <= float64(_b.Literal))
			vm.StackPush(r)
			return nil
		case *value.Float:
			r := value.NewBool(_a.Literal <= _b.Literal)
			vm.StackPush(r)
			return nil
		default:
			return ErrInvalidOperandType
		}
	default:
		return ErrInvalidOperandType
	}
}
