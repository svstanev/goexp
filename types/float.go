package types

type Float float64

func NewFloat(value float64) Float {
	return Float(value)
}
