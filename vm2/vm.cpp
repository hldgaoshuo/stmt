//
// Created by gaoshuo on 2025/12/1.
//

#include "vm.h"
#include <cmath>
#include <fmt/core.h>

VM::VM(Object::Chunk* chunk) {
    for (const uint8_t b : chunk->code()) {
        _code_emit(b);
    }
    _code_emit(OP_RETURN);
    for (int i = 0; i < chunk->constants_size(); i++) {
        Object::Object* o = chunk->mutable_constants(i);
        _constant_add(o);
    }
}

void VM::_code_emit(const uint8_t byte) {
    code.push_back(byte);
}
uint8_t VM::code_next() {
    return code[ip++];
}

void VM::_constant_add(Object::Object* value) {
    constants.push_back(value);
}
Object::Object* VM::constant_get(const uint8_t index) const {
    return constants[index];
}

void VM::stack_push(Object::Object* value) {
    stack.push_back(value);
}
Object::Object* VM::stack_pop() {
    Object::Object* value = stack.back();
    stack.pop_back();
    return value;
}

std::pair<Object::Object*, Error> VM::run() {
    for (;;) {
        switch (uint8_t instruction = code_next()) {
            case OP_RETURN: {
                // todo 临时将 result 返回
                auto result = stack_pop();
                return {result, Error::SUCCESS};
            }
            case OP_CONSTANT: {
                const uint8_t constant_index = code_next();
                Object::Object* constant = constant_get(constant_index);
                stack_push(constant);
                break;
            }
            case OP_NEGATE: {
                const Object::Object* value = stack_pop();
                const auto result = new Object::Object();
                if (value->has_literal_int()) {
                    result->set_literal_int(-value->literal_int());
                }
                else if (value->has_literal_float()) {
                    result->set_literal_float(-value->literal_float());
                }
                else {
                    fmt::print("Invalid operand for OP_NEGATE\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete value;
                break;
            }
            case OP_ADD: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() + b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() + b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_ADD\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_SUBTRACT: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() - b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() - b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_SUBTRACT\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_MULTIPLY: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() * b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() * b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_MULTIPLY\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_DIVIDE: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() / b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() / b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_DIVIDE\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_MODULO: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() % b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(fmod(a->literal_float(), b->literal_float()));
                }
                else {
                    fmt::print("Invalid operands for OP_MODULO\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_TRUE: {
                const auto result = new Object::Object();
                result->set_literal_bool(true);
                stack_push(result);
                break;
            }
            case OP_FALSE: {
                const auto result = new Object::Object();
                result->set_literal_bool(false);
                stack_push(result);
                break;
            }
            case OP_NIL: {
                const auto result = new Object::Object();
                result->set_literal_nil("");
                stack_push(result);
                break;
            }
            case OP_NOT: {
                const Object::Object* value = stack_pop();
                const auto result = new Object::Object();
                if (value->has_literal_bool()) {
                    result->set_literal_bool(!value->literal_bool());
                }
                else {
                    fmt::print("Invalid operand for OP_NOT\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete value;
                break;
            }
            case OP_EQ: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() == b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() == b->literal_float());
                }
                else if (a->has_literal_bool() && b->has_literal_bool()) {
                    result->set_literal_bool(a->literal_bool() == b->literal_bool());
                }
                else if (a->has_literal_nil() && b->has_literal_nil()) {
                    result->set_literal_bool(true);

                }
                else {
                    fmt::print("Invalid operands for OP_EQ\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_GT: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() > b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() > b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_GT\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_LT: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() < b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() < b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_LT\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_GE: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() >= b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() >= b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_GE\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            case OP_LE: {
                const Object::Object* b = stack_pop();
                const Object::Object* a = stack_pop();
                const auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() <= b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() <= b->literal_float());
                }
                else {
                    fmt::print("Invalid operands for OP_LE\n");
                    return {nullptr, Error::ERROR};
                }
                stack_push(result);
                delete a;
                delete b;
                break;
            }
            default:
                fmt::print("Unknown opcode {}\n", instruction);
                return {nullptr, Error::ERROR};
        }
    }
}
