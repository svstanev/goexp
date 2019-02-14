package goexp

import "fmt"

func ExampleEval() {
	if expr, err := Parse("1 + 2"); err == nil {
		res, err := Eval(expr, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	// Output: Integer(3)
}
