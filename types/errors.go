package types

import "fmt"

func notSupportedOperationError(op string, x, y interface{}) error {
	return fmt.Errorf(`Operation "%s" not supported for %T and %T`, op, x, y)
}
