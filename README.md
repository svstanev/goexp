# goexp

Recursive descent expression parser in Go

## Installation

```
go get -u github.com/svstanev/goexp
```

## Usage

```golang
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
	context.AddMethod("max", func(args ...types.Integer) (interface{}, error) {
		var res int64 = math.MinInt64
		for _, value := range args {
			n := int64(value)
			if n > res {
				res = n
			}
		}
		return types.Integer(res), nil
	})

	res, err := goexp.EvalString("max(x, y, z)", context)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```

## Expression language

### Syntax Grammar

```
expression      -> logical_or
logical_or      -> logical_and (("||") logical_and)*;
logical_and     -> equality (("&&") equality)*;
equality        -> comparison (("==" | "!=") comparison)*;
comparison      -> addition (("<" | "<=" | ">" | ">=") addition)*;
addition        -> multiplication (("+" | "-") multiplication)*;
multiplication  -> unary (("*" | "/") unary)*;
unary           -> ("!" | "-")? call;
call            -> primary (("(" arguments? ")") | ("." IDENTIFIER))*;
primary         -> "false" | "true" | "nil" | IDENTIFIER | NUMBER | STRING | "(" expression ")";

arguments       -> expression ("," expression)*;
```

### Lexical Grammar

```
IDENTIFIER      -> ALPHA (ALPHA | DIGIT)*;
NUMBER          -> DIGIT* ("." DIGIT*)?;
STRING          -> "'" <any char except "'">* "'"
                  | '"' <any char except '"'>* '"';

DIGIT           -> '0'...'9'
ALPHA           -> 'a'...'z'|'A'...'Z'|'_'
```

