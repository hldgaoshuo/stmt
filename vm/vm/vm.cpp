#include <vector>
#include <stack>
#include <fstream>
#include <fmt/core.h>
#include <google/protobuf/stubs/common.h>
#include "object.pb.h"

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
    void _code_emit(uint8_t byte) {
        code.push_back(byte);
    }
    uint8_t code_next() {
        return code[ip++];
    }

	// constants
    std::vector<Object::Object> constants;
    void _constant_add(Object::Object value) {
        constants.push_back(value);
    }
    Object::Object constant_get(uint8_t index) {
        return constants[index];
    }

    // stack
	std::vector<Object::Object> stack;
    void stack_push(Object::Object value) {
        stack.push_back(value);
    }
    Object::Object stack_pop() {
        Object::Object value = stack.back();
        stack.pop_back();
        return value;
    }

    void interpret() {
       for (;;) {
            uint8_t instruction = code_next();
            switch (instruction) {
                case OP_RETURN: {
				// todo 这里的打印是现阶段为了方便观察结果
                    Object::Object result = stack_pop();
                    fmt::print("result {}\n", result.literal_int());
                    return;
                }
                case OP_CONSTANT: {
                    uint8_t constant_index = code_next();
                    Object::Object constant = constant_get(constant_index);
                    stack_push(constant);
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
    GOOGLE_PROTOBUF_VERIFY_VERSION;

    Object::Chunk chunk;
    std::ifstream row("1.bin", std::ios::in | std::ios::binary);
    if (!row) {
        fmt::print("Failed to open 1.bin\n");
        return 1;
    }
    if (!chunk.ParseFromIstream(&row)) {
        fmt::print("Failed to parse Chunk from file\n");
        return 1;
    }

    VM vm;
    for (uint8_t b : chunk.code()) {
        vm._code_emit(b);
	}
	vm._code_emit(OP_RETURN);
    for (int i = 0; i < chunk.constants_size(); i++) {
        Object::Object o = chunk.constants(i);
        vm._constant_add(o);
    }
    vm.interpret();

    google::protobuf::ShutdownProtobufLibrary();
    return 0;
}
