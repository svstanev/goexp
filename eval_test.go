package goexp

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/svstanev/goexp/types"
)

func reduce(values []types.Integer, f func(x, y types.Integer) types.Integer, init types.Integer) types.Integer {
	res := init
	for _, n := range values {
		res = f(res, n)
	}
	return res
}

func max(x, y types.Integer) types.Integer {
	if x > y {
		return x
	}
	return y
}

func TestEval(t *testing.T) {
	ctx := NewEvalContext(nil)
	ctx.AddName("x", types.Integer(1))
	ctx.AddName("y", types.Integer(2))
	ctx.AddMethod("max", func(values ...types.Integer) (types.Integer, error) {
		return reduce(values, max, math.MinInt64), nil
	})

	tests := []struct {
		expr   string
		result interface{}
		err    error
	}{
		{"1 + 2", types.Integer(3), nil},
		{"max(x + y, 5, 3 * x * y)", types.Integer(6), nil},
		{"2 ** 3", types.Integer(8), nil},
		{"2.5 ** 3", types.Float(15.625), nil},

		{"1 < 'foo'", nil, binaryOpNotSupportedError(types.Integer(1), types.String("foo"), "cmp")},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%v %v", i, test.expr), func(t *testing.T) {
			res, err := eval(test.expr, ctx)
			if !reflect.DeepEqual(err, test.err) {
				t.Fatalf(`Expected "%v" error but got "%v" error`, test.err, err)
			}
			if err == nil && test.err == nil && !reflect.DeepEqual(res, test.result) {
				t.Fatalf(`Expected %v but got %v`, test.result, res)
			}
		})
	}

	// res, err := eval("max(x + y, 5, 3 * x * y)", ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// expected := types.Integer(6)
	// if res != expected {
	// 	t.Fatalf("Expected %#v but got %#v", expected, res)
	// }
}

func eval(expr string, context Context) (interface{}, error) {
	scanner := newScanner(expr)
	tokens, err := scanner.scan()
	if err != nil {
		return nil, err
	}

	parser := newParser(tokens)
	x, err := parser.parse()
	if err != nil {
		return nil, err
	}

	interpreter := newInterpreter(context)
	return interpreter.eval(x)
}
