//
// Created by gaoshuo on 2025/12/1.
//

#include "vm.h"
#include <cmath>
#include <fmt/core.h>

VM::VM(Object::Chunk* chunk) {
    for (uint8_t b : chunk->code()) {
        _code_emit(b);
    }
    for (int i = 0; i < chunk->constants_size(); i++) {
        Object::Object* o = chunk->mutable_constants(i);
        _constant_add(o);
    }
    globals.resize(chunk->globals_count());
}

void VM::_code_emit(uint8_t byte) {
    code.push_back(byte);
}
uint8_t VM::code_next() {
    return code[ip++];
}

void VM::_constant_add(Object::Object* value) {
    retain(value);
    constants.push_back(value);
}
Object::Object* VM::constant_get(uint8_t index) {
    return constants[index];
}

void VM::stack_push(Object::Object* value) {
    retain(value);
    stack.push_back(value);
}
Object::Object* VM::stack_pop() {
    Object::Object* value = stack.back();
    stack.pop_back();
    return value;
}
Object::Object* VM::stack_peek() {
    Object::Object* value = stack.back();
    return value;
}
void VM::stack_set(uint8_t index, Object::Object* value) {
    stack[index] = value;
}
Object::Object* VM::stack_get(uint8_t index) {
    return stack[index];
}

void VM::retain(Object::Object* obj) {
    auto ref_count = obj->ref_count();
    ref_count += 1;
    obj->set_ref_count(ref_count);
}
void VM::release(Object::Object* obj) {
    auto ref_count = obj->ref_count();
    ref_count -= 1;
    obj->set_ref_count(ref_count);
    if (ref_count == 0) {
        delete obj;
    }
}

Error VM::interpret() {
    while (ip < code.size()) {
        switch (uint8_t instruction = code_next()) {
            case OP_CONSTANT: {
                auto constant_index = code_next();
                auto constant = constant_get(constant_index);
                stack_push(constant);
                break;
            }
            case OP_NEGATE: {
                auto value = stack_pop();
                auto result = new Object::Object();
                if (value->has_literal_int()) {
                    result->set_literal_int(-value->literal_int());
                }
                else if (value->has_literal_float()) {
                    result->set_literal_float(-value->literal_float());
                }
                else {
                    fmt::print("Invalid operand for OP_NEGATE\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(value);
                break;
            }
            case OP_ADD: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() + b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() + b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) + b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() + static_cast<double>(b->literal_int()));
                }
                else if (a->has_literal_string() && b->has_literal_string()) {
                    result->set_literal_string(a->literal_string() + b->literal_string());
                }
                else {
                    fmt::print("Invalid operands for OP_ADD\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_SUBTRACT: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() - b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() - b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) - b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() - static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_SUBTRACT\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_MULTIPLY: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() * b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() * b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) * b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() * static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_MULTIPLY\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_DIVIDE: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() / b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() / b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) / b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() / static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_DIVIDE\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_MODULO: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() % b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_float(fmod(a->literal_float(), b->literal_float()));
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(fmod(static_cast<double>(a->literal_int()), b->literal_float()));
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(fmod(a->literal_float(), static_cast<double>(b->literal_int())));
                }
                else {
                    fmt::print("Invalid operands for OP_MODULO\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_TRUE: {
                auto result = new Object::Object();
                result->set_literal_bool(true);
                stack_push(result);
                break;
            }
            case OP_FALSE: {
                auto result = new Object::Object();
                result->set_literal_bool(false);
                stack_push(result);
                break;
            }
            case OP_NIL: {
                auto result = new Object::Object();
                result->set_literal_nil("");
                stack_push(result);
                break;
            }
            case OP_NOT: {
                auto value = stack_pop();
                auto result = new Object::Object();
                if (value->has_literal_bool()) {
                    result->set_literal_bool(!value->literal_bool());
                }
                else {
                    fmt::print("Invalid operand for OP_NOT\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(value);
                break;
            }
            case OP_EQ: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() == b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() == b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) == b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() == static_cast<double>(b->literal_int()));
                }
                else if (a->has_literal_bool() && b->has_literal_bool()) {
                    result->set_literal_bool(a->literal_bool() == b->literal_bool());
                }
                else if (a->has_literal_nil() && b->has_literal_nil()) {
                    result->set_literal_bool(true);
                }
                else {
                    fmt::print("Invalid operands for OP_EQ\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_GT: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() > b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() > b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) > b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() > static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_GT\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_LT: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() < b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() < b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) < b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() < static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_LT\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_GE: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() >= b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() >= b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) >= b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() >= static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_GE\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_LE: {
                auto b = stack_pop();
                auto a = stack_pop();
                auto result = new Object::Object();
                if (a->has_literal_int() && b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() <= b->literal_int());
                }
                else if (a->has_literal_float() && b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() <= b->literal_float());
                }
                else if (a->has_literal_int() && b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) <= b->literal_float());
                }
                else if (a->has_literal_float() && b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() <= static_cast<double>(b->literal_int()));
                }
                else {
                    fmt::print("Invalid operands for OP_LE\n");
                    return Error::ERROR;
                }
                stack_push(result);
                release(a);
                release(b);
                break;
            }
            case OP_POP: {
                auto value = stack_pop();
                release(value);
                break;
            }
            case OP_PRINT: {
                Object::Object* value = stack_pop();
                if (value->has_literal_int()) {
                    fmt::print("{}\n", value->literal_int());
                }
                else if (value->has_literal_float()) {
                    fmt::print("{}\n", value->literal_float());
                }
                else if (value->has_literal_string()) {
                    fmt::print("{}\n", value->literal_string());
                }
                else if (value->has_literal_bool()) {
                    fmt::print("{}\n", value->literal_bool());
                }
                else if (value->has_literal_nil()) {
                    fmt::print("nil\n");
                }
                break;
            }
            case OP_SET_GLOBAL: {
                auto global_index = code_next();
                auto global_value = stack_pop();
                globals[global_index] = global_value;
                break;
            }
            case OP_GET_GLOBAL: {
                auto global_index = code_next();
                auto global_value = globals[global_index];
                stack_push(global_value);
                break;
            }
            case OP_SET_LOCAL: {
                auto local_index = code_next();
                auto local_value = stack_pop();
                stack_set(local_index, local_value);
                break;
            }
            case OP_GET_LOCAL: {
                auto local_index = code_next();
                auto local_value = stack_get(local_index);
                stack_push(local_value);
                break;
            }
            case OP_JUMP_FALSE: {
                auto _ip = code_next();
                if (auto cond = stack_peek(); cond->has_literal_bool()) {
                    if (!cond->literal_bool()) {
                        ip = _ip;
                    }
                } else {
                    fmt::print("Invalid operands for OP_JUMP_FALSE\n");
                    return Error::ERROR;
                }
                break;
            }
            case OP_JUMP: {
                auto _ip = code_next();
                ip = _ip;
                break;
            }
            case OP_LOOP: {
                auto _ip = code_next();
                ip = _ip;
                break;
            }
            default: {
                fmt::print("Unknown opcode {}\n", instruction);
                return Error::ERROR;
            }
        }
    }
    return Error::SUCCESS;
}
