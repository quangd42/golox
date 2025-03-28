package lox

type loopType bool

type Resolver struct {
	errorReporter ErrorReporter
	interpreter   *Interpreter
	scopes        *scopeStack
	currentFn     fnType
	currentClass  classType
	loopStack     *loopStack
}

func NewResolver(er ErrorReporter, i *Interpreter) *Resolver {
	return &Resolver{
		errorReporter: er,
		interpreter:   i,
		scopes:        newScopeStack(),
		currentFn:     fnTypeNONE,
		currentClass:  classTypeNONE,
		loopStack:     newLoopStack(),
	}
}

func (r *Resolver) Resolve(stmts []stmt) error {
	return r.resolveStmtList(stmts)
}

func (r *Resolver) resolveStmtList(stmts []stmt) error {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
	return nil
}

func (r *Resolver) resolveExpr(e expr) (any, error) {
	return e.accept(r)
}

func (r *Resolver) resolveStmt(s stmt) error {
	return s.accept(r)
}

func (r *Resolver) resolveFunction(s functionStmt, ft fnType) error {
	enclosingFn := r.currentFn
	r.currentFn = ft
	defer func(r *Resolver) {
		r.currentFn = enclosingFn
	}(r)
	r.beginScope()
	defer r.endScope()
	for _, param := range s.params {
		r.declare(param)
		r.define(param)
	}
	return r.resolveStmtList(s.body)
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

func (r *Resolver) beginLoop(label string) {
	r.loopStack.push(label)
}

func (r *Resolver) endLoop() {
	r.loopStack.pop()
}

func (r *Resolver) visitBinaryExpr(e binaryExpr) (any, error) {
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
	return r.resolveExpr(e.expr)
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
	return r.resolveExpr(e.right)
}

func (r *Resolver) visitVariableExpr(e variableExpr) (any, error) {
	if currentScope, err := r.scopes.peek(); err == nil {
		varInitialized, ok := currentScope[e.name.lexeme]
		if ok && !varInitialized {
			r.errorReporter.ParseError(e.name, "Can't read local variable in its own initializer.")
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

func (r *Resolver) visitTernaryExpr(e ternaryExpr) (any, error) {
	r.resolveExpr(e.condition)
	r.resolveExpr(e.thenExpr)
	r.resolveExpr(e.elseExpr)
	return nil, nil
}

func (r *Resolver) visitGetExpr(e getExpr) (any, error) {
	r.resolveExpr(e.object)
	return nil, nil
}

func (r *Resolver) visitSetExpr(e setExpr) (any, error) {
	r.resolveExpr(e.value)
	r.resolveExpr(e.object)
	return nil, nil
}

func (r *Resolver) visitThisExpr(e thisExpr) (any, error) {
	if r.currentClass == classTypeNONE {
		r.errorReporter.ParseError(e.keyword, "Can't use 'this' outside of a class.")
	}
	r.resolveLocal(e, e.keyword)
	return nil, nil
}

func (r *Resolver) visitExprStmt(s exprStmt) error {
	r.resolveExpr(s.expr)
	return nil
}

func (r *Resolver) visitFunctionStmt(s functionStmt) error {
	r.declare(s.name)
	r.define(s.name)
	return r.resolveFunction(s, fnTypeFUNCTION)
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
	if r.currentFn == fnTypeNONE {
		r.errorReporter.ParseError(s.keyword, "Can't return from top-level code.")
	}
	if s.value != nil {
		if r.currentFn == fnTypeINITIALIZER {
			r.errorReporter.ParseError(s.keyword, "Can't return value from an initializer.")
		}
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
	if s.label.lexeme != "" && r.loopStack.contains(s.label.lexeme) {
		r.errorReporter.ParseError(s.label, "Label already belongs to outer loops.")
	}
	r.beginLoop(s.label.lexeme)
	defer r.endLoop()
	r.resolveExpr(s.condition)
	r.resolveStmt(s.body)
	if s.increment != nil {
		r.resolveStmt(s.increment)
	}
	return nil
}

func (r *Resolver) visitForStmt(s forStmt) error {
	r.beginScope()
	defer r.endScope()
	if s.initializer != nil {
		r.resolveStmt(s.initializer)
	}
	r.resolveStmt(s.whileBody)
	return nil
}

func (r *Resolver) visitBreakStmt(s breakStmt) error {
	if r.loopStack.isEmpty() {
		r.errorReporter.ParseError(s.keyword, "Break statement must be in a loop.")
	}
	if s.label.lexeme != "" && !r.loopStack.contains(s.label.lexeme) {
		r.errorReporter.ParseError(s.label, "Invalid break label.")
	}
	return nil
}

func (r *Resolver) visitContinueStmt(s continueStmt) error {
	if r.loopStack.isEmpty() {
		r.errorReporter.ParseError(s.keyword, "Continue statement must be in a loop.")
	}
	if s.label.lexeme != "" && !r.loopStack.contains(s.label.lexeme) {
		r.errorReporter.ParseError(s.label, "Invalid continue label.")
	}
	return nil
}

func (r *Resolver) visitBlockStmt(s blockStmt) error {
	r.beginScope()
	defer r.endScope()
	return r.resolveStmtList(s.statements)
}

func (r *Resolver) visitClassStmt(s classStmt) error {
	enclosingClass := r.currentClass
	r.currentClass = classTypeCLASS
	defer func(r *Resolver) {
		r.currentClass = enclosingClass
	}(r)

	r.declare(s.name)
	r.define(s.name)
	r.beginScope()
	defer r.endScope()
	currentScope, _ := r.scopes.peek() // after begining a scope this cannot fail
	currentScope["this"] = true
	for _, method := range s.methods {
		methodType := fnTypeMETHOD
		if method.name.lexeme == "init" {
			methodType = fnTypeINITIALIZER
		}
		r.resolveFunction(method, methodType)
	}
	return nil
}
