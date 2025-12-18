package value

type Value interface {
}

type Int struct {
	Literal int64
}

type Float struct {
	Literal float64
}

type Bool struct {
	Literal bool
}

type Nil struct {
}

type String struct {
	Literal string
}

type Function struct {
	Code        []byte
	NumParams   uint64
	NumUpvalues uint64
}
