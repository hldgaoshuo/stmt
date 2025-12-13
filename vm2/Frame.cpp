//
// Created by gaoshuo on 2025/12/11.
//

#include "Frame.h"

Frame::Frame(Object::Function *fun, std::size_t bp) {
    function = fun;
    base_pointer = bp;
    ip = 0;
}

std::size_t Frame::code_size() {
    auto code = function->code();
    return code.size();
}

uint8_t Frame::code_next() {
    auto code = function->code();
    return code[ip++];
}
