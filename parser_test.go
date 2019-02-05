package goexp

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParse(t *testing.T) {
	tests := []struct {
		str  string
		expr Expr
		err  error
	}{
		{
			"1 + 2",
			BinaryExpr{
				Left:     IntegerLiteralExpr{int64(1)},
				Right:    IntegerLiteralExpr{int64(2)},
				Operator: Token{Type: Add, Lexeme: "+", Pos: 2},
			},
			nil,
		},

		{
			"5 - 2 - 1",
			BinaryExpr{
				Left: BinaryExpr{
					Left:     IntegerLiteralExpr{int64(5)},
					Right:    IntegerLiteralExpr{int64(2)},
					Operator: Token{Sub, "-", nil, 2},
				},
				Right:    IntegerLiteralExpr{int64(1)},
				Operator: Token{Sub, "-", nil, 6},
			},
			nil,
		},

		{
			"1 + 2 == 3",
			BinaryExpr{
				Left: BinaryExpr{
					Left:     IntegerLiteralExpr{int64(1)},
					Right:    IntegerLiteralExpr{int64(2)},
					Operator: Token{Type: Add, Lexeme: "+", Pos: 2},
				},
				Right:    IntegerLiteralExpr{int64(3)},
				Operator: Token{Equal, "==", nil, 6},
			},
			nil,
		},

		{
			"(a || b) && c",
			BinaryExpr{
				Left: GroupingExpr{
					BinaryExpr{
						Left:     IdentifierExpr{Name: "a"},
						Right:    IdentifierExpr{Name: "b"},
						Operator: Token{Or, "||", nil, 3},
					}},
				Right:    IdentifierExpr{Name: "c"},
				Operator: Token{And, "&&", nil, 9},
			},
			nil,
		},

		{
			"foo()",
			CallExpr{
				Name: IdentifierExpr{"foo", nil},
				Args: []Expr{},
			},
			nil,
		},

		{
			"foo.bar(1, 2, 3)",
			CallExpr{
				Name: IdentifierExpr{"bar", IdentifierExpr{"foo", nil}},
				Args: []Expr{
					IntegerLiteralExpr{int64(1)},
					IntegerLiteralExpr{int64(2)},
					IntegerLiteralExpr{int64(3)},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			testParse(t, test.str, test.expr, test.err)
		})
	}
}

func testParse(t *testing.T, str string, expectedExpr Expr, expectedErr error) {
	s := NewScanner(str)
	tokens, err := s.Scan()
	if err != nil {
		t.Error(err)
	}

	p := NewParser(tokens)
	expr, err := p.Parse()

	if diff := deep.Equal(err, expectedErr); diff != nil {
		t.Error(diff)
	}

	if diff := deep.Equal(expr, expectedExpr); diff != nil {
		t.Error(diff)
	}
}
