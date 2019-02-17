package goexp

import (
	"fmt"
	"reflect"

	"github.com/svstanev/goexp/types"
)

type Var interface {
	Value() (interface{}, error)
}

type varx struct {
	value interface{}
}

func (v varx) Value() (interface{}, error) {
	return v.value, nil
}

type Method interface {
	Invoke(args []interface{}) (interface{}, error)
}

type methodx struct {
	fn interface{}
}

func (m methodx) Invoke(args []interface{}) (interface{}, error) {
	fn := reflect.ValueOf(m.fn)
	in := make([]reflect.Value, len(args))
	for i, value := range args {
		in[i] = reflect.ValueOf(value)
	}
	res := fn.Call(in)
	return res[0].Interface(), nil
}

type Context interface {
	ResolveName(name string) (Var, bool)
	ResolveMethod(name string) (Method, bool)
}

type EvalContext interface {
	Context
	AddName(name string, value interface{}) error
	AddMethod(name string, fn interface{}) error
}

type context struct {
	parent  Context
	vars    map[string]Var
	methods map[string]Method
}

func NewEvalContext(parent Context) EvalContext {
	return &context{
		parent:  parent,
		vars:    make(map[string]Var),
		methods: make(map[string]Method),
	}
}

func (ctx *context) AddName(name string, value interface{}) error {
	if _, present := ctx.vars[name]; present {
		return fmt.Errorf("Var %s already exists", name)
	}
	ctx.vars[name] = varx{value}
	return nil
}

func (ctx *context) AddMethod(name string, fn interface{}) error {
	if _, present := ctx.methods[name]; present {
		return fmt.Errorf("Method %s already exists", name)
	}
	ctx.methods[name] = methodx{fn}
	return nil
}

func (ctx *context) ResolveName(name string) (Var, bool) {
	if x, ok := ctx.vars[name]; ok {
		return x, true
	}
	if ctx.parent != nil {
		return ctx.parent.ResolveName(name)
	}
	return nil, false
}

func (ctx *context) ResolveMethod(name string) (Method, bool) {
	if m, ok := ctx.methods[name]; ok {
		return m, true
	}
	if ctx.parent != nil {
		return ctx.parent.ResolveMethod(name)
	}
	return nil, false
}

// type Function func(args ...[]interface{}) (interface{}, error)

type interpreter struct {
	context Context
}

func newInterpreter(context Context) *interpreter {
	return &interpreter{
		context,
	}
}

func (i *interpreter) eval(expr Expr) (interface{}, error) {
	e := newEvaluator()
	return e.Eval(expr, i.context)
}

type evaluator struct{}

func newEvaluator() *evaluator {
	return &evaluator{}
}

func (eval *evaluator) Eval(expr Expr, context VisitorContext) (res interface{}, err error) {
	return expr.Accept(eval, context)
}

func (eval *evaluator) VisitStringLiteralExpr(expr StringLiteralExpr, context VisitorContext) (interface{}, error) {
	return types.NewString(expr.Value), nil
}

func (eval *evaluator) VisitIntegerLiteralExpr(expr IntegerLiteralExpr, context VisitorContext) (interface{}, error) {
	return types.NewInteger(expr.Value), nil
}

func (eval *evaluator) VisitFloatLiteralExpr(e FloatLiteralExpr, context VisitorContext) (interface{}, error) {
	return types.NewFloat(e.Value), nil
}

func (eval *evaluator) VisitBooleanLiteralExpr(e BooleanLiteralExpr, context VisitorContext) (interface{}, error) {
	return types.NewBoolean(e.Value), nil
}

func (eval *evaluator) VisitNilLiteralExpr(e NilLiteralExpr, context VisitorContext) (interface{}, error) {
	return types.Null(), nil
}

func (eval *evaluator) VisitBinaryExpr(e BinaryExpr, context VisitorContext) (interface{}, error) {
	var left, right interface{}
	var err error
	if left, err = eval.Eval(e.Left, context); err != nil {
		return nil, err
	}
	if right, err = eval.Eval(e.Right, context); err != nil {
		return nil, err
	}
	return binaryOp(left, right, e.Operator)
}

