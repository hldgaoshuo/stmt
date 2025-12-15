//
// Created by gaoshuo on 2025/12/11.
//

#include "Frame.h"

Frame::Frame(Object::Closure *clo, std::size_t bp) {
    closure = clo;
    base_pointer = bp;
    ip = 0;
}

std::size_t Frame::code_size() {
    auto function = closure->mutable_function();
    auto code = function->code();
    return code.size();
}

uint8_t Frame::code_next() {
    auto function = closure->mutable_function();
    auto code = function->code();
    return code[ip++];
}
