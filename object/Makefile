# JB IDE，需要设置 make 所在路径

.PHONY: all go cpp clean

all: go cpp

go:
	protoc -I. --go_out=. object.proto

cpp:
	protoc -I. --cpp_out=../vm2 object.proto

clean:
	rm -f *.pb.go
	rm -f ../vm2/*.pb.cc ../vm2/*.pb.h

help:
	@echo "可用命令:"
	@echo "  make all   - 生成所有语言代码"
	@echo "  make go    - 生成 Go 代码"
	@echo "  make cpp   - 生成 C++ 代码"
	@echo "  make clean - 清理生成文件"
