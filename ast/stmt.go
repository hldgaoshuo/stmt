package ast

type ExpressionStatement struct {
	Expression Node
}

type If struct {
	Line       int
	Condition  Node
	ThenBranch Node
	ElseBranch Node
}

type Print struct {
	Line       int
	Expression Node
}

type Block struct {
	Line         int
	Declarations []Node
}

type While struct {
	Line      int
	Condition Node
	Body      Node
}

type Return struct {
	Line       int
	Expression Node
}
