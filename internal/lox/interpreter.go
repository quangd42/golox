package lox

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Interpreter struct {
	er      ErrorReporter
	globals *environment
	locals  map[expr]int
	env     *environment
}

func NewInterpreter(er ErrorReporter) *Interpreter {
	globals := newGlobalEnvironment()
	defineClockFn(globals)
	return &Interpreter{
		er:      er,
		globals: globals,
		locals:  make(map[expr]int, 0),
		env:     globals,
	}
}

func (i *Interpreter) Interpret(stmts []stmt) error {
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			var rtErr RuntimeError
			if errors.As(err, &rtErr) {
				i.er.RuntimeError(rtErr)
			} else {
				return err
			}
		}
	}
	return nil
}

func (i *Interpreter) resolve(e expr, depth int) {
	i.locals[e] = depth
}

func (i *Interpreter) lookUpVariable(name token, e expr) (any, error) {
	distance, ok := i.locals[e]
	if !ok {
		return i.globals.get(name)
	}
	return i.env.getAt(distance, name.lexeme)
}

func (i *Interpreter) evaluate(e expr) (any, error) {
	return e.accept(i)
}

func (i *Interpreter) visitLiteralExpr(e literalExpr) (any, error) {
	return e.value, nil
}

func (i *Interpreter) visitUnaryExpr(e unaryExpr) (any, error) {
	val, err := i.evaluate(e.right)
	if err != nil {
		return nil, err
	}
	switch e.operator.tokenType {
	case MINUS:
		numI, err := i.assertInt(val)
		if err == nil {
			return -numI, nil
		}
		numF, err := i.assertFloat(val)
		if err == nil {
			return -numF, nil
		}
		return nil, NewRuntimeError(e.operator, "Operand must be a number.")
	case BANG:
		return !i.isTruthy(val), nil
	default:
		return nil, NewRuntimeError(e.operator, "Undefined unary operator.")
	}
}

func (i *Interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) visitBinaryExpr(e binaryExpr) (any, error) {
	left, err := i.evaluate(e.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(e.right)
	if err != nil {
		return nil, err
	}
	numErr := NewRuntimeError(e.operator, "Operands must be numbers.")
	switch e.operator.tokenType {
	case SLASH:
		leftInt, rightInt, err := i.assertIntOperands(left, right)
		if err == nil {
			if rightInt == 0 {
				return nil, NewRuntimeError(e.operator, "Divisor must not be zero.")
			}
			return leftInt / rightInt, nil
		}
		leftFloat, rightFloat, err := i.assertFloatOperands(left, right)
		if err == nil {
			if rightFloat == 0 {
				return nil, NewRuntimeError(e.operator, "Divisor must not be zero.")
			}
			return leftFloat / rightFloat, nil
		}
		return nil, numErr
	case STAR:
		leftInt, rightInt, err := i.assertIntOperands(left, right)
		if err == nil {
			return leftInt * rightInt, nil
		}
		leftFloat, rightFloat, err := i.assertFloatOperands(left, right)
		if err == nil {
			return leftFloat * rightFloat, nil
		}
		return nil, numErr
	case MINUS:
		leftInt, rightInt, err := i.assertIntOperands(left, right)
		if err == nil {
			return leftInt - rightInt, nil
		}
		leftFloat, rightFloat, err := i.assertFloatOperands(left, right)
		if err == nil {
			return leftFloat - rightFloat, nil
		}
		return nil, numErr
	case PLUS:
		leftInt, rightInt, err := i.assertIntOperands(left, right)
		if err == nil {
			return leftInt + rightInt, nil
		}
		leftFloat, rightFloat, err := i.assertFloatOperands(left, right)
		if err == nil {
			return leftFloat + rightFloat, nil
		}
		leftStr, rightStr, err := i.assertStringOperands(left, right)
		if err == nil {
			return leftStr + rightStr, nil
		}
		return nil, NewRuntimeError(e.operator, "Operands must be either numbers or strings.")
	case GREATER:
		leftNum, rightNum, err := i.assertFloatOperands(left, right)
		if err != nil {
			return nil, numErr
		}
		return leftNum > rightNum, nil
	case GREATER_EQUAL:
		leftNum, rightNum, err := i.assertFloatOperands(left, right)
		if err != nil {
			return nil, numErr
		}
		return leftNum >= rightNum, nil
	case LESS:
		leftNum, rightNum, err := i.assertFloatOperands(left, right)
		if err != nil {
			return nil, numErr
		}
		return leftNum < rightNum, nil
	case LESS_EQUAL:
		leftNum, rightNum, err := i.assertFloatOperands(left, right)
		if err != nil {
			return nil, numErr
		}
		return leftNum <= rightNum, nil
	case BANG_EQUAL:
		leftNum, rightNum, err := i.assertFloatOperands(left, right)
		if err == nil {
			return leftNum != rightNum, nil
		}
		return left != right, nil
	case EQUAL_EQUAL:
		leftNum, rightNum, err := i.assertFloatOperands(left, right)
		if err == nil {
			return leftNum == rightNum, nil
		}
		return left == right, nil
	default:
		return nil, NewRuntimeError(e.operator, "Undefined binary operator.")
	}
}

func (i *Interpreter) assertFloat(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, errors.New("not a float")
	}
}

