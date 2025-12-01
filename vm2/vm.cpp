//
// Created by gaoshuo on 2025/12/1.
//

#include "vm.h"
#include <fmt/core.h>

VM::VM(const Object::Chunk& chunk) {
    for (const uint8_t b : chunk.code()) {
        _code_emit(b);
    }
    _code_emit(OP_RETURN);
    for (int i = 0; i < chunk.constants_size(); i++) {
        const Object::Object& o = chunk.constants(i);
        _constant_add(o);
    }
}

void VM::_code_emit(const uint8_t byte) {
    code.push_back(byte);
}
uint8_t VM::code_next() {
    return code[ip++];
}

void VM::_constant_add(const Object::Object& value) {
    constants.push_back(value);
}
Object::Object VM::constant_get(const uint8_t index) {
    return constants[index];
}

void VM::stack_push(const Object::Object& value) {
    stack.push_back(value);
}
Object::Object VM::stack_pop() {
    Object::Object value = stack.back();
    stack.pop_back();
    return value;
}

std::pair<Object::Object, Error> VM::run() {
    for (;;) {
        switch (uint8_t instruction = code_next()) {
            case OP_RETURN: {
                Object::Object result = stack_pop();
                return {result, Error::SUCCESS};
            }
            case OP_CONSTANT: {
                const uint8_t constant_index = code_next();
                Object::Object constant = constant_get(constant_index);
                stack_push(constant);
                break;
            }
            case OP_NEGATE: {
                Object::Object value = stack_pop();
                Object::Object result;
                if (value.has_literal_int()) {
                    result.set_literal_int(-value.literal_int());
                }
                else if (value.has_literal_float()) {
                    result.set_literal_float(-value.literal_float());
                }
                else {
                    fmt::print("Invalid operand for NEGATE\n");
                    return {{}, Error::ERROR};
                }
                stack_push(result);
                break;
            }
            case OP_ADD: {
                Object::Object b = stack_pop();
                Object::Object a = stack_pop();
                Object::Object result;
                if (a.has_literal_int() && b.has_literal_int()) {
                    result.set_literal_int(a.literal_int() + b.literal_int());
                }
                else if (a.has_literal_float() && b.has_literal_float()) {
                    result.set_literal_float(a.literal_float() + b.literal_float());
                }
                else {
                    fmt::print("Invalid operands for ADD\n");
                    return {{}, Error::ERROR};
                }
                stack_push(result);
                break;
            }
            case OP_SUBTRACT: {
                Object::Object b = stack_pop();
                Object::Object a = stack_pop();
                Object::Object result;
                if (a.has_literal_int() && b.has_literal_int()) {
                    result.set_literal_int(a.literal_int() - b.literal_int());
                }
                else if (a.has_literal_float() && b.has_literal_float()) {
                    result.set_literal_float(a.literal_float() - b.literal_float());
                }
                else {
                    fmt::print("Invalid operands for SUBTRACT\n");
                    return {{}, Error::ERROR};
                }
                stack_push(result);
                break;
            }
            case OP_MULTIPLY: {
                Object::Object b = stack_pop();
                Object::Object a = stack_pop();
                Object::Object result;
                if (a.has_literal_int() && b.has_literal_int()) {
                    result.set_literal_int(a.literal_int() * b.literal_int());
                }
                else if (a.has_literal_float() && b.has_literal_float()) {
                    result.set_literal_float(a.literal_float() * b.literal_float());
                }
                else {
                    fmt::print("Invalid operands for MULTIPLY\n");
                    return {{}, Error::ERROR};
                }
                stack_push(result);
                break;
            }
            case OP_DIVIDE: {
                Object::Object b = stack_pop();
                Object::Object a = stack_pop();
                Object::Object result;
                if (a.has_literal_int() && b.has_literal_int()) {
                    result.set_literal_int(a.literal_int() / b.literal_int());
                }
                else if (a.has_literal_float() && b.has_literal_float()) {
                    result.set_literal_float(a.literal_float() / b.literal_float());
                }
                else {
                    fmt::print("Invalid operands for DIVIDE\n");
                    return {{}, Error::ERROR};
                }
                stack_push(result);
                break;
            }
            case OP_TRUE: {
                Object::Object result;
                result.set_literal_bool(true);
                stack_push(result);
                break;
            }
            case OP_FALSE: {
                Object::Object result;
                result.set_literal_bool(false);
                stack_push(result);
                break;
            }
            case OP_NIL: {
                Object::Object result;
                result.set_literal_nil("");
                stack_push(result);
                break;
            }
            default:
                fmt::print("Unknown opcode {}\n", instruction);
                return {{}, Error::ERROR};
        }
    }
}
