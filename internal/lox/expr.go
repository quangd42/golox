// Generated with AST generator.

package lox

type expr interface {
	accept(visitor exprVisitor) (any, error)
}

type exprVisitor interface {
	visitArrayExpr(e arrayExpr) (any, error)
	visitAssignExpr(e assignExpr) (any, error)
	visitBinaryExpr(e binaryExpr) (any, error)
	visitCallExpr(e callExpr) (any, error)
	visitFunctionExpr(e functionExpr) (any, error)
	visitGetExpr(e getExpr) (any, error)
	visitGroupingExpr(e groupingExpr) (any, error)
	visitIndexExpr(e indexExpr) (any, error)
	visitLiteralExpr(e literalExpr) (any, error)
	visitLogicalExpr(e logicalExpr) (any, error)
	visitSetExpr(e setExpr) (any, error)
	visitSuperExpr(e superExpr) (any, error)
	visitTernaryExpr(e ternaryExpr) (any, error)
	visitThisExpr(e thisExpr) (any, error)
	visitUnaryExpr(e unaryExpr) (any, error)
	visitVariableExpr(e variableExpr) (any, error)
}

type arrayExpr struct {
	value []expr
}

func (e arrayExpr) accept(v exprVisitor) (any, error) {
	return v.visitArrayExpr(e)
}

type assignExpr struct {
	name  token
	value expr
}

func (e assignExpr) accept(v exprVisitor) (any, error) {
	return v.visitAssignExpr(e)
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

type functionExpr struct {
	params []token
	body   []stmt
}

func (e functionExpr) accept(v exprVisitor) (any, error) {
	return v.visitFunctionExpr(e)
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

type indexExpr struct {
	callee  expr
	bracket token
	index   expr
}

func (e indexExpr) accept(v exprVisitor) (any, error) {
	return v.visitIndexExpr(e)
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

type superExpr struct {
	keyword token
	method  token
}

func (e superExpr) accept(v exprVisitor) (any, error) {
	return v.visitSuperExpr(e)
}

type ternaryExpr struct {
	condition expr
	thenExpr  expr
	elseExpr  expr
}

func (e ternaryExpr) accept(v exprVisitor) (any, error) {
	return v.visitTernaryExpr(e)
}

type thisExpr struct {
	keyword token
}

func (e thisExpr) accept(v exprVisitor) (any, error) {
	return v.visitThisExpr(e)
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
