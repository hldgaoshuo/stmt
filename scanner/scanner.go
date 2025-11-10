package scanner

import (
	"log/slog"
	"stmt/token"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []*token.Token
	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []*token.Token{},
		line:    1,
		start:   0,
		current: 0,
	}
}

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.IsAtEnd() {
		s.start = s.current
		s.ScanToken()
	}
	tokenEof := token.New(token.EOF, "", nil, s.line)
	s.tokens = append(s.tokens, tokenEof)
	return s.tokens
}

func (s *Scanner) ScanToken() {
	char := s.Advance()
	switch char {
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '(':
		s.AddToken(token.LEFT_PAREN, nil)
	case ')':
		s.AddToken(token.RIGHT_PAREN, nil)
	case '{':
		s.AddToken(token.LEFT_BRACE, nil)
	case '}':
		s.AddToken(token.RIGHT_BRACE, nil)
	case ',':
		s.AddToken(token.COMMA, nil)
	case '.':
		s.AddToken(token.DOT, nil)
	case ';':
		s.AddToken(token.SEMICOLON, nil)
	case '+':
		s.AddToken(token.PLUS, nil)
	case '-':
		s.AddToken(token.MINUS, nil)
	case '*':
		s.AddToken(token.STAR, nil)
	case '/':
		if s.Match('/') {
			for s.Peek() != '\n' && !s.IsAtEnd() {
				s.Advance()
			}
		} else {
			s.AddToken(token.SLASH, nil)
		}
	case '!':
		if s.Match('=') {
			s.AddToken(token.BANG_EQUAL, nil)
		} else {
			s.AddToken(token.BANG, nil)
		}
	case '=':
		if s.Match('=') {
			s.AddToken(token.EQUAL_EQUAL, nil)
		} else {
			s.AddToken(token.EQUAL, nil)
		}
	case '<':
		if s.Match('=') {
			s.AddToken(token.LESS_EQUAL, nil)
		} else {
			s.AddToken(token.LESS, nil)
		}
	case '>':
		if s.Match('=') {
			s.AddToken(token.GREATER_EQUAL, nil)
		} else {
			s.AddToken(token.GREATER, nil)
		}
	case '"':
		s.String()
	default:
		if s.IsDigit(char) {
			s.Number()
		} else if s.IsAlpha(char) {
			s.Identifier()
		} else {
			slog.Error("Unexpected character.", "line", s.line)
		}
	}
}

func (s *Scanner) Advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) AddToken(tokenType string, literal any) {
	lexeme := s.source[s.start:s.current]
	token_ := token.New(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, token_)
}

func (s *Scanner) Match(expected byte) bool {
	if s.IsAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) Peek() byte {
	if s.IsAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) PeekNext() byte {
	if s.IsNextAtEnd() {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) IsAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) IsNextAtEnd() bool {
	return s.current+1 >= len(s.source)
}

func (s *Scanner) String() {
	for s.Peek() != '"' && !s.IsAtEnd() {
		if s.Peek() == '\n' {
			s.line++
		}
		s.Advance()
	}
	if s.IsAtEnd() {
		slog.Error("Unterminated string.", "line", s.line)
		return
	}

	// The closing ".
	s.Advance()

	// Trim the surrounding quotes.
	literal := s.source[s.start+1 : s.current-1]

	s.AddToken(token.STRING, literal)
}

func (s *Scanner) IsDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) Number() {
	for s.IsDigit(s.Peek()) {
		s.Advance()
	}
	if s.Peek() == '.' && s.IsDigit(s.PeekNext()) {
		s.Advance()
		for s.IsDigit(s.Peek()) {
			s.Advance()
		}
	}
	literalStr := s.source[s.start:s.current]
	literalFloat, err := strconv.ParseFloat(literalStr, 64)
	if err != nil {
		slog.Error("Invalid number format.", "value", literalStr, "line", s.line)
		return
	}
	s.AddToken(token.NUMBER, literalFloat)
}

func (s *Scanner) IsAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func (s *Scanner) IsKeyword() (bool, string) {
	lexeme := s.source[s.start:s.current]
	tokenType, ok := keywords[lexeme]
	return ok, tokenType
}

func (s *Scanner) Identifier() {
	for s.IsAlpha(s.Peek()) || s.IsDigit(s.Peek()) {
		s.Advance()
	}
	isKeyword, tokenType := s.IsKeyword()
	if isKeyword {
		s.AddToken(tokenType, nil)
	} else {
		s.AddToken(token.IDENTIFIER, nil)
	}
}
