#include <string>
#include <google/protobuf/stubs/common.h>
#include <fmt/core.h>
#include "vm.h"

static bool test_literal_int() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk.set_code(code);

    Object::Object* c1 = chunk.add_constants();
    c1->set_literal_int(1);

    VM vm(chunk);
    const Object::Object result = vm.run();
    return result.has_literal_int() && result.literal_int() == 1;
}

static bool test_literal_float() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    chunk.set_code(code);

    Object::Object* c1 = chunk.add_constants();
    c1->set_literal_float(1.5);

    VM vm(chunk);
    const Object::Object result = vm.run();
    return result.has_literal_float() && result.literal_float() == 1.5;
}

static bool test_negate() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_NEGATE);
    chunk.set_code(code);

    Object::Object* c1 = chunk.add_constants();
    c1->set_literal_int(5);

    VM vm(chunk);
    const Object::Object result = vm.run();
    return result.has_literal_int() && result.literal_int() == -5;
}

static bool test_add() {
    Object::Chunk chunk;
    std::string code;
    code.push_back(OP_CONSTANT); code.push_back(0);
    code.push_back(OP_CONSTANT); code.push_back(1);
    code.push_back(OP_ADD);
    chunk.set_code(code);

    Object::Object* c1 = chunk.add_constants();
    c1->set_literal_int(1);
    Object::Object* c2 = chunk.add_constants();
    c2->set_literal_int(2);

    VM vm(chunk);
    const Object::Object result = vm.run();
    return result.has_literal_int() && result.literal_int() == 3;
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
    return (passed == total) ? 0 : 1;
}
