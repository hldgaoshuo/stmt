package ast

type Stmt interface {
	Node
}

type ExpressionStatement struct {
	Expression Expr
}

type If struct {
	Line       int
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type Print struct {
	Line       int
	Expression Expr
}

type Block struct {
	Line         int
	Declarations []Decl
}

type While struct {
	Line      int
	Condition Expr
	Body      Stmt
}