func (i *Interpreter) assertInt(val any) (int, error) {
	out, ok := val.(int)
	if !ok {
		return 0, errors.New("not an int")
	}
	return out, nil
}

func (i *Interpreter) assertString(val any) (string, error) {
	strVal, ok := val.(string)
	if ok {
		return strVal, nil
	}
	numVal, err := i.assertFloat(val)
	if err == nil {
		return strconv.FormatFloat(numVal, 'g', 'g', 64), nil
	}
	return "", errors.New("not a string")
}

func (i *Interpreter) assertIntOperands(left, right any) (leftNum, rightNum int, err error) {
	err = errors.New("operands are not int")
	leftNum, nErr := i.assertInt(left)
	if nErr != nil {
		return 0, 0, err
	}
	rightNum, nErr = i.assertInt(right)
	if nErr != nil {
		return 0, 0, err
	}
	return leftNum, rightNum, nil
}

func (i *Interpreter) assertFloatOperands(left, right any) (leftNum, rightNum float64, err error) {
	err = errors.New("operands are not float")
	leftNum, nErr := i.assertFloat(left)
	if nErr != nil {
		return 0, 0, err
	}
	rightNum, nErr = i.assertFloat(right)
	if nErr != nil {
		return 0, 0, err
	}
	return leftNum, rightNum, nil
}

func (i *Interpreter) assertStringOperands(left, right any) (leftStr, rightStr string, err error) {
	err = errors.New("operands cannot be interpreted as strings")
	leftStr, sErr := i.assertString(left)
	if sErr != nil {
		return "", "", err
	}
	rightStr, sErr = i.assertString(right)
	if sErr != nil {
		return "", "", err
	}
	return leftStr, rightStr, nil
}

func (i *Interpreter) visitGroupingExpr(e groupingExpr) (any, error) {
	return i.evaluate(e.expr)
}

func (i *Interpreter) visitVariableExpr(e variableExpr) (any, error) {
	return i.lookUpVariable(e.name, e)
}

func (i *Interpreter) visitAssignExpr(e assignExpr) (any, error) {
	val, err := i.evaluate(e.value)
	if err != nil {
		return nil, err
	}
	distance, ok := i.locals[e]
	if ok {
		err = i.env.assignAt(distance, e.name, val)
	} else {
		err = i.globals.assign(e.name, val)
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (i *Interpreter) visitLogicalExpr(e logicalExpr) (any, error) {
	leftVal, err := i.evaluate(e.left)
	if err != nil {
		return nil, err
	}
	// Short-circuit
	if e.operator.tokenType == OR {
		if i.isTruthy(leftVal) {
			return leftVal, nil
		}
	} else if !i.isTruthy(leftVal) {
		return leftVal, nil
	}

	return i.evaluate(e.right)
}

func (i *Interpreter) visitCallExpr(e callExpr) (any, error) {
	callee, err := i.evaluate(e.callee)
	if err != nil {
		return nil, err
	}
	args := make([]any, len(e.arguments))
	for idx, argExpr := range e.arguments {
		arg, err := i.evaluate(argExpr)
		if err != nil {
			return nil, err
		}
		args[idx] = arg
	}
	function, ok := callee.(callable)
	if !ok {
		return nil, NewRuntimeError(e.paren, "Can only call functions and classes.")
	}
	if function.arity() != len(args) {
		return nil, NewRuntimeError(
			e.paren,
			fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(args)),
		)
	}
	return function.call(i, args)
}

func (i *Interpreter) visitTernaryExpr(e ternaryExpr) (any, error) {
	condition, err := i.evaluate(e.condition)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(condition) {
		return i.evaluate(e.thenExpr)
	} else {
		return i.evaluate(e.elseExpr)
	}
}

func (i *Interpreter) visitGetExpr(e getExpr) (any, error) {
	object, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}
	instance, ok := object.(*instance)
	if !ok {
		return nil, NewRuntimeError(e.name, "Only instances have properties.")
	}
	val, ok := instance.fields[e.name.lexeme]
	if ok {
		return val, nil
	}
	method, ok := instance.class.findMethod(e.name.lexeme)
	if ok {
		return method.bind(instance), nil
	}
	return nil, NewRuntimeError(e.name, fmt.Sprintf("Undefined properties '%s'", e.name.lexeme))
}

