//
// Created by gaoshuo on 2025/12/11.
//

#ifndef VM2_FRAME_H
#define VM2_FRAME_H

#include "object.pb.h"


class Frame {
public:
    Frame(Object::Function* fun, std::size_t bp);

    Object::Function* function;
    std::size_t base_pointer;
    std::size_t ip;

    std::size_t code_size();
    uint8_t code_next();
};


#endif //VM2_FRAME_H