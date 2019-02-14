package types

import (
	"fmt"
)

type Integer int64

func NewInteger(value int64) Integer {
	return Integer(value)
}

func (this Integer) Add(other interface{}) (interface{}, error) {
	switch other.(type) {
	case Integer:
		value := other.(Integer)
		return Integer(int64(this) + int64(value)), nil
	case Float:
		value := other.(Float)
		return Float(float64(int64(this)) + float64(value)), nil
	case String:
		return fmt.Sprintf("%d%s", this, other), nil
	default:
		return nil, fmt.Errorf("Cannot add %T and %T", this, other)
	}
}

func (n Integer) String() string {
	return fmt.Sprintf("Integer(%d)", int64(n))
}
