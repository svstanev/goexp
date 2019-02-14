package types

// Adder interface
type Adder interface {
	Add(other interface{}) (interface{}, error)
}

// Subtractor interface
type Subtractor interface {
	Sub(other interface{}) (interface{}, error)
}

// Multiplexor interface
type Multiplexor interface {
	Mul(other interface{}) (interface{}, error)
}

// Divider interface
type Divider interface {
	Div(other interface{}) (interface{}, error)
}

// Moduler interface
type Moduler interface {
	Mod(other interface{}) (interface{}, error)
}

// Inverter interface
type Inverter interface {
	Not() (interface{}, error)
}

// Negator interface
type Negator interface {
	Negate() (interface{}, error)
}

// EqualityComparer interface
type EqualityComparer interface {
	Equals(other interface{}) (bool, error)
}

// Comparer interface
type Comparer interface {
	Compare(other interface{}) (int, error)
}

type BooleanConverter interface {
	ToBoolean() Boolean
}
