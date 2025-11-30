package compiler

const (
	OBJ_INT uint8 = iota
	OBJ_FLOAT
)

type Object struct {
	ObjectType uint8 `json:"object_type"`
	Literal    any   `json:"literal"`
}
