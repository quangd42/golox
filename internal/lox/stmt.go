// Generated with AST generator.

package lox

type stmt interface {
	accept(visitor stmtVisitor) error
}

type stmtVisitor interface {
	visitExprStmt(e exprStmt) error
	visitIfStmt(e ifStmt) error
	visitPrintStmt(e printStmt) error
	visitVarStmt(e varStmt) error
	visitBlockStmt(e blockStmt) error
}

type exprStmt struct {
	expr expr
}

func (e exprStmt) accept(v stmtVisitor) error {
	return v.visitExprStmt(e)
}

type ifStmt struct {
	condition  expr
	thenBranch stmt
	elseBranch stmt
}

func (e ifStmt) accept(v stmtVisitor) error {
	return v.visitIfStmt(e)
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

type blockStmt struct {
	statements []stmt
}

func (e blockStmt) accept(v stmtVisitor) error {
	return v.visitBlockStmt(e)
}
