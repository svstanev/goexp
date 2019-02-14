package goexp

import (
	"errors"
	"fmt"
)

type Type interface {
}

type Object interface {
	getType() Type
}

type TInteger int64
type TFloat float64
type TBoolean bool
type TString string

func NewInteger(value int64) TInteger {
	return TInteger(value)
}

func (this TInteger) Add(other interface{}) (interface{}, error) {
	switch other.(type) {
	case TInteger:
		value := other.(TInteger)
		return TInteger(int64(this) + int64(value)), nil
	case TFloat:
		value := other.(TFloat)
		return TFloat(float64(int64(this)) + float64(value)), nil
	case TString:
		return fmt.Sprintf("%d%s", this, other), nil
	default:
		return nil, errors.New(fmt.Sprintf("Cannot add %T and %T", this, other))
	}
}

func NewFloat(value float64) TFloat {
	return TFloat(value)
}

func NewBoolean(value bool) TBoolean {
	return TBoolean(value)
}
