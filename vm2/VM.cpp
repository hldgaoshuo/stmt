//
// Created by gaoshuo on 2025/12/1.
//

#include "VM.h"
#include <cmath>
#include <fmt/core.h>

#include "Frame.h"

VM::VM(Object::Chunk* chunk) {
    // frames
    auto frame = new Frame(chunk->mutable_closure(), 0);
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
Frame* VM::frame_pop() {
    frames.pop_back();
    auto frame = frames.back();
    return frame;
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
    auto value = stack.back();
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
    auto result = size - offset;
    return result;
}
uint8_t VM::stack_local_index(Frame* frame) {
    return frame->code_next() + frame->base_pointer;
}
void VM::stack_resize(std::size_t offset) {
    stack.resize(offset);
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
                    fmt::print("Invalid operands for OP_ADD, a is ({}), b is ({})\n", a->DebugString(), b->DebugString());
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
                else if (value->has_literal_function()) {
                    fmt::print("function {}\n", value->DebugString());
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
                auto local_index = stack_local_index(frame);
                auto local_value = stack_pop();
                stack_set(local_index, local_value);
                break;
            }
            case OP_GET_LOCAL: {
                auto local_index = stack_local_index(frame);
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
                if (not obj->has_literal_closure()) {
                    fmt::print("Invalid constant for OP_CALL\n");
                    return Error::ERROR;
                }
                auto base_pointer = stack_base_pointer(arg_count);
                frame = new Frame(obj->mutable_literal_closure(), base_pointer);
                frame_push(frame);
                break;
            }
            case OP_RETURN: {
                auto result = stack_pop();
                // fmt::print("OP_RETURN result: {}\n", result->DebugString());
                // fmt::print("OP_RETURN current frame base_pointer: {}\n", frame->base_pointer);
                stack_resize(frame->base_pointer);
                stack_push(result);
                frame = frame_pop();
                break;
            }
            case OP_CLOSURE: {
                auto function_index = frame->code_next();
                auto function_obj = constant_get(function_index);
                if (not function_obj->has_literal_function()) {
                    fmt::print("Invalid constant for OP_CLOSURE\n");
                    return Error::ERROR;
                }
                auto function = function_obj->mutable_literal_function();
                auto obj = new Object::Object();
                auto closure = new Object::Closure();
                closure->set_allocated_function(function);

                for (uint64_t i = 0; i < function->num_upvalues(); i++) {
                    auto is_local = frame->code_next();
                    auto index = frame->code_next();
                    auto upvalue = new Object::Object();
                    if (is_local == 1) {
                        auto local_index = frame->base_pointer + index;
                        auto local_value = stack_get(local_index);
                        upvalue->CopyFrom(*local_value);
                    } else {
                        auto parent_upvalue = frame->closure->upvalues(index);
                        upvalue->CopyFrom(parent_upvalue);
                    }
                    closure->add_upvalues()->CopyFrom(*upvalue);
                }

                obj->set_allocated_literal_closure(closure);
                stack_push(obj);
                break;
            }
            case OP_GET_UPVALUE: {
                auto upvalue_index = frame->code_next();
                auto upvalue = frame->closure->mutable_upvalues(upvalue_index);
                stack_push(upvalue);
                break;
            }
            case OP_SET_UPVALUE: {
                auto upvalue_index = frame->code_next();
                auto value = stack_pop();
                auto upvalue = frame->closure->mutable_upvalues(upvalue_index);
                upvalue->CopyFrom(*value);
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

void VM::stack_show() {
    fmt::print("== Stack Debug Info ==\n");
    for (std::size_t i = 0; i < stack.size(); ++i) {
        Object::Object* obj = stack[i];
        fmt::print("[{}] ", i);
        if (obj->has_literal_int()) {
            fmt::print("int: {}\n", obj->literal_int());
        }
        else if (obj->has_literal_float()) {
            fmt::print("float: {}\n", obj->literal_float());
        }
        else if (obj->has_literal_string()) {
            fmt::print("string: {}\n", obj->literal_string());
        }
        else if (obj->has_literal_bool()) {
            fmt::print("bool: {}\n", obj->literal_bool());
        }
        else if (obj->has_literal_nil()) {
            fmt::print("nil\n");
        }
        else if (obj->has_literal_function()) {
            fmt::print("function {}\n", obj->DebugString());
        }
        else {
            fmt::print("unknown type\n");
        }
    }
    fmt::print("======================\n");
}

void VM::frame_show() {
    fmt::print("== Frames Debug Info ==\n");
    for (std::size_t i = 0; i < frames.size(); ++i) {
        Frame* frame = frames[i];
        fmt::print("[Frame {}] IP: {}, Base Pointer: {}, Function bytes: ",
                   i,
                   frame->ip,
                   frame->base_pointer);
        const std::string& code = frame->closure->mutable_function()->code();
        fmt::print("[");
        for (std::size_t j = 0; j < code.size(); ++j) {
            if (j > 0) fmt::print(", ");
            fmt::print("{}", static_cast<uint8_t>(code[j]));
        }
        fmt::print("]\n");
    }
    fmt::print("======================\n");
}
