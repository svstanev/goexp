package types

import (
	"fmt"
	"testing"
)

func TestIsNull(t *testing.T) {
	tests := []struct {
		isNull bool
		value  interface{}
	}{
		{true, Null()},

		{false, nil},
		{false, 1},
		{false, "abc"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, expected := IsNull(test.value), test.isNull
			if actual != expected {
				t.Errorf("Expected '%v' but got '%v'", expected, actual)
			}
		})
	}
}
