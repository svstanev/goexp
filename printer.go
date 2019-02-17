package goexp

import (
	"errors"
	"fmt"
	"strings"
)

var ops = map[TokenType]string{
	Add:          "+",
	Sub:          "-",
	Mul:          "*",
	Div:          "/",
	Modulo:       "%",
	Power:        "**",
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
	right, err := p.printExpr(e.Value, context)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%s%s", ops[e.Operator.Type], right), nil
}

func (p *printer) VisitBinaryExpr(e BinaryExpr, context VisitorContext) (interface{}, error) {
	left, err := p.printExpr(e.Left, context)
	if err != nil {
		return nil, err
	}

	right, err := p.printExpr(e.Right, context)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%s %s %s", left, ops[e.Operator.Type], right), nil
}

func (p *printer) VisitCallExpr(e CallExpr, context VisitorContext) (interface{}, error) {
	name, err := p.printExpr(e.Name, context)
	if err != nil {
		return nil, err
	}
	args, err := p.printArguments(e.Args, context)
	if err != nil {
		return nil, err
	}
	str := fmt.Sprintf("%s(%s)", name, args)
	return str, nil
}

func (p *printer) VisitIdentifierExpr(e IdentifierExpr, context VisitorContext) (interface{}, error) {
	var res = ""
	if e.Expr != nil {
		s, err := p.printExpr(e.Expr, context)
		if err != nil {
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
	res, err := expr.Accept(p, context)
	if err != nil {
		return "", err
	}
	str, ok := res.(string)
	if !ok {
		err = errors.New("Unable to print expression")
	}
	return str, err
}

func (p *printer) printArguments(args []Expr, context VisitorContext) (string, error) {
	strs, err := p.printMany(args, context)
	if err != nil {
		return "", err
	}
	return strings.Join(strs, ", "), nil
}

func (p *printer) printMany(exprs []Expr, context VisitorContext) ([]string, error) {
	var strs = make([]string, len(exprs))
	for i, x := range exprs {
		str, err := p.printExpr(x, context)
		if err != nil {
			return nil, err
		} else {
			strs[i] = str
		}
	}
	return strs, nil
}
