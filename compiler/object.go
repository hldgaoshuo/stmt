package compiler

import (
	"encoding/json"
	"fmt"
)

const (
	INT   = "INT"
	FLOAT = "FLOAT"
)

type Object struct {
	ObjectType string `json:"object_type"`
	Literal    any    `json:"literal"`
}

func objectsEmit(objects []*Object) {
	// 将对象序列化为 JSON 格式
	data, err := json.Marshal(objects)
	if err != nil {
		fmt.Println("Error serializing objects to JSON:", err)
		return
	}

	// 打印 JSON 数据
	fmt.Println(string(data))
}
