package goexp

import (
	"fmt"
	"strconv"
	"strings"
)

var keywords = map[string]TokenType{
	"true":  True,
	"false": False,
	"nil":   Nil,
	"and":   And,
	"or":    Or,
	"not":   Not,
}

type ScannerError struct {
	Message string
	Pos     int
}

func (e *ScannerError) Error() string {
	return fmt.Sprintf("Error @ %d: %s", e.Pos, e.Message)
}

type Scanner struct {
	source  []rune
	start   int
	current int
	length  int
	tokens  []Token
	err     error
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  []rune(source),
		start:   0,
		current: 0,
		length:  len(source),
		tokens:  make([]Token, 0),
	}
}

func (s *Scanner) Scan() ([]Token, error) {
	s.reset()

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	if s.err == nil {
		s.addTok(Token{Type: EOF, Pos: s.length})
	}

	return s.tokens, s.err
}

func (s *Scanner) scanToken() {
	var c = s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen, nil)
		break

	case ')':
		s.addToken(RightParen, nil)
		break

	case '+':
		s.addToken(Add, nil)
		break

	case '-':
		s.addToken(Sub, nil)
		break

	case '*':
		s.addToken(Mul, nil)
		break

	case '/':
		s.addToken(Div, nil)
		break

	case '%':
		s.addToken(Modulo, nil)
		break

	case '.':
		s.addToken(Period, nil)
		break

	case ',':
		s.addToken(Comma, nil)
		break

	case '!':
		if s.match('=') {
			s.addToken(NotEqual, nil)
		} else {
			s.addToken(Not, nil)
		}
		break

	case '=':
		if s.match('=') {
			s.addToken(Equal, nil)
		}
		break

	case '<':
		if s.match('=') {
			s.addToken(LessEqual, nil)
		} else {
			s.addToken(Less, nil)
		}
		break

	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual, nil)
		} else {
			s.addToken(Greater, nil)
		}
		break

	case '&':
		if s.match('&') {
			s.addToken(And, nil)
		}
		break

	case '|':
		if s.match('|') {
			s.addToken(Or, nil)
		}
		break

	case '"', '\'':
		s.readStringLiteral(c)
		break

	case ' ', '\t':
		// ignore whitespaces
		break

	default:
		if isDigit(c) {
			s.readNumber()
		} else if isAlpha(c) {
			s.readIdentifier()
		} else {
			s.error("Unexpected character")
		}
		break
	}

}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.advance()
	return true
}

func (s *Scanner) readNumber() {
	isFloat := false

	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' {
		isFloat = true
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	str := string(s.source[s.start:s.current])

	if isFloat {
		n, err := strconv.ParseFloat(str, 64)
		if err != nil {
			s.error("Invalid integer number: %s (%s)", str, err.Error())
		} else {
			s.addToken(Float, n)
		}
	} else {
		n, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			s.error("Invalid floating point number: %s (%s)", str, err.Error())
		} else {
			s.addToken(Integer, n)
		}
	}
}

func (s *Scanner) readIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	str := string(s.source[s.start:s.current])
	tokenType, isKeyword := keywords[strings.ToLower(str)]
	if !isKeyword {
		tokenType = Identifier
	}
	s.addToken(tokenType, nil)
}

func (s *Scanner) readStringLiteral(term rune) {
	for (s.peek() != term || s.source[s.current-1] == '\\') && !s.isAtEnd() {
		s.advance()
	}

	if s.isAtEnd() {
		s.error("Unterminated string")
	} else {
		s.advance()
		str := string(s.source[s.start+1 : s.current-1])
		s.addToken(String, str)
	}
}

func (s *Scanner) error(message string, args ...interface{}) {
	s.err = &ScannerError{
		Message: fmt.Sprintf(message, args...),
		Pos:     s.current,
	}
}

func (s *Scanner) addToken(t TokenType, literal interface{}) {
	lex := string(s.source[s.start:s.current])
	tok := Token{
		Type:    t,
		Lexeme:  lex,
		Literal: literal,
		Pos:     s.start,
	}
	s.addTok(tok)
}

func (s *Scanner) addTok(tok Token) {
	s.tokens = append(s.tokens, tok)
}

func (s *Scanner) reset() {
	s.current = 0
	s.start = 0
	s.tokens = make([]Token, 0)
	s.err = nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= s.length
}

func (s *Scanner) peekNext() rune {
	var i = s.current + 1
	if i >= s.length {
		return 0
	}
	return s.source[i]
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) advance() rune {
	s.current++
	return s.source[s.current-1]
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}

func isAlpha(char rune) bool {
	return ('a' <= char && char <= 'z') ||
		('A' <= char && char <= 'Z') ||
		'_' == char
}

func isAlphaNumeric(char rune) bool {
	return isAlpha(char) || isDigit(char)
}
