package token

const (
	EOF = "EOF"

	// Single-character tokens.
	LEFT_PAREN  = "LEFT_PAREN"
	RIGHT_PAREN = "RIGHT_PAREN"
	LEFT_BRACE  = "LEFT_BRACE"
	RIGHT_BRACE = "RIGHT_BRACE"
	COMMA       = "COMMA"
	DOT         = "DOT"
	MINUS       = "MINUS"
	PLUS        = "PLUS"
	SEMICOLON   = "SEMICOLON"
	SLASH       = "SLASH"
	STAR        = "STAR"
	PERCENTAGE  = "PERCENTAGE"

	// One or two character tokens.
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"

	IDENTIFIER = "IDENTIFIER"

	// Literals.
	STRING_LITERAL = "STRING_LITERAL"
	INT_LITERAL    = "INT_LITERAL"
	FLOAT_LITERAL  = "FLOAT_LITERAL"

	// Keywords.
	AND      = "AND"
	CLASS    = "CLASS"
	ELSE     = "ELSE"
	FALSE    = "FALSE"
	FUN      = "FUN"
	FOR      = "FOR"
	IF       = "IF"
	NIL      = "NIL"
	OR       = "OR"
	PRINT    = "PRINT"
	RETURN   = "RETURN"
	SUPER    = "SUPER"
	THIS     = "THIS"
	TRUE     = "TRUE"
	VAR      = "VAR"
	WHILE    = "WHILE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

type Token struct {
	TokenType string
	Lexeme    string
	Literal   any
	Line      int
}

func New(tokenType string, lexeme string, literal any, line int) *Token {
	return &Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}
