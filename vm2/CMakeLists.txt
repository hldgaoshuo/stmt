cmake_minimum_required(VERSION 4.0)
project(vm2)

set(CMAKE_CXX_STANDARD 26)

find_package(fmt CONFIG REQUIRED)
find_package(Protobuf REQUIRED)

add_executable(vm2
        main.cpp
        vm.cpp
        object.pb.cc
)

target_link_libraries(vm2 PRIVATE
        fmt::fmt
        protobuf::libprotobuf
)
