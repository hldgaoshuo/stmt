#include <string>
#include <google/protobuf/stubs/common.h>
#include <fmt/core.h>
#include "vm.h"

static bool test_literal_int() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_int()) {
        return false;
    }
    if (result->literal_int() != 1) {
        return false;
    }
    return true;
}

static bool test_literal_float() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_float(1.5);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_float()) {
        return false;
    }
    if (result->literal_float() != 1.5) {
        return false;
    }
    return true;
}

static bool test_negate() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_NEGATE);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(5);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_int()) {
        return false;
    }
    if (result->literal_int() != -5) {
        return false;
    }
    return true;
}

static bool test_add() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_ADD);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(1);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(2);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_int()) {
        return false;
    }
    if (result->literal_int() != 3) {
        return false;
    }
    return true;
}

static bool test_literal_true() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_TRUE);
    chunk->set_code(code);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_literal_false() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_FALSE);
    chunk->set_code(code);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != false) {
        return false;
    }
    return true;
}

static bool test_literal_nil() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_NIL);
    chunk->set_code(code);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_nil()) {
        return false;
    }
    if (!result->literal_nil().empty()) {
        return false;
    }
    return true;
}

static bool test_not() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_NOT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_bool(true);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != false) {
        return false;
    }
    return true;
}

static bool test_eq() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_EQ);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_bool(true);
    const auto c2 = chunk->add_constants();
    c2->set_literal_bool(true);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_gt() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_GT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(2);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_literal_string() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_string("abc");

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_string()) {
        return false;
    }
    if (result->literal_string() != "abc") {
        return false;
    }
    return true;
}

static bool test_add_string() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_ADD);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_string("abc");
    const auto c2 = chunk->add_constants();
    c2->set_literal_string("def");

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_string()) {
        return false;
    }
    if (result->literal_string() != "abcdef") {
        return false;
    }
    return true;
}

static bool test_print() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_PRINT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }
    return true;
}

static bool test_var() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_SET_GLOBAL); code.push_back(0);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_PRINT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(1);

    chunk->set_globals_count(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }
    return true;
}

static bool test_assign() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_SET_GLOBAL); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_SET_GLOBAL); code.push_back(0);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_PRINT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(1);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(2);

    chunk->set_globals_count(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }
    return true;
}

static bool test_block() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_SET_GLOBAL); code.push_back(0);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_PRINT);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_SET_LOCAL); code.push_back(0);
    code.push_back(OP_GET_LOCAL); code.push_back(0);
    code.push_back(OP_PRINT);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_PRINT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(1);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(2);

    chunk->set_globals_count(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }
    return true;
}

static bool test_if() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_TRUE);
    code.push_back(OP_JUMP_FALSE); code.push_back(9);
    code.push_back(OP_POP);
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_PRINT);
    code.push_back(OP_JUMP); code.push_back(10);
    code.push_back(OP_POP);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_PRINT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(10);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(20);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }
    return true;
}

static bool test_if_else() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_FALSE);
    code.push_back(OP_JUMP_FALSE); code.push_back(9);
    code.push_back(OP_POP);
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_PRINT);
    code.push_back(OP_JUMP); code.push_back(13);
    code.push_back(OP_POP);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_PRINT);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(10);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(20);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }
    return true;
}

static bool test_and() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_TRUE);
    code.push_back(OP_JUMP_FALSE); code.push_back(5);
    code.push_back(OP_POP);
    code.push_back(OP_TRUE);
    chunk->set_code(code);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_or() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_TRUE);
    code.push_back(OP_JUMP_FALSE); code.push_back(5);
    code.push_back(OP_JUMP); code.push_back(7);
    code.push_back(OP_POP);
    code.push_back(OP_TRUE);
    chunk->set_code(code);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        return false;
    }

    auto result = vm.stack_pop();
    if (!result->has_literal_bool()) {
        return false;
    }
    if (result->literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_while() {
    const auto chunk = new Object::Chunk();
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_SET_GLOBAL); code.push_back(0);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_LT);
    code.push_back(OP_JUMP_FALSE); code.push_back(24);
    code.push_back(OP_POP);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_PRINT);
    code.push_back(OP_GET_GLOBAL); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(2);
    code.push_back(OP_ADD);
    code.push_back(OP_SET_GLOBAL); code.push_back(0);
    code.push_back(OP_LOOP); code.push_back(4);
    code.push_back(OP_POP);
    chunk->set_code(code);

    const auto c1 = chunk->add_constants();
    c1->set_literal_int(0);
    const auto c2 = chunk->add_constants();
    c2->set_literal_int(5);
    const auto c3 = chunk->add_constants();
    c3->set_literal_int(1);

    chunk->set_globals_count(1);

    VM vm(chunk);
    auto err = vm.interpret();
    if (err != Error::SUCCESS) {
        fmt::print("Error occurred: {}\n", static_cast<int>(err));
        return false;
    }
    return true;
}

int main() {
    GOOGLE_PROTOBUF_VERIFY_VERSION;

    int total = 0;
    int passed = 0;

    struct Test { const char* name; bool (*fn)(); } tests[] = {
        {"test_literal_int", test_literal_int},
        {"test_literal_float", test_literal_float},
        {"test_negate", test_negate},
        {"test_add", test_add},
        {"test_literal_true", test_literal_true},
        {"test_literal_false", test_literal_false},
        {"test_literal_nil", test_literal_nil},
        {"test_not", test_not},
        {"test_eq", test_eq},
        {"test_gt", test_gt},
        {"test_literal_string", test_literal_string},
        {"test_add_string", test_add_string},
        {"test_print", test_print},
        {"test_var", test_var},
        {"test_assign", test_assign},
        {"test_if", test_if},
        {"test_block", test_block},
        {"test_if_else", test_if_else},
        {"test_and", test_and},
        {"test_or", test_or},
        {"test_while", test_while},
    };

    for (auto &t : tests) {
        ++total;
        bool ok = false;
        try { 
            fmt::print("Running test: {}\n", t.name);
            ok = t.fn(); 
        } catch (...) { 
            fmt::print("Exception occurred in test: {}\n", t.name);
            ok = false; 
        }
        if (ok) {
            ++passed;
            fmt::print("[PASS] {}\n", t.name);
        } else {
            fmt::print("[FAIL] {}\n", t.name);
        }
    }

    fmt::print("Summary: {} / {} tests passed\n", passed, total);

    google::protobuf::ShutdownProtobufLibrary();
    return passed == total ? 0 : 1;
}
