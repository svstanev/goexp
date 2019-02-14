package types

type NullType struct{}

var null = &NullType{}

func Null() interface{} {
	return null
}

func IsNull(value interface{}) bool {
	return value == null
}
