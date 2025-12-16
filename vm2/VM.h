//
// Created by gaoshuo on 2025/12/1.
//

#ifndef VM2_VM_H
#define VM2_VM_H

#include <vector>
#include <utility>
#include "object.pb.h"
#include "Frame.h"

typedef enum {
    OP_CONSTANT,
    OP_NEGATE,
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
    OP_MODULO,
    OP_TRUE,
    OP_FALSE,
    OP_NIL,
    OP_NOT,
    OP_EQ,
    OP_GT,
    OP_LT,
    OP_GE,
    OP_LE,
    OP_POP,
    OP_PRINT,
    OP_SET_GLOBAL,
    OP_GET_GLOBAL,
    OP_SET_LOCAL,
    OP_GET_LOCAL,
    OP_JUMP_FALSE,
    OP_JUMP,
    OP_LOOP,
    OP_CALL,
    OP_RETURN,
    OP_CLOSURE,
} OpCode;

enum class Error {
    SUCCESS = 0,
    ERROR = 1,
};

class VM {
public:
    VM(Object::Chunk* chunk);

    // frames
    std::vector<Frame*> frames;
    void frame_show();
    void frame_push(Frame* frame);
    Frame* frame_pop();

    // constants
    std::vector<Object::Object*> constants;
    void _constant_add(Object::Object* value);
    Object::Object* constant_get(uint8_t index);

    // stack
    std::vector<Object::Object*> stack;
    void stack_show();
    void stack_push(Object::Object* value);
    Object::Object* stack_pop();
    Object::Object* stack_peek(uint8_t num);
    void stack_set(uint8_t index, Object::Object* value);
    Object::Object* stack_get(uint8_t index);
    uint8_t stack_base_pointer(uint8_t offset);
    static uint8_t stack_local_index(Frame* frame);
    void stack_resize(std::size_t offset);

    // globals
    std::vector<Object::Object*> globals;

    // interpret
    Error interpret();

    // gc
    static void release(Object::Object* obj);
    static void retain(Object::Object* obj);
};

#endif //VM2_VM_H