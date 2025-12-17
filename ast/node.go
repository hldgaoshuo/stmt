package ast

// Node is the base interface for all AST nodes
type Node interface {
	node()
	Pos() int // returns line number for error reporting
}

// Expr represents an expression node
type Expr interface {
	Node
	expr()
}

// Stmt represents a statement node
type Stmt interface {
	Node
	stmt()
}
