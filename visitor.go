package goexp

type VisitorContext interface{}

type Visitable interface {
	Accept(visitor Visitor, context VisitorContext) (interface{}, error)
}

/*
Visitor interface
*/
type Visitor interface {
	VisitStringLiteralExpr(e StringLiteralExpr, context VisitorContext) (interface{}, error)
	VisitIntegerLiteralExpr(e IntegerLiteralExpr, context VisitorContext) (interface{}, error)
	VisitFloatLiteralExpr(e FloatLiteralExpr, context VisitorContext) (interface{}, error)
	VisitBooleanLiteralExpr(e BooleanLiteralExpr, context VisitorContext) (interface{}, error)
	VisitNilLiteralExpr(e NilLiteralExpr, context VisitorContext) (interface{}, error)
	VisitBinaryExpr(e BinaryExpr, context VisitorContext) (interface{}, error)
	VisitUnaryExpr(e UnaryExpr, context VisitorContext) (interface{}, error)
	VisitCallExpr(e CallExpr, context VisitorContext) (interface{}, error)
	VisitIdentifierExpr(e IdentifierExpr, context VisitorContext) (interface{}, error)
	VisitGroupingExpr(e GroupingExpr, context VisitorContext) (interface{}, error)
}
