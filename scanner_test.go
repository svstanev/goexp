package goexp

import (
	"testing"

	"github.com/go-test/deep"
)

func TestIsDigit(t *testing.T) {
	digits := []rune("0123456789")
	for _, ch := range digits {
		if isDigit(ch) != true {
			t.Errorf("Expected '%s'(%v) to be a digit", string(ch), ch)
		}
	}

	nonDigits := []rune("abcx.,/?")
	for _, ch := range nonDigits {
		if isDigit(ch) != false {
			t.Errorf("Expected '%s'(%v) NOT to be a digit", string(ch), ch)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	for _, ch := range []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_") {
		if isAlpha(ch) != true {
			t.Errorf("Expected '%s'(%v) to be an letter of underscore", string(ch), ch)
		}
	}

	for _, ch := range []rune("1234567.,[];`") {
		if isAlpha(ch) != false {
			t.Errorf("Expected '%s'(%v) to NOT be an letter of underscore", string(ch), ch)
		}
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		src            string
		expectedTokens []Token
		expectedError  error
	}{
		{"", []Token{Token{Type: EOF, Pos: 0}}, nil},
		{"1", []Token{Token{Integer, string("1"), int64(1), 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"123", []Token{Token{Integer, "123", int64(123), 0}, Token{Type: EOF, Pos: 3}}, nil},
		{"1.23", []Token{Token{Float, "1.23", float64(1.23), 0}, Token{Type: EOF, Pos: 4}}, nil},
		{"''", []Token{Token{String, "''", "", 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"'abc'", []Token{Token{String, "'abc'", "abc", 0}, Token{Type: EOF, Pos: 5}}, nil},
		{"'ab\\'c'", []Token{Token{String, "'ab\\'c'", "ab\\'c", 0}, Token{Type: EOF, Pos: 7}}, nil},
		{"\"\"", []Token{Token{String, "\"\"", "", 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"\"abc\"", []Token{Token{String, "\"abc\"", "abc", 0}, Token{Type: EOF, Pos: 5}}, nil},
		{"\"ab\\\"c\"", []Token{Token{String, "\"ab\\\"c\"", "ab\\\"c", 0}, Token{Type: EOF, Pos: 7}}, nil},
		{"'ab", []Token{}, &ScannerError{Message: "Unterminated string", Pos: 3}},
		{"foo", []Token{Token{Identifier, "foo", nil, 0}, Token{Type: EOF, Pos: 3}}, nil},
		{"<", []Token{Token{Less, "<", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"<=", []Token{Token{LessEqual, "<=", nil, 0}, Token{Type: EOF, Pos: 2}}, nil},
		{">", []Token{Token{Greater, ">", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{">=", []Token{Token{GreaterEqual, ">=", nil, 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"==", []Token{Token{Equal, "==", nil, 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"!=", []Token{Token{NotEqual, "!=", nil, 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"+", []Token{Token{Add, "+", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"+++", []Token{Token{Add, "+", nil, 0}, Token{Add, "+", nil, 1}, Token{Add, "+", nil, 2}, Token{Type: EOF, Pos: 3}}, nil},
		{"-", []Token{Token{Sub, "-", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"*", []Token{Token{Mul, "*", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"/", []Token{Token{Div, "/", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"%", []Token{Token{Modulo, "%", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"!", []Token{Token{Not, "!", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{"&&", []Token{Token{And, "&&", nil, 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"||", []Token{Token{Or, "||", nil, 0}, Token{Type: EOF, Pos: 2}}, nil},
		{"true", []Token{Token{True, "true", nil, 0}, Token{Type: EOF, Pos: 4}}, nil},
		{"false", []Token{Token{False, "false", nil, 0}, Token{Type: EOF, Pos: 5}}, nil},
		{"nil", []Token{Token{Nil, "nil", nil, 0}, Token{Type: EOF, Pos: 3}}, nil},
		{"(", []Token{Token{LeftParen, "(", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{")", []Token{Token{RightParen, ")", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{".", []Token{Token{Period, ".", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},
		{",", []Token{Token{Comma, ",", nil, 0}, Token{Type: EOF, Pos: 1}}, nil},

		{"(1 + 2)", []Token{
			Token{LeftParen, "(", nil, 0},
			Token{Integer, "1", int64(1), 1},
			Token{Add, "+", nil, 3},
			Token{Integer, "2", int64(2), 5},
			Token{RightParen, ")", nil, 6},
			Token{Type: EOF, Pos: 7},
		}, nil},

		{
			"foo.bar(x, y, 5)",
			[]Token{
				Token{Identifier, "foo", nil, 0},
				Token{Period, ".", nil, 3},
				Token{Identifier, "bar", nil, 4},
				Token{LeftParen, "(", nil, 7},
				Token{Identifier, "x", nil, 8},
				Token{Comma, ",", nil, 9},
				Token{Identifier, "y", nil, 11},
				Token{Comma, ",", nil, 12},
				Token{Integer, "5", int64(5), 14},
				Token{RightParen, ")", nil, 15},
				Token{Type: EOF, Pos: 16},
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			s := NewScanner(test.src)
			tokens, err := s.Scan()

			if diff := deep.Equal(err, test.expectedError); diff != nil {
				t.Error(diff)
			}

			if diff := deep.Equal(tokens, test.expectedTokens); diff != nil {
				t.Error(diff)
			}
		})
	}

}
