package lox

import "fmt"

type Resolver struct {
	interpreter *Interpreter
	scopes      *scopeStack
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      newScopeStack(),
	}
}

func (r *Resolver) Resolve(stmts []stmt) error {
	return r.resolveStmtList(stmts)
}

func (r *Resolver) resolveStmtList(stmts []stmt) error {
	for _, stmt := range stmts {
		// TODO: after logger is added, statements after the error
		// should also be resolved
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveExpr(e expr) (any, error) {
	return e.accept(r)
}

func (r *Resolver) resolveStmt(s stmt) error {
	return s.accept(r)
}

func (r *Resolver) resolveFunction(s functionStmt) error {
	r.beginScope()
	defer r.endScope()
	for _, param := range s.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmtList(s.body)
	// NOTE: Explicitly return nil so that endScope() runs AFTER resolveStmtList()
	return nil
}

func (r *Resolver) resolveLocal(e expr, name token) {
	scopeLen := r.scopes.size()
	for i := scopeLen - 1; i >= 0; i-- {
		scope, _ := r.scopes.get(i)
		if _, ok := scope[name.lexeme]; ok {
			r.interpreter.resolve(e, scopeLen-1-i)
			return
		}
	}
}

func (r *Resolver) beginScope() {
	r.scopes.push(make(map[string]bool, 0))
}

func (r *Resolver) endScope() {
	if !r.scopes.isEmpty() {
		return
	}
	r.scopes.pop()
}

func (r *Resolver) declare(name token) {
	scope, err := r.scopes.peek()
	if err != nil {
		return
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name token) {
	scope, err := r.scopes.peek()
	if err != nil {
		return
	}
	scope[name.lexeme] = true
}

func (r *Resolver) visitBinaryExpr(e binaryExpr) (any, error) {
	// TODO: after logger, report error here
	r.resolveExpr(e.left)
	r.resolveExpr(e.right)
	return nil, nil
}

func (r *Resolver) visitCallExpr(e callExpr) (any, error) {
	r.resolveExpr(e.callee)
	for _, a := range e.arguments {
		r.resolveExpr(a)
	}
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(e groupingExpr) (any, error) {
	r.resolveExpr(e.expr)
	return nil, nil
}

func (r *Resolver) visitLiteralExpr(e literalExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(e logicalExpr) (any, error) {
	r.resolveExpr(e.left)
	r.resolveExpr(e.right)
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(e unaryExpr) (any, error) {
	r.resolveExpr(e.right)
	return nil, nil
}

func (r *Resolver) visitVariableExpr(e variableExpr) (any, error) {
	if currentScope, err := r.scopes.peek(); err != nil {
		exists, ok := currentScope[e.name.lexeme]
		if ok && !exists {
			// TODO: return nil, NewParseError(e.name, "Can't read local variable in its own initializer.")
			fmt.Printf("[line %d] Error: Can't read local variable in its own initializer.", e.name.line)
		}
	}
	r.resolveLocal(e, e.name)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(e assignExpr) (any, error) {
	r.resolveExpr(e.value)
	r.resolveLocal(e, e.name)
	return nil, nil
}

func (r *Resolver) visitExprStmt(s exprStmt) error {
	r.resolveExpr(s.expr)
	return nil
}

func (r *Resolver) visitFunctionStmt(s functionStmt) error {
	r.declare(s.name)
	r.define(s.name)
	r.resolveFunction(s)
	return nil
}

func (r *Resolver) visitIfStmt(s ifStmt) error {
	r.resolveExpr(s.condition)
	r.resolveStmt(s.thenBranch)
	if s.elseBranch != nil {
		r.resolveStmt(s.elseBranch)
	}
	return nil
}

func (r *Resolver) visitPrintStmt(s printStmt) error {
	r.resolveExpr(s.expr)
	return nil
}

func (r *Resolver) visitReturnStmt(s returnStmt) error {
	if s.value != nil {
		r.resolveExpr(s.value)
	}
	return nil
}

func (r *Resolver) visitVarStmt(s varStmt) error {
	r.declare(s.name)
	if s.initializer != nil {
		r.resolveExpr(s.initializer)
	}
	r.define(s.name)
	return nil
}

func (r *Resolver) visitWhileStmt(s whileStmt) error {
	r.resolveExpr(s.condition)
	r.resolveStmt(s.body)
	return nil
}

func (r *Resolver) visitBlockStmt(s blockStmt) error {
	r.beginScope()
	defer r.endScope()
	err := r.resolveStmtList(s.statements)
	if err != nil {
		return err
	}
	// NOTE: Explicitly return nil so that endScope() runs AFTER resolveStmtList()
	return nil
}
