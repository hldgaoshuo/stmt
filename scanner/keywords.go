package scanner

import "stmt/token"

var keywords = map[string]string{
	"and":      token.AND,
	"class":    token.CLASS,
	"else":     token.ELSE,
	"false":    token.FALSE,
	"for":      token.FOR,
	"fun":      token.FUN,
	"if":       token.IF,
	"nil":      token.NIL,
	"or":       token.OR,
	"print":    token.PRINT,
	"return":   token.RETURN,
	"super":    token.SUPER,
	"this":     token.THIS,
	"true":     token.TRUE,
	"var":      token.VAR,
	"while":    token.WHILE,
	"break":    token.BREAK,
	"continue": token.CONTINUE,
}
