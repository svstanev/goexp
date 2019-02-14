package main

import (
	"fmt"
	"math"

	"github.com/svstanev/goexp"
	"github.com/svstanev/goexp/types"
)

func main() {
	context := goexp.NewEvalContext(nil)
	context.AddName("x", types.Integer(1))
	context.AddName("y", types.Integer(3))
	context.AddName("z", types.Integer(5))
	context.AddMethod("max", max)

	res, err := goexp.EvalString("max(x, y, z)", context)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func max(args ...types.Integer) (interface{}, error) {
	var res int64 = math.MinInt64
	for _, value := range args {
		n := int64(value)
		if n > res {
			res = n
		}
	}
	return types.Integer(res), nil
}
