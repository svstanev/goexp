package goexp

// Parse the given string and returns the expression's AST
func Parse(expr string) (Expr, error) {
	scanner := newScanner(expr)
	tokens, err := scanner.scan()
	if err != nil {
		return nil, err
	}

	parser := newParser(tokens)
	return parser.parse()
}

// Eval returns the result of the evaluation of the given expression
func Eval(expr Expr, context Context) (interface{}, error) {
	in := newInterpreter(context)
	return in.eval(expr)
}

// EvalString parses and then evaluates the given string
func EvalString(s string, context Context) (interface{}, error) {
	expr, err := Parse(s)
	if err != nil {
		return nil, err
	}
	in := newInterpreter(context)
	return in.eval(expr)
}
