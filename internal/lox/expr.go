// Generated with AST generator.

package lox

type expr interface {
	accept(visitor exprVisitor) (any, error)
}

type exprVisitor interface {
	visitBinaryExpr(e binaryExpr) (any, error)
	visitGroupingExpr(e groupingExpr) (any, error)
	visitLiteralExpr(e literalExpr) (any, error)
	visitUnaryExpr(e unaryExpr) (any, error)
	visitVariableExpr(e variableExpr) (any, error)
}

type binaryExpr struct {
	left     expr
	operator token
	right    expr
}

func (e binaryExpr) accept(v exprVisitor) (any, error) {
	return v.visitBinaryExpr(e)
}

type groupingExpr struct {
	expr expr
}

func (e groupingExpr) accept(v exprVisitor) (any, error) {
	return v.visitGroupingExpr(e)
}

type literalExpr struct {
	value any
}

func (e literalExpr) accept(v exprVisitor) (any, error) {
	return v.visitLiteralExpr(e)
}

type unaryExpr struct {
	operator token
	right    expr
}

func (e unaryExpr) accept(v exprVisitor) (any, error) {
	return v.visitUnaryExpr(e)
}

type variableExpr struct {
	name token
}

func (e variableExpr) accept(v exprVisitor) (any, error) {
	return v.visitVariableExpr(e)
}
