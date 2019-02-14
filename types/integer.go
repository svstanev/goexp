package types

import (
	"fmt"
)

type Integer int64

func NewInteger(value int64) Integer {
	return Integer(value)
}

// Add returns the sum of the current and the other value
func (n Integer) Add(other interface{}) (interface{}, error) {
	x := int64(n)
	switch other.(type) {
	case Integer:
		y := other.(Integer)
		return Integer(x + int64(y)), nil
	case Float:
		value := other.(Float)
		return Float(float64(x) + float64(value)), nil
	case String:
		return fmt.Sprintf("%d%s", n, other), nil
	default:
		return nil, notSupportedOperationError("+", n, other)
	}
}

// Mul returns the product of the current and the other value
func (n Integer) Mul(other interface{}) (interface{}, error) {
	x := int64(n)
	switch other.(type) {
	case Integer:
		y := int64(other.(Integer))
		return Integer(x * y), nil
	case Float:
		y := float64(other.(float64))
		return Float(float64(x) * y), nil
	default:
		return nil, notSupportedOperationError("*", n, other)
	}
}

func (n Integer) Div(other interface{}) (interface{}, error) {
	x := int64(n)
	switch other.(type) {
	case Integer:
		y := int64(other.(Integer))
		return Integer(x / y), nil
	case Float:
		y := float64(other.(float64))
		return Float(float64(x) / y), nil
	default:
		return nil, notSupportedOperationError("/", n, other)
	}
}

func (n Integer) String() string {
	return fmt.Sprintf("Integer(%d)", int64(n))
}
