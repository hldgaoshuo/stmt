#include <string>
#include <google/protobuf/stubs/common.h>
#include <fmt/core.h>
#include "vm.h"

static bool test_literal_int() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_int(1);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_int()) {
        return false;
    }
    if (result.literal_int() != 1) {
        return false;
    }
    return true;
}

static bool test_literal_float() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_float(1.5);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_float()) {
        return false;
    }
    if (result.literal_float() != 1.5) {
        return false;
    }
    return true;
}

static bool test_negate() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_NEGATE);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_int(5);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_int()) {
        return false;
    }
    if (result.literal_int() != -5) {
        return false;
    }
    return true;
}

static bool test_add() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_ADD);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_int(1);
    const auto c2 = chunk.add_constants();
    c2->set_literal_int(2);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_int()) {
        return false;
    }
    if (result.literal_int() != 3) {
        return false;
    }
    return true;
}

static bool test_literal_true() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_TRUE);
    chunk.set_code(code);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_bool()) {
        return false;
    }
    if (result.literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_literal_false() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_FALSE);
    chunk.set_code(code);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_bool()) {
        return false;
    }
    if (result.literal_bool() != false) {
        return false;
    }
    return true;
}

static bool test_literal_nil() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_NIL);
    chunk.set_code(code);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_nil()) {
        return false;
    }
    if (result.literal_nil() != "") {
        return false;
    }
    return true;
}

static bool test_not() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_NOT);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_bool(true);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_bool()) {
        return false;
    }
    if (result.literal_bool() != false) {
        return false;
    }
    return true;
}

static bool test_eq() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_EQ);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_bool(true);
    const auto c2 = chunk.add_constants();
    c2->set_literal_bool(true);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_bool()) {
        return false;
    }
    if (result.literal_bool() != true) {
        return false;
    }
    return true;
}

static bool test_gt() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_GT);
    chunk.set_code(code);

    const auto c1 = chunk.add_constants();
    c1->set_literal_int(2);
    const auto c2 = chunk.add_constants();
    c2->set_literal_int(1);

    VM vm(chunk);
    auto [result, err] = vm.run();
    if (err != Error::SUCCESS) {
        return false;
    }
    if (!result.has_literal_bool()) {
        return false;
    }
    if (result.literal_bool() != true) {
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
    };

    for (auto &t : tests) {
        ++total;
        bool ok = false;
        try { ok = t.fn(); } catch (...) { ok = false; }
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
