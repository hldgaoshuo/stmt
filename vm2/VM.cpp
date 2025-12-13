//
// Created by gaoshuo on 2025/12/1.
//

#include "VM.h"
#include <cmath>
#include <fmt/core.h>

#include "Frame.h"

VM::VM(Object::Chunk* chunk) {
    // frames
    auto frame = new Frame(chunk->mutable_function(), 0);
    frames.push_back(frame);

    // constants
    for (int i = 0; i < chunk->constants_size(); i++) {
        Object::Object* o = chunk->mutable_constants(i);
        _constant_add(o);
    }

    // globals
    globals.resize(chunk->globals_count());
}

void VM::frame_push(Frame* frame) {
    frames.push_back(frame);
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
Object::Object* VM::stack_peek(uint8_t num) {
    auto size = stack.size();
    Object::Object* value = stack[size-1-num];
    return value;
}
void VM::stack_set(uint8_t index, Object::Object* value) {
    stack[index] = value;
}
Object::Object* VM::stack_get(uint8_t index) {
    return stack[index];
}
uint8_t VM::stack_base_pointer(uint8_t offset) {
    auto size = stack.size();
    auto result = size - 1 - offset;
    return result;
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
    auto frame = frames.back();
    while (frame->ip < frame->code_size()) {
        switch (uint8_t instruction = frame->code_next()) {
            case OP_CONSTANT: {
                auto constant_index = frame->code_next();
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() + b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() + b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) + b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() + static_cast<double>(b->literal_int()));
                }
                else if (a->has_literal_string() and b->has_literal_string()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() - b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() - b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) - b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() * b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() * b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) * b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() / b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_float(a->literal_float() / b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) / b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_int(a->literal_int() % b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_float(fmod(a->literal_float(), b->literal_float()));
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(fmod(static_cast<double>(a->literal_int()), b->literal_float()));
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                    result->set_literal_bool(not value->literal_bool());
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() == b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() == b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) == b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
                    result->set_literal_float(a->literal_float() == static_cast<double>(b->literal_int()));
                }
                else if (a->has_literal_bool() and b->has_literal_bool()) {
                    result->set_literal_bool(a->literal_bool() == b->literal_bool());
                }
                else if (a->has_literal_nil() and b->has_literal_nil()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() > b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() > b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) > b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() < b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() < b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) < b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() >= b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() >= b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) >= b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                if (a->has_literal_int() and b->has_literal_int()) {
                    result->set_literal_bool(a->literal_int() <= b->literal_int());
                }
                else if (a->has_literal_float() and b->has_literal_float()) {
                    result->set_literal_bool(a->literal_float() <= b->literal_float());
                }
                else if (a->has_literal_int() and b->has_literal_float()) {
                    result->set_literal_float(static_cast<double>(a->literal_int()) <= b->literal_float());
                }
                else if (a->has_literal_float() and b->has_literal_int()) {
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
                auto global_index = frame->code_next();
                auto global_value = stack_pop();
                globals[global_index] = global_value;
                break;
            }
            case OP_GET_GLOBAL: {
                auto global_index = frame->code_next();
                auto global_value = globals[global_index];
                stack_push(global_value);
                break;
            }
            case OP_SET_LOCAL: {
                auto local_index = frame->code_next();
                auto local_value = stack_pop();
                stack_set(local_index, local_value);
                break;
            }
            case OP_GET_LOCAL: {
                auto local_index = frame->code_next();
                auto local_value = stack_get(local_index);
                stack_push(local_value);
                break;
            }
            case OP_JUMP_FALSE: {
                auto _ip = frame->code_next();
                if (auto cond = stack_peek(0); cond->has_literal_bool()) {
                    if (not cond->literal_bool()) {
                        frame->ip = _ip;
                    }
                } else {
                    fmt::print("Invalid operands for OP_JUMP_FALSE\n");
                    return Error::ERROR;
                }
                break;
            }
            case OP_JUMP: {
                auto _ip = frame->code_next();
                frame->ip = _ip;
                break;
            }
            case OP_LOOP: {
                auto _ip = frame->code_next();
                frame->ip = _ip;
                break;
            }
            case OP_CALL: {
                auto arg_count = frame->code_next();
                auto obj = stack_peek(arg_count);
                if (not obj->has_literal_function()) {
                    fmt::print("Invalid constant for OP_CALL\n");
                    return Error::ERROR;
                }
                auto base_pointer = stack_base_pointer(arg_count);
                frame = new Frame(obj->mutable_literal_function(), base_pointer);
                frame_push(frame);
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
