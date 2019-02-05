package goexp

/*
Expr interface for all types expression nodes
*/
type Expr interface {
	exprNode()

	Accept(v Visitor, context VisitorContext) (res interface{}, err error)
}

/*
StringLiteralExpr represents a string literal
*/
type StringLiteralExpr struct {
	Value string
}

type IntegerLiteralExpr struct {
	Value int64
}

type FloatLiteralExpr struct {
	Value float64
}

type BooleanLiteralExpr struct {
	Value bool
}

type NilLiteralExpr struct {
}

type GroupingExpr struct {
	Expr Expr
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

type UnaryExpr struct {
	Value    Expr
	Operator Token
}

type CallExpr struct {
	Name Expr
	Args []Expr
}

type IdentifierExpr struct {
	Name string
	Expr Expr
}

func (StringLiteralExpr) exprNode()  {}
func (IntegerLiteralExpr) exprNode() {}
func (FloatLiteralExpr) exprNode()   {}
func (BooleanLiteralExpr) exprNode() {}
func (NilLiteralExpr) exprNode()     {}
func (UnaryExpr) exprNode()          {}
func (BinaryExpr) exprNode()         {}
func (CallExpr) exprNode()           {}
func (IdentifierExpr) exprNode()     {}
func (GroupingExpr) exprNode()       {}

func (s StringLiteralExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitStringLiteralExpr(s, context)
}

func (e IntegerLiteralExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitIntegerLiteralExpr(e, context)
}

func (e FloatLiteralExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitFloatLiteralExpr(e, context)
}

func (e BooleanLiteralExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitBooleanLiteralExpr(e, context)
}

func (e NilLiteralExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitNilLiteralExpr(e, context)
}

func (e BinaryExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitBinaryExpr(e, context)
}

func (e UnaryExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitUnaryExpr(e, context)
}

func (e IdentifierExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitIdentifierExpr(e, context)
}

func (e CallExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitCallExpr(e, context)
}

func (e GroupingExpr) Accept(v Visitor, context VisitorContext) (interface{}, error) {
	return v.VisitGroupingExpr(e, context)
}
