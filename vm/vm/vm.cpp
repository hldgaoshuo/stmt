#include <vector>
#include <stack>
#include <fmt/core.h>

typedef enum {
    OP_RETURN,
    OP_CONSTANT,
    OP_NEGATE,
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
} OpCode;

class VM {
public:
	// code
    std::vector<uint8_t> code;
    std::size_t ip = 0;
    void _code_write(uint8_t byte) {
        code.push_back(byte);
    }
    uint8_t code_next() {
        return code[ip++];
    }

	// constants
    std::vector<int64_t> constants;
    void _constant_add(int64_t value) {
        constants.push_back(value);
    }
    int64_t constant_get(uint8_t index) {
        return constants[index];
    }

    // stack
	std::vector<int64_t> stack;
    void stack_push(int64_t value) {
        stack.push_back(value);
    }
    int64_t stack_pop() {
        int64_t value = stack.back();
        stack.pop_back();
        return value;
    }

    void interpret() {
       for (;;) {
            uint8_t instruction = code_next();
            switch (instruction) {
                case OP_RETURN: {
					// todo 这里的打印是现阶段为了方便观察结果
                    int64_t result = stack_pop();
                    fmt::print("result {}\n", result);
                    return;
                }
                case OP_CONSTANT: {
                    uint8_t constant_index = code_next();
                    int64_t constant = constant_get(constant_index);
                    stack_push(constant);
                    break;
                }
                case OP_NEGATE: {
                    int64_t value = stack_pop();
                    stack_push(-value);
                    break;
				}
                case OP_ADD: {
                    int64_t b = stack_pop();
                    int64_t a = stack_pop();
                    stack_push(a + b);
                    break;
				}
                case OP_SUBTRACT: {
                    int64_t b = stack_pop();
                    int64_t a = stack_pop();
                    stack_push(a - b);
                    break;
                }
                case OP_MULTIPLY: {
                    int64_t b = stack_pop();
                    int64_t a = stack_pop();
                    stack_push(a * b);
                    break;
				}
                case OP_DIVIDE: {
                    int64_t b = stack_pop();
                    int64_t a = stack_pop();
                    stack_push(a / b);
                    break;
                }
                default:
                    fmt::print("Unknown opcode {}\n", instruction);
                    return;
            }
	   }
    }
};

int main() {
    VM vm;
    vm._code_write(OP_CONSTANT);
    vm._code_write(0);
    vm._code_write(OP_CONSTANT);
    vm._code_write(1);
	vm._code_write(OP_SUBTRACT);
    vm._code_write(OP_RETURN);
    
    vm._constant_add(3);
    vm._constant_add(2);
    
    vm.interpret();

    return 0;
}
