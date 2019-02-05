package goexp

import (
	"errors"
	"fmt"
)

var ops = map[TokenType]string{
	Add:          "+",
	Sub:          "-",
	Mul:          "*",
	Div:          "/",
	Modulo:       "%",
	Less:         "<",
	LessEqual:    "<=",
	Greater:      ">",
	GreaterEqual: ">=",
	Equal:        "==",
	NotEqual:     "!=",
	Not:          "!",
	And:          "&&",
	Or:           "||",
}

func Print(node Expr) (string, error) {
	p := newPrinter()

	var res interface{}
	var err error
	var s string

	context := newPrinterContext()
	if res, err = node.Accept(p, context); err == nil {
		s = res.(string)
	}

	return s, err
}

type printer struct {
	Visitor
}

type printerContext struct {
	VisitorContext
}

func newPrinter() *printer {
	return &printer{}
}

func newPrinterContext() *printerContext {
	return &printerContext{}
}

func (p *printer) VisitGroupingExpr(e GroupingExpr, context VisitorContext) (interface{}, error) {
	s, err := p.printExpr(e.Expr, context)
	if err != nil {
		return "", err
	}
	return "(" + s + ")", nil
}

func (p *printer) VisitStringLiteralExpr(e StringLiteralExpr, c VisitorContext) (interface{}, error) {
	return fmt.Sprintf("\"%s\"", e.Value), nil
}

func (p *printer) VisitIntegerLiteralExpr(e IntegerLiteralExpr, context VisitorContext) (interface{}, error) {
	return fmt.Sprintf("%d", e.Value), nil
}

func (p *printer) VisitFloatLiteralExpr(e FloatLiteralExpr, context VisitorContext) (interface{}, error) {
	return fmt.Sprintf("%f", e.Value), nil
}

func (p *printer) VisitBooleanLiteralExpr(e BooleanLiteralExpr, context VisitorContext) (interface{}, error) {
	var res string

	if e.Value {
		res = "True"
	} else {
		res = "False"
	}

	return res, nil
}

func (p *printer) VisitNilLiteralExpr(e NilLiteralExpr, context VisitorContext) (interface{}, error) {
	return "Nil", nil
}

func (p *printer) VisitUnaryExpr(e UnaryExpr, context VisitorContext) (interface{}, error) {
	str, err := p.printExpr(e.Value, context)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s%s", ops[e.Operator.Type], str), nil
}

func (p *printer) VisitBinaryExpr(e BinaryExpr, context VisitorContext) (interface{}, error) {
	var err error
	var l, r string

	if l, err = p.printExpr(e.Left, context); err != nil {
		return nil, err
	}

	if r, err = p.printExpr(e.Right, context); err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s %s %s", l, ops[e.Operator.Type], r), nil
}

func (p *printer) VisitCallExpr(e CallExpr, context VisitorContext) (interface{}, error) {
	name, err := p.printExpr(e.Name, context)
	if err != nil {
		return nil, err
	}

	var res = ""
	res += name
	res += "("

	for i, arg := range e.Args {
		s, err := p.printExpr(arg, context)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			res += ", "
		}
		res += s
	}

	res += ")"
	return res, nil
}

func (p *printer) VisitIdentifierExpr(e IdentifierExpr, context VisitorContext) (interface{}, error) {
	var res = ""
	if e.Expr != nil {
		var s string
		var err error
		if s, err = p.printExpr(e.Expr, context); err != nil {
			return "", err
		}
		res += s
	}
	if len(res) > 0 {
		res += "."
	}
	res += e.Name

	return res, nil
}

func (p *printer) printExpr(expr Expr, context VisitorContext) (string, error) {
	var err error
	var res interface{}
	if res, err = expr.Accept(p, context); err == nil {
		if s, ok := res.(string); ok {
			return s, nil
		}
		err = errors.New("Unable to print expression")
	}
	return "", err
}
