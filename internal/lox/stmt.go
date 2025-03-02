// Generated with AST generator.

package lox

type stmt interface {
	accept(visitor stmtVisitor) error
}

type stmtVisitor interface {
	visitExprStmt(e exprStmt) error
	visitPrintStmt(e printStmt) error
	visitVarStmt(e varStmt) error
}

type exprStmt struct {
	expr expr
}

func (e exprStmt) accept(v stmtVisitor) error {
	return v.visitExprStmt(e)
}

type printStmt struct {
	expr expr
}

func (e printStmt) accept(v stmtVisitor) error {
	return v.visitPrintStmt(e)
}

type varStmt struct {
	name        token
	initializer expr
}

func (e varStmt) accept(v stmtVisitor) error {
	return v.visitVarStmt(e)
}
