package ast

type Stmt interface {
	Node
}

type ExpressionStatement struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}
