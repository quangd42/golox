// Generated with AST generator.

package lox

type expr interface {
	accept(visitor exprVisitor) (any, error)
}

type exprVisitor interface {
	visitBinaryExpr(e binaryExpr) (any, error)
	visitCallExpr(e callExpr) (any, error)
	visitGetExpr(e getExpr) (any, error)
	visitGroupingExpr(e groupingExpr) (any, error)
	visitLiteralExpr(e literalExpr) (any, error)
	visitLogicalExpr(e logicalExpr) (any, error)
	visitSetExpr(e setExpr) (any, error)
	visitUnaryExpr(e unaryExpr) (any, error)
	visitVariableExpr(e variableExpr) (any, error)
	visitAssignExpr(e assignExpr) (any, error)
}

type binaryExpr struct {
	left     expr
	operator token
	right    expr
}

func (e binaryExpr) accept(v exprVisitor) (any, error) {
	return v.visitBinaryExpr(e)
}

type callExpr struct {
	callee    expr
	paren     token
	arguments []expr
}

func (e callExpr) accept(v exprVisitor) (any, error) {
	return v.visitCallExpr(e)
}

type getExpr struct {
	object expr
	name   token
}

func (e getExpr) accept(v exprVisitor) (any, error) {
	return v.visitGetExpr(e)
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

type logicalExpr struct {
	left     expr
	operator token
	right    expr
}

func (e logicalExpr) accept(v exprVisitor) (any, error) {
	return v.visitLogicalExpr(e)
}

type setExpr struct {
	object expr
	name   token
	value  expr
}

func (e setExpr) accept(v exprVisitor) (any, error) {
	return v.visitSetExpr(e)
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

type assignExpr struct {
	name  token
	value expr
}

func (e assignExpr) accept(v exprVisitor) (any, error) {
	return v.visitAssignExpr(e)
}
