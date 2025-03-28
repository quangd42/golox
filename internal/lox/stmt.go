// Generated with AST generator.

package lox

type stmt interface {
	accept(visitor stmtVisitor) error
}

type stmtVisitor interface {
	visitExprStmt(e exprStmt) error
	visitFunctionStmt(e functionStmt) error
	visitIfStmt(e ifStmt) error
	visitPrintStmt(e printStmt) error
	visitReturnStmt(e returnStmt) error
	visitVarStmt(e varStmt) error
	visitWhileStmt(e whileStmt) error
	visitForStmt(e forStmt) error
	visitBreakStmt(e breakStmt) error
	visitContinueStmt(e continueStmt) error
	visitBlockStmt(e blockStmt) error
	visitClassStmt(e classStmt) error
}

type exprStmt struct {
	expr expr
}

func (e exprStmt) accept(v stmtVisitor) error {
	return v.visitExprStmt(e)
}

type functionStmt struct {
	name    token
	literal functionExpr
}

func (e functionStmt) accept(v stmtVisitor) error {
	return v.visitFunctionStmt(e)
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

type returnStmt struct {
	keyword token
	value   expr
}

func (e returnStmt) accept(v stmtVisitor) error {
	return v.visitReturnStmt(e)
}

type varStmt struct {
	name        token
	initializer expr
}

func (e varStmt) accept(v stmtVisitor) error {
	return v.visitVarStmt(e)
}

type whileStmt struct {
	condition expr
	body      stmt
	label     token
	increment stmt
}

func (e whileStmt) accept(v stmtVisitor) error {
	return v.visitWhileStmt(e)
}

type forStmt struct {
	initializer stmt
	whileBody   whileStmt
}

func (e forStmt) accept(v stmtVisitor) error {
	return v.visitForStmt(e)
}

type breakStmt struct {
	keyword token
	label   token
}

func (e breakStmt) accept(v stmtVisitor) error {
	return v.visitBreakStmt(e)
}

type continueStmt struct {
	keyword token
	label   token
}

func (e continueStmt) accept(v stmtVisitor) error {
	return v.visitContinueStmt(e)
}

type blockStmt struct {
	statements []stmt
}

func (e blockStmt) accept(v stmtVisitor) error {
	return v.visitBlockStmt(e)
}

type classStmt struct {
	name       token
	superclass variableExpr
	methods    []functionStmt
}

func (e classStmt) accept(v stmtVisitor) error {
	return v.visitClassStmt(e)
}
