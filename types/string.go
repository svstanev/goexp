package types

import (
	"fmt"
	"strings"
)

type String string

func NewString(value string) String {
	return String(value)
}

func (this String) String() string {
	return string(this)
}

func (this String) Add(other interface{}) (res interface{}, err error) {
	switch other.(type) {
	case Integer, Float, String, Boolean:
		res = String(fmt.Sprintf("%s%s", this, other))
	default:
		if IsNull(other) {
			res = this
		} else {
			err = fmt.Errorf("Cannot add %T to String", other)
		}
	}
	return
}

// Less returns true if the String s is less then the String other
func (s String) Less(other interface{}) (res bool, err error) {
	switch other.(type) {
	case String:
		str := other.(String)
		res = compare(s, str) < 0
	default:
		err = fmt.Errorf("Cannot compare String and %T", other)
	}
	return
}

// Equals return true if two strings match
func (s String) Equals(other interface{}) (res bool, err error) {
	switch other.(type) {
	case String:
		str := other.(String)
		res = compare(s, str) == 0
	default:
		err = fmt.Errorf("Cannot compare String and %T", other)
	}
	return
}

// Compare returns an integer comparing two strings lexicographically
func (s String) Compare(other interface{}) (res int, err error) {
	if str, isString := other.(String); isString {
		res = compare(s, str)
	} else {
		err = fmt.Errorf("Cannot compare String and %T", other)
	}
	return
}

func compare(s1, s2 String) int {
	return strings.Compare(string(s1), string(s2))
}
