package types

import (
	"fmt"
	"testing"
)

var (
	TRUE  = Boolean(true)
	FALSE = Boolean(false)
)

func TestAnd(t *testing.T) {
	tests := []struct {
		left, right, expected Boolean
	}{
		{TRUE, TRUE, TRUE},
		{TRUE, FALSE, FALSE},
		{FALSE, TRUE, FALSE},
		{FALSE, FALSE, FALSE},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d: %t AND %t", i, test.left, test.right), func(t *testing.T) {
			result := test.left.And(test.right)
			if result != test.expected {
				t.Errorf("Need %t but got %t", test.expected, result)
			}
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		left, right, expected Boolean
	}{
		{TRUE, TRUE, TRUE},
		{TRUE, FALSE, TRUE},
		{FALSE, TRUE, TRUE},
		{FALSE, FALSE, FALSE},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d: %t OR %t", i, test.left, test.right), func(t *testing.T) {
			result := test.left.Or(test.right)
			if result != test.expected {
				t.Errorf("Need %t but got %t", test.expected, result)
			}
		})
	}
}
