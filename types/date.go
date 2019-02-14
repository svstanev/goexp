package types

import (
	"fmt"
	"time"
)

// Date in ms since 1.1.1970 00:00.00
type Date time.Time

func (date Date) Add(other interface{}) (interface{}, error) {
	switch other.(type) {
	case Integer:
		n := other.(Integer)
		return Date(time.Time(date).Add(time.Duration(n))), nil

	case Duration:
		n := other.(Duration)
		return Date(time.Time(date).Add(time.Duration(n))), nil

	default:
		return nil, fmt.Errorf("Add operation not supported for types %T and %T", date, other)
	}
}