func (i *Interpreter) visitSetExpr(e setExpr) (any, error) {
	val, err := i.evaluate(e.value)
	if err != nil {
		return nil, err
	}
	object, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}
	instance, ok := object.(*instance)
	if !ok {
		return nil, NewRuntimeError(e.name, "Only instances have fields.")
	}
	instance.fields[e.name.lexeme] = val
	return val, nil
}

func (i *Interpreter) visitThisExpr(e thisExpr) (any, error) {
	return i.lookUpVariable(e.keyword, e)
}

func (i *Interpreter) visitSuperExpr(e superExpr) (any, error) {
	distance, ok := i.locals[e]
	if !ok {
		return nil, errors.New("could not find super expr in locals")
	}
	val, err := i.env.getAt(distance, "super")
	if err != nil {
		return nil, err
	}
	superclass, ok := val.(*class)
	if !ok {
		return nil, errors.New("'super' value must be '*class'")
	}
	val, err = i.env.getAt(distance-1, "this")
	if err != nil {
		return nil, err
	}
	object, ok := val.(*instance)
	if !ok {
		return nil, errors.New("'this' value must be '*instance'")
	}
	method, ok := superclass.findMethod(e.method.lexeme)
	if !ok {
		return nil, NewRuntimeError(e.method, fmt.Sprintf("Undefined property '%s'.", e.method.lexeme))
	}
	return method.bind(object), nil
}

func (i *Interpreter) visitFunctionExpr(e functionExpr) (any, error) {
	return newAnonymousFunction(e, i.env), nil
}

func (i *Interpreter) visitArrayExpr(e arrayExpr) (any, error) {
	out := make([]any, len(e.value))
	for idx, X := range e.value {
		val, err := i.evaluate(X)
		if err != nil {
			return nil, err
		}
		out[idx] = val
	}
	return out, nil
}

func (i *Interpreter) visitIndexExpr(e indexExpr) (any, error) {
	callee, err := i.evaluate(e.callee)
	if err != nil {
		return nil, err
	}
	index, err := i.evaluate(e.index)
	if err != nil {
		return nil, err
	}
	array, ok := callee.([]any)
	if !ok {
		return nil, NewRuntimeError(e.bracket, "Can only index arrays.")
	}
	indexInt, ok := index.(int)
	if !ok {
		return nil, NewRuntimeError(e.bracket, "Index must be an integer.")
	}
	if indexInt < 0 {
		return nil, NewRuntimeError(e.bracket, fmt.Sprintf("Index out of range [%d].", indexInt))
	} else if len(array) == 0 || indexInt > len(array) {
		return nil, NewRuntimeError(e.bracket, fmt.Sprintf("Index out of range [%d] with length %d.", indexInt, len(array)))
	}
	return array[indexInt], nil
}

func (i *Interpreter) execute(s stmt) error {
	return s.accept(i)
}

