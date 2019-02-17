package types

import "math"

type Float float64

func NewFloat(value float64) Float {
	return Float(value)
}

func (f Float) Power(other interface{}) (interface{}, error) {
	x := float64(f)
	switch other.(type) {
	case Integer:
		y := int64(other.(Integer))
		return Float(math.Pow(x, float64(y))), nil

	case Float:
		y := float64(other.(Float))
		return Float(math.Pow(x, y)), nil

	default:
		return nil, notSupportedOperationError("**", x, other)
	}
}