func (eval *evaluator) VisitUnaryExpr(e UnaryExpr, context VisitorContext) (interface{}, error) {
	var value interface{}
	var err error

	if value, err = eval.Eval(e.Value, context); err != nil {
		return nil, err
	}
	return unaryOp(value, e.Operator)
}

func (eval *evaluator) EvalMany(exprs []Expr, context VisitorContext) ([]interface{}, error) {
	var values = make([]interface{}, 0)
	var value interface{}
	var err error

	for _, x := range exprs {
		if value, err = eval.Eval(x, context); err == nil {
			values = append(values, value)
		} else {
			return nil, err
		}
	}

	return values, nil
}

func (eval *evaluator) VisitCallExpr(e CallExpr, context VisitorContext) (interface{}, error) {
	var err error
	id, ok := e.Name.(IdentifierExpr)
	if !ok {
		return nil, fmt.Errorf("Expected IdentifierExpr")
	}

	var val interface{} = context
	if id.Expr != nil {
		val, err = eval.Eval(id.Expr, context)
		if err != nil {
			return nil, err
		}
	}

	ctx, ok := val.(Context)
	if !ok {
		return nil, fmt.Errorf("Cannot resolve method %s", id.Name)
	}

	m, ok := ctx.ResolveMethod(id.Name)
	if !ok {
		return nil, fmt.Errorf("Method not found %s", id.Name)
	}

	args, err := eval.EvalMany(e.Args, context)
	if err != nil {
		return nil, err
	}

	return m.Invoke(args)
}

func (eval *evaluator) VisitIdentifierExpr(e IdentifierExpr, context VisitorContext) (interface{}, error) {
	var val interface{} = context
	if e.Expr != nil {
		var err error
		if val, err = eval.Eval(e.Expr, context); err != nil {
			return nil, err
		}
	}
	if ctx, ok := val.(Context); ok {
		if n, present := ctx.ResolveName(e.Name); present {
			return n.Value()
		}
		return nil, fmt.Errorf("%s not defined", e.Name)
	}
	return nil, fmt.Errorf("Cannot resolve %s", e.Name)
}

func (eval *evaluator) VisitGroupingExpr(e GroupingExpr, context VisitorContext) (interface{}, error) {
	return eval.Eval(e.Expr, context)
}

// type Lazy func() (interface{}, error)

func unaryOp(x interface{}, op Token) (interface{}, error) {
	switch op.Type {
	case Not:
		return not(x)
	case Sub:
		return negate(x)
	default:
		return nil, fmt.Errorf("Unknown operation: %s", op.Lexeme)
	}
}

func binaryOp(x, y interface{}, op Token) (interface{}, error) {
	switch op.Type {
	case Add:
		return add(x, y)
	case Sub:
		return sub(x, y)
	case Mul:
		return mul(x, y)
	case Div:
		return div(x, y)
	case Modulo:
		return mod(x, y)
	case Power:
		return pow(x, y)
	case Less:
		return lt(x, y)
	case LessEqual:
		return lte(x, y)
	case Greater:
		return gt(x, y)
	case GreaterEqual:
		return gte(x, y)
	case Equal:
		return equals(x, y)
	case NotEqual:
		return ne(x, y)
	case And:
		return and(x, y)
	case Or:
		return or(x, y)
	default:
		return nil, fmt.Errorf("Unknown operation: %s", op.Lexeme)
	}
}

func toBoolean(x interface{}) (b types.Boolean, ok bool) {
	if b, ok = x.(types.Boolean); ok {
		// all good
	} else if bc, ok := x.(types.BooleanConverter); ok {
		b = bc.ToBoolean()
	} else {
		ok = false
	}
	return
}

func and(x, y interface{}) (res interface{}, err error) {
	var l, r types.Boolean
	var ok bool
	if l, ok = toBoolean(x); !ok {
		err = fmt.Errorf("Cannot convert %v to Boolean", x)
	} else if r, ok = toBoolean(y); !ok {
		err = fmt.Errorf("Cannot convert %v to Boolean", y)
	} else {
		res = l.And(r)
	}
	return
}

