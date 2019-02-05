package goexp

import (
	"testing"
)

var printExprTests = []struct {
	expected string
	expr     Expr
}{
	{
		"1 + 2",
		BinaryExpr{
			Left:     IntegerLiteralExpr{int64(1)},
			Operator: Token{Type: Add},
			Right:    IntegerLiteralExpr{int64(2)},
		},
	},
	{
		"a",
		IdentifierExpr{"a", nil},
	},
	{
		"a.b.c",
		IdentifierExpr{
			Name: "c",
			Expr: IdentifierExpr{
				Name: "b",
				Expr: IdentifierExpr{
					Name: "a",
				},
			},
		},
	},
	{
		"add(1, 2)",
		CallExpr{
			Name: IdentifierExpr{"add", nil},
			Args: []Expr{
				IntegerLiteralExpr{int64(1)},
				IntegerLiteralExpr{int64(2)},
			},
		},
	},
}

func TestPrintExpr(t *testing.T) {
	for i, tt := range printExprTests {
		t.Run(tt.expected, func(t *testing.T) {
			s, err := Print(tt.expr)
			if err != nil {
				t.Errorf("[#%d] Error: %s", i, err.Error())
			}

			if s != tt.expected {
				t.Errorf("[#%d] Expected %s but got %s", i, tt.expected, s)
			}
		})
	}
}
