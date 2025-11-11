package ast

type Stmt interface {
	Node
}

type Expression struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}
