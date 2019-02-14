package types

type Boolean bool

func NewBoolean(value bool) Boolean {
	return Boolean(value)
}

func (b Boolean) And(other Boolean) Boolean {
	return Boolean(bool(b) && bool(other))
}

func (b Boolean) Or(other Boolean) Boolean {
	return Boolean(bool(b) || bool(other))
}

func (b Boolean) Not() Boolean {
	return Boolean(!bool(b))
}
