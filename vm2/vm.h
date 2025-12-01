//
// Created by gaoshuo on 2025/12/1.
//

#ifndef VM2_VM_H
#define VM2_VM_H

#include <cstdint>
#include <vector>
#include <utility>
#include "object.pb.h"

typedef enum {
    OP_RETURN,
    OP_CONSTANT,
    OP_NEGATE,
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
    OP_TRUE,
    OP_FALSE,
    OP_NIL,
} OpCode;

enum class Error {
    SUCCESS = 0,
    ERROR = 1,
};

class VM {
public:
    explicit VM(const Object::Chunk& chunk);

    // code
    std::vector<uint8_t> code;
    std::size_t ip = 0;
    void _code_emit(uint8_t byte);
    uint8_t code_next();

    // constants
    std::vector<Object::Object> constants;
    void _constant_add(const Object::Object& value);
    Object::Object constant_get(uint8_t index);

    // stack
    std::vector<Object::Object> stack;
    void stack_push(const Object::Object& value);
    Object::Object stack_pop();

    std::pair<Object::Object, Error> run();
};

#endif //VM2_VM_H