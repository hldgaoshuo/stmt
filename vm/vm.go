package vm

import (
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

func (vm *VM) Run() error {
	frame := vm.FramesBack()
	for frame.Ip < frame.CodeSize() {
		op := frame.Opcode()
		switch op {
		case opcode.OP_CONSTANT:
			globalIndex, err := frame.Operand(op)
			if err != nil {
				return err
			}
			global := vm.Constants[globalIndex]
			vm.StackPush(global)
		default:
			return ErrInvalidOpcodeType
		}
	}
	return nil
}
