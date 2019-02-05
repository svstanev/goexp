package goexp

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	token := Token{
		Type:    Integer,
		Lexeme:  "123",
		Literal: 123,
		Pos:     5,
	}
	s := fmt.Sprintf("%+v", token)
	if s != "{Type:21 Lexeme:123 Literal:123 Pos:5}" {
		t.Errorf("Invalid Token.String()")
	}
}
