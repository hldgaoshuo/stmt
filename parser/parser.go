package parser

import (
	"errors"
	"log/slog"
	"stmt/ast"
	"stmt/token"
)

var (
	ErrUnexpectedEof           = errors.New("unexpected end of file")
	ErrExpectExpression        = errors.New("expect expression")
	ErrInvalidAssignmentTarget = errors.New("invalid assignment target")
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]ast.Node, error) {
	var decls []ast.Node
	for !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			// todo 可能需要打印一些东西
			continue
		}
		decls = append(decls, decl)
	}
	return decls, nil
}

func (p *Parser) declaration() (ast.Node, error) {
	if p.match(token.CLASS) {
		return p.class()
	}
	if p.match(token.FUN) {
		return p.fun()
	}
	if p.match(token.VAR) {
		return p.var_()
	}
	return p.statement()
}

func (p *Parser) class() (ast.Node, error) {
	kw := p.previous()
	name, err := p.consume(token.IDENTIFIER, "Expect class name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}
	var methods []*ast.Function
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		method, err := p._method()
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	_, err = p.consume(token.RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}
	return &ast.Class{
		Line:    kw.Line,
		Name:    name,
		Methods: methods,
	}, nil
}

func (p *Parser) _method() (*ast.Function, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect method name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after method name.")
	if err != nil {
		return nil, err
	}
	var parameters []*token.Token
	if !p.match(token.RIGHT_PAREN) {
		parameters, err = p._parameters()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before method body.")
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &ast.Function{
		Line:   name.Line,
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) fun() (ast.Node, error) {
	kw := p.previous()
	name, err := p.consume(token.IDENTIFIER, "Expect function name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after function name.")
	if err != nil {
		return nil, err
	}
	var parameters []*token.Token
	if !p.match(token.RIGHT_PAREN) {
		parameters, err = p._parameters()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before function body.")
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &ast.Function{
		Line:   kw.Line,
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) _parameters() ([]*token.Token, error) {
	var parameters []*token.Token
	for {
		name, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, name)
		if !p.match(token.COMMA) {
			break
		}
	}
	return parameters, nil
}

func (p *Parser) var_() (ast.Node, error) {
	kw := p.previous()
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer ast.Node
	if p.match(token.EQUAL) {
		initializer, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &ast.Var{
		Line:        kw.Line,
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (ast.Node, error) {
	if p.match(token.PRINT) {
		return p.print()
	}
	if p.match(token.LEFT_BRACE) {
		return p.block()
	}
	if p.match(token.IF) {
		return p.if_()
	}
	if p.match(token.WHILE) {
		return p.while()
	}
	if p.match(token.FOR) {
		return p.for_()
	}
	if p.match(token.RETURN) {
		return p.return_()
	}
	return p.expressionStatement()
}

func (p *Parser) print() (ast.Node, error) {
	kw := p.previous()
	value, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ast.Print{
		Line:       kw.Line,
		Expression: value,
	}, nil
}

func (p *Parser) block() (ast.Node, error) {
	kw := p.previous()
	var decls []ast.Node
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		decl, err := p.declaration()
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)
	}
	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return &ast.Block{
		Line:         kw.Line,
		Declarations: decls,
	}, nil
}

func (p *Parser) if_() (ast.Node, error) {
	kw := p.previous()
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch ast.Node
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &ast.If{
		Line:       kw.Line,
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) while() (ast.Node, error) {
	kw := p.previous()
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &ast.While{
		Line:      kw.Line,
		Body:      body,
		Condition: condition,
	}, nil
}

func (p *Parser) for_() (ast.Node, error) {
	kw := p.previous()
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}
	var initializer ast.Node
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.var_()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}
	var condition ast.Node
	if !p.match(token.SEMICOLON) {
		condition, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}
	var increment ast.Node
	if !p.match(token.RIGHT_PAREN) {
		increment, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
		body = &ast.Block{
			Line: kw.Line,
			Declarations: []ast.Node{body, &ast.ExpressionStatement{
				Expression: increment,
			}},
		}
	}
	if condition == nil {
		condition = &ast.Literal{
			Value: true,
		}
	}
	body = &ast.While{
		Line:      kw.Line,
		Body:      body,
		Condition: condition,
	}
	if initializer != nil {
		body = &ast.Block{
			Line:         kw.Line,
			Declarations: []ast.Node{initializer, body},
		}
	}
	return body, nil
}

func (p *Parser) return_() (ast.Node, error) {
	kw := p.previous()
	var value ast.Node
	var err error
	if !p.check(token.SEMICOLON) {
		value, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}
	return &ast.Return{
		Line:       kw.Line,
		Expression: value,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Node, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after Expression.")
	if err != nil {
		return nil, err
	}
	return &ast.ExpressionStatement{
		Expression: expr,
	}, nil
}

func (p *Parser) Expression() (ast.Node, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Node, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.Expression()
		if err != nil {
			return nil, err
		}
		switch _expr := expr.(type) {
		case *ast.Variable:
			return &ast.Assign{
				Name:  _expr.Name,
				Value: value,
			}, nil
		case *ast.Get:
			return &ast.Set{
				Name:   _expr.Name,
				Object: _expr.Object,
				Value:  value,
			}, nil
		default:
			slog.Error("Invalid assignment target.", "line", equals.Line, "equals", equals)
			return nil, ErrInvalidAssignmentTarget
		}
	}
	return expr, nil
}

func (p *Parser) or() (ast.Node, error) {
	left, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		left = &ast.Logical{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) and() (ast.Node, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		left = &ast.Logical{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) equality() (ast.Node, error) {
	left, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) comparison() (ast.Node, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) term() (ast.Node, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) factor() (ast.Node, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) unary() (ast.Node, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}
	return p.call()
}

func (p *Parser) call() (ast.Node, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(token.LEFT_PAREN) {
			var arguments []ast.Node
			if !p.check(token.RIGHT_PAREN) {
				arguments, err = p._arguments()
				if err != nil {
					return nil, err
				}
			}
			paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
			if err != nil {
				return nil, err
			}
			expr = &ast.Call{
				Arguments: arguments,
				Callee:    expr,
				Paren:     paren,
			}
		} else if p.match(token.DOT) {
			name, err := p.consume(token.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &ast.Get{
				Name:   name,
				Object: expr,
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) _arguments() ([]ast.Node, error) {
	var arguments []ast.Node
	for {
		argument, err := p.Expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argument)
		if !p.match(token.COMMA) {
			break
		}
	}
	return arguments, nil
}

func (p *Parser) primary() (ast.Node, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{
			Value: false,
		}, nil
	}
	if p.match(token.TRUE) {
		return &ast.Literal{
			Value: true,
		}, nil
	}
	if p.match(token.NIL) {
		return &ast.Literal{
			Value: nil,
		}, nil
	}
	if p.match(token.NUMBER, token.STRING) {
		token_ := p.previous()
		return &ast.Literal{
			Value: token_.Literal,
		}, nil
	}
	if p.match(token.IDENTIFIER) {
		token_ := p.previous()
		return &ast.Variable{
			Name: token_,
		}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.Expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after Expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Grouping{
			Expression: expr,
		}, nil
	}
	return nil, ErrExpectExpression
}

// utils
func (p *Parser) match(tokenTypes ...string) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType string) bool {
	if p.isAtEnd() {
		return false
	}
	token_ := p.peek()
	return token_.TokenType == tokenType
}

func (p *Parser) isAtEnd() bool {
	token_ := p.peek()
	return token_.TokenType == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType string, message string) (*token.Token, error) {
	if p.check(tokenType) {
		token_ := p.advance()
		return token_, nil
	} else {
		token_ := p.peek()
		if token_.TokenType == token.EOF {
			slog.Error("Unexpected end of file.", "line", token_.Line, "message", message)
			return nil, ErrUnexpectedEof
		} else {
			slog.Error("Unexpected token.", "line", token_.Line, "message", message, "token", token_)
			return nil, errors.New(message)
		}
	}
}
