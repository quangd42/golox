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
	err := r.resolveStmtList(stmts)
	if err != nil {
		fmt.Print(err)
	}
	return nil
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
	err := r.resolveStmtList(s.body)
	if err != nil {
		return err
	}
	// NOTE: Explicitly return nil so that endScope() runs AFTER resolveStmtList()
	return nil
}

func (r *Resolver) resolveLocal(e expr, name token) {
	scopeLen := r.scopes.size()
	for i := range scopeLen {
		scope, _ := r.scopes.get(i)
		if _, ok := scope[name.lexeme]; ok {
			r.interpreter.resolve(e, i)
			return
		}
	}
}

func (r *Resolver) beginScope() {
	r.scopes.push(make(map[string]bool, 0))
}

func (r *Resolver) endScope() {
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
	_, err := r.resolveExpr(e.left)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(e.right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) visitCallExpr(e callExpr) (any, error) {
	_, err := r.resolveExpr(e.callee)
	if err != nil {
		return nil, err
	}
	for _, a := range e.arguments {
		_, err := r.resolveExpr(a)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) visitGroupingExpr(e groupingExpr) (any, error) {
	return r.resolveExpr(e.expr)
}

func (r *Resolver) visitLiteralExpr(e literalExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(e logicalExpr) (any, error) {
	_, err := r.resolveExpr(e.left)
	if err != nil {
		return nil, err
	}
	_, err = r.resolveExpr(e.right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) visitUnaryExpr(e unaryExpr) (any, error) {
	return r.resolveExpr(e.right)
}

func (r *Resolver) visitVariableExpr(e variableExpr) (any, error) {
	if currentScope, err := r.scopes.peek(); err != nil {
		exists, ok := currentScope[e.name.lexeme]
		if ok && !exists {
			return nil, NewParseError(e.name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(e, e.name)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(e assignExpr) (any, error) {
	_, err := r.resolveExpr(e.value)
	if err != nil {
		return nil, err
	}
	r.resolveLocal(e, e.name)
	return nil, nil
}

func (r *Resolver) visitExprStmt(s exprStmt) error {
	_, err := r.resolveExpr(s.expr)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) visitFunctionStmt(s functionStmt) error {
	r.declare(s.name)
	r.define(s.name)
	return r.resolveFunction(s)
}

func (r *Resolver) visitIfStmt(s ifStmt) error {
	_, err := r.resolveExpr(s.condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(s.thenBranch)
	if err != nil {
		return err
	}
	if s.elseBranch != nil {
		err = r.resolveStmt(s.elseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) visitPrintStmt(s printStmt) error {
	_, err := r.resolveExpr(s.expr)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) visitReturnStmt(s returnStmt) error {
	if s.value != nil {
		_, err := r.resolveExpr(s.value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) visitVarStmt(s varStmt) error {
	r.declare(s.name)
	defer r.define(s.name)
	if s.initializer != nil {
		_, err := r.resolveExpr(s.initializer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) visitWhileStmt(s whileStmt) error {
	_, err := r.resolveExpr(s.condition)
	if err != nil {
		return err
	}
	return r.resolveStmt(s.body)
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
