package goexp

import "fmt"

type parseError struct {
	token   Token
	message string
}

func (pe parseError) Error() string {
	return fmt.Sprintf("Parse Error at pos %d: %s", pe.token.Pos, pe.message)
}

type parser struct {
	tokens  []Token
	current int
}

func newParser(tokens []Token) *parser {
	return &parser{
		tokens: tokens,
	}
}

func (p *parser) parse() (Expr, error) {
	p.current = 0
	return p.expression()
}

func (p *parser) peek() Token {
	return p.tokens[p.current]
}

func (p *parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *parser) check(tokenType TokenType) bool {
	return !p.isAtEnd() && p.peek().Type == tokenType
}

func (p *parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) consume(tokenType TokenType, msg string) (token Token, err error) {
	if p.check(tokenType) {
		token = p.advance()
	} else {
		err = parseError{
			token:   p.peek(),
			message: msg,
		}
	}
	return
}

func (p *parser) expression() (Expr, error) {
	return p.logicalOr()
}

func (p *parser) logicalOr() (Expr, error) {
	// or = and "||" and

	left, err := p.logicalAnd()
	if err != nil {
		return nil, err
	}

	for p.match(Or) {
		op := p.previous()

		right, err := p.logicalAnd()
		if err != nil {
			return nil, err
		}

		left = BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	return left, nil
}

func (p *parser) logicalAnd() (Expr, error) {
	// logicalAnd = logicalNot "&&" logicalNot

	left, err := p.logicalNot()
	if err != nil {
		return nil, err
	}

	for p.match(And) {
		op := p.previous()
		right, err := p.logicalNot()
		if err != nil {
			return nil, err
		}
		left = BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	return left, nil
}

func (p *parser) logicalNot() (Expr, error) {
	// logicalNot = (Not)? equality
	var op Token
	if p.match(Not) {
		op = p.previous()
	}
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if op.Type != Unknown {
		expr = UnaryExpr{
			Value:    expr,
			Operator: op,
		}
	}
	return expr, nil
}

func (p *parser) equality() (Expr, error) {
	// equality = logicalNot (("==" | "!=") logicalNot)*
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(Equal, NotEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}

	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	// comparison = addition "<" | "<=" | ">" | ">=" addition
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}
	for p.match(Less, LessEqual, Greater, GreaterEqual) {
		op := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}
	return expr, nil
}

func (p *parser) addition() (Expr, error) {
	// addition = mult "+" | "-" mult
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}
	for p.match(Add, Sub) {
		op := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}

	return expr, nil
}

func (p *parser) multiplication() (Expr, error) {
	// addition = unary ("*" | "/" | "%") unary
	expr, err := p.power()
	if err != nil {
		return nil, err
	}
	for p.match(Mul, Div, Modulo) {
		op := p.previous()
		right, err := p.power()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}
	return expr, nil
}

func (p *parser) power() (Expr, error) {
	expr, err := p.negate()
	if err != nil {
		return nil, err
	}
	for p.match(Power) {
		op := p.previous()
		right, err := p.negate()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}
	return expr, nil
}

func (p *parser) negate() (Expr, error) {
	// unary =  ("-")? call

	var op Token
	if p.match(Sub) {
		op = p.previous()
	}
	expr, err := p.call()
	if err != nil {
		return nil, err
	}
	if op.Type != Unknown {
		expr = UnaryExpr{
			Operator: op,
			Value:    expr,
		}
	}
	return expr, nil
}

func (p *parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	done := false
	for !done {
		if p.match(LeftParen) {
			// function call
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(Period) {
			name, err := p.consume(Identifier, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = IdentifierExpr{name.Lexeme, expr}
		} else {
			done = true
		}
	}
	return expr, nil
}

func (p *parser) finishCall(callee Expr) (Expr, error) {
	args := make([]Expr, 0)
	if !p.check(RightParen) {
		done := false
		for !done {
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
			done = !p.match(Comma)
		}
	}
	_, err := p.consume(RightParen, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return CallExpr{callee, args}, nil
}

func (p *parser) primary() (Expr, error) {
	// primary = NUMBER | STRING | "false" | "true" | "nil" | "(" expression ")"

	if p.match(False) {
		return BooleanLiteralExpr{false}, nil
	}
	if p.match(True) {
		return BooleanLiteralExpr{true}, nil
	}
	if p.match(Nil) {
		return NilLiteralExpr{}, nil
	}
	if p.match(Integer) {
		value := p.previous().Literal.(int64)
		return IntegerLiteralExpr{value}, nil
	}
	if p.match(Float) {
		value := p.previous().Literal.(float64)
		return FloatLiteralExpr{value}, nil
	}
	if p.match(String) {
		value := p.previous().Literal.(string)
		return StringLiteralExpr{value}, nil
	}
	if p.match(Identifier) {
		return IdentifierExpr{Name: p.previous().Lexeme}, nil
	}
	if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RightParen, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return GroupingExpr{expr}, nil
	}

	return nil, parseError{p.previous(), "Unknown token"}
}