func (i *Interpreter) executeBlock(s blockStmt, blockEnv *environment) error {
	prevEnv := i.env
	i.env = blockEnv
	defer func(i *Interpreter) {
		i.env = prevEnv
	}(i)
	for _, stmt := range s.statements {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) visitVarStmt(s varStmt) error {
	var val any
	var err error
	if s.initializer != nil {
		val, err = i.evaluate(s.initializer)
		if err != nil {
			return err
		}
	}
	i.env.define(s.name.lexeme, val)
	return nil
}

func (i *Interpreter) visitExprStmt(s exprStmt) error {
	_, err := i.evaluate(s.expr)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) visitFunctionStmt(s functionStmt) error {
	i.env.define(s.name.lexeme, newFunction(s.name, s.literal, i.env, false))
	return nil
}

func (i *Interpreter) visitIfStmt(s ifStmt) error {
	condition, err := i.evaluate(s.condition)
	if err != nil {
		return err
	}
	if i.isTruthy(condition) {
		return i.execute(s.thenBranch)
	} else if s.elseBranch != nil {
		return i.execute(s.elseBranch)
	} else {
		return nil
	}
}

func (i *Interpreter) visitPrintStmt(s printStmt) error {
	val, err := i.evaluate(s.expr)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", val)
	return nil
}

// Use error to exit execution early
func (i *Interpreter) visitReturnStmt(s returnStmt) error {
	var val any
	var err error
	if s.value != nil {
		val, err = i.evaluate(s.value)
		if err != nil {
			return err
		}
	}
	return &functionReturn{value: val}
}

func (i *Interpreter) visitWhileStmt(s whileStmt) error {
	for {
		condVal, err := i.evaluate(s.condition)
		if err != nil {
			return err
		}
		if !i.isTruthy(condVal) {
			return nil
		}
		err = i.execute(s.body)
		if err != nil {
			var breakErr *loopBreak
			if errors.As(err, &breakErr) {
				if breakErr.label.lexeme != "" && breakErr.label.lexeme != s.label.lexeme {
					return err
				}
				return nil
			}
			var contErr *loopContinue
			if errors.As(err, &contErr) {
				if contErr.label.lexeme != "" && contErr.label.lexeme != s.label.lexeme {
					return err
				}
				if s.increment != nil {
					err = i.execute(s.increment)
					if err != nil {
						return err
					}
				}
				continue
			}
			return err
		}
		if s.increment != nil {
			err = i.execute(s.increment)
			if err != nil {
				return err
			}
		}
	}
}

func (i *Interpreter) visitForStmt(s forStmt) error {
	return i.executeBlock(blockStmt{[]stmt{
		s.initializer,
		s.whileBody,
	}}, newEnvironment(i.env))
}

func (i *Interpreter) visitBreakStmt(s breakStmt) error {
	return &loopBreak{keyword: s.keyword, label: s.label}
}

func (i *Interpreter) visitContinueStmt(s continueStmt) error {
	return &loopContinue{keyword: s.keyword, label: s.label}
}

func (i *Interpreter) visitBlockStmt(s blockStmt) error {
	return i.executeBlock(s, newEnvironment(i.env))
}

func (i *Interpreter) visitClassStmt(s classStmt) error {
	var superclass *class
	var ok bool
	if s.superclass != (variableExpr{}) {
		val, err := i.evaluate(s.superclass)
		if err != nil {
			return err
		}
		superclass, ok = val.(*class)
		if !ok {
			return NewRuntimeError(s.superclass.name, "Superclass must be a class.")
		}
	}
	// two-stage variable binding process allows references to the class
	// inside its own methods
	i.env.define(s.name.lexeme, nil)
	if s.superclass != (variableExpr{}) {
		i.env = newEnvironment(i.env)
		i.env.define("super", superclass)
	}
	methods := make(map[string]*function, len(s.methods))
	for _, m := range s.methods {
		methods[m.name.lexeme] = newFunction(m.name, m.literal, i.env, m.name.lexeme == "init")
	}
	if s.superclass != (variableExpr{}) {
		i.env = i.env.enclosing
	}
	i.env.assign(s.name, newClass(s.name.lexeme, superclass, methods))
	return nil
}

func defineClockFn(env *environment) {
	env.define("clock", nativeFn{
		arityFn: func() int { return 0 },
		callFn: func(i *Interpreter, args []any) (any, error) {
			return time.Now().Unix(), nil
		},
		stringFn: func() string { return "<native fn clock>" },
	})
}