func or(x, y interface{}) (res interface{}, err error) {
	var l, r types.Boolean
	var ok bool
	if l, ok = toBoolean(x); !ok {
		err = fmt.Errorf("Cannot convert %v to Boolean", x)
	} else if r, ok = toBoolean(y); !ok {
		err = fmt.Errorf("Cannot convert %v to Boolean", y)
	} else {
		res = l.Or(r)
	}
	return
}

func ne(x, y interface{}) (bool, error) {
	var res bool
	var err error
	if res, err = equals(x, y); err != nil {
		return false, err
	}
	return !res, nil
}

func lt(x, y interface{}) (bool, error) {
	var res int
	var err error
	if res, err = compare(x, y); err != nil {
		return false, err
	}
	return res < 0, nil
}

func lte(x, y interface{}) (bool, error) {
	var res int
	var err error
	if res, err = compare(x, y); err != nil {
		return false, err
	}
	return res <= 0, nil
}

func gt(x, y interface{}) (bool, error) {
	var res int
	var err error
	if res, err = compare(x, y); err != nil {
		return false, err
	}
	return res > 0, nil
}

func gte(x, y interface{}) (bool, error) {
	var res int
	var err error
	if res, err = compare(x, y); err != nil {
		return false, err
	}
	return res >= 0, nil
}

func add(x, y interface{}) (result interface{}, err error) {
	if adder, ok := x.(types.Adder); ok {
		result, err = adder.Add(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "+")
	}
	return
}

func sub(x, y interface{}) (res interface{}, err error) {
	if subtractor, ok := x.(types.Subtractor); ok {
		res, err = subtractor.Sub(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "-")
	}
	return
}

func mul(x, y interface{}) (res interface{}, err error) {
	if multiplexor, ok := x.(types.Multiplexor); ok {
		res, err = multiplexor.Mul(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "*")
	}
	return
}

func div(x, y interface{}) (res interface{}, err error) {
	if divider, ok := x.(types.Divider); ok {
		res, err = divider.Div(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "/")
	}
	return
}

func mod(x, y interface{}) (res interface{}, err error) {
	if modulo, ok := x.(types.Moduler); ok {
		res, err = modulo.Mod(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "%")
	}
	return
}

func pow(x, y interface{}) (res interface{}, err error) {
	if pow, ok := x.(types.SupportsPower); ok {
		res, err = pow.Power(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "**")
	}
	return
}

func not(x interface{}) (res interface{}, err error) {
	if b, ok := toBoolean(x); ok {
		res, err = b.Not(), nil
	} else {
		err = unaryOpNotSupportedError(x, "not")
	}
	return
}

func negate(x interface{}) (res interface{}, err error) {
	if negator, ok := x.(types.Negator); ok {
		res, err = negator.Negate()
	} else {
		err = unaryOpNotSupportedError(x, "-")
	}
	return
}

func equals(x, y interface{}) (res bool, err error) {
	if ec, ok := x.(types.EqualityComparer); ok {
		res, err = ec.Equals(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "==")
	}
	return
}

func compare(x, y interface{}) (res int, err error) {
	if comparer, ok := x.(types.Comparer); ok {
		res, err = comparer.Compare(y)
	} else {
		err = binaryOpNotSupportedError(x, y, "cmp")
	}
	return
}

type notSupportedOpError struct {
	x, y interface{}
	op   string
}

func binaryOpNotSupportedError(x, y interface{}, op string) notSupportedOpError {
	return notSupportedOpError{x, y, op}
}

func unaryOpNotSupportedError(x interface{}, op string) notSupportedOpError {
	return notSupportedOpError{x, nil, op}
}

func (err notSupportedOpError) Error() string {
	if err.y != nil {
		return fmt.Sprintf("Operation \"%v\" not supported for types %T and %T", err.op, err.x, err.y)
	}
	return fmt.Sprintf("Operation \"%v\" not supported for type %T", err.op, err.x)
}
