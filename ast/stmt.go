package ast

type Print struct {
	Line       int
	Expression Expr
}

func (p *Print) node()    {}
func (p *Print) stmt()    {}
func (p *Print) Pos() int { return p.Line }

type Block struct {
	Line         int
	Declarations []Stmt
}

func (b *Block) node()    {}
func (b *Block) stmt()    {}
func (b *Block) Pos() int { return b.Line }

type If struct {
	Line       int
	Condition  Expr
	ThenBranch *Block
	ElseBranch *Block // can be nil
}

func (i *If) node()    {}
func (i *If) stmt()    {}
func (i *If) Pos() int { return i.Line }

type While struct {
	Line      int
	Condition Expr
	Body      *Block
}

func (w *While) node()    {}
func (w *While) stmt()    {}
func (w *While) Pos() int { return w.Line }

type Return struct {
	Line       int
	Expression Expr // can be nil for empty return
}

func (r *Return) node()    {}
func (r *Return) stmt()    {}
func (r *Return) Pos() int { return r.Line }

type Break struct {
	Line int
}

func (b *Break) node()    {}
func (b *Break) stmt()    {}
func (b *Break) Pos() int { return b.Line }

type Continue struct {
	Line int
}

func (c *Continue) node()    {}
func (c *Continue) stmt()    {}
func (c *Continue) Pos() int { return c.Line }

type ExpressionStatement struct {
	Line       int
	Expression Expr
}

func (e *ExpressionStatement) node()    {}
func (e *ExpressionStatement) stmt()    {}
func (e *ExpressionStatement) Pos() int { return e.Line }
