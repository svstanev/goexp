package goexp

import (
	"testing"

	"github.com/svstanev/goexp/types"
)

func TestEval(t *testing.T) {
	ctx := newContext(nil)
	ctx.addName("x", types.Integer(1))
	ctx.addName("y", types.Integer(2))

	res, err := eval("x + y", ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := types.Integer(3)
	if res != expected {
		t.Fatalf("Expected %#v but got %#v", expected, res)
	}
}

func eval(expr string, context *context) (interface{}, error) {
	scanner := NewScanner(expr)
	tokens, err := scanner.Scan()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	x, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	interpreter := newInterpreter(context)
	return interpreter.Eval(x)
}
