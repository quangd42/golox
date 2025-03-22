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
	return i.env.getAt(distance, name)
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
		num, err := i.assertNumber(val)
		if err != nil {
			return nil, NewRuntimeError(e.operator, "Operand must be a number.")
		}
		return -num, nil
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
	switch e.operator.tokenType {
	case SLASH:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum / rightNum, nil
	case STAR:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum * rightNum, nil
	case MINUS:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum - rightNum, nil
	case PLUS:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err == nil {
			return leftNum + rightNum, nil
		}
		leftStr, rightStr, err := i.assertStringOperands(e.operator, left, right)
		if err == nil {
			return leftStr + rightStr, nil
		}
		return nil, NewRuntimeError(e.operator, "Operands must be either numbers or strings.")
	case GREATER:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum > rightNum, nil
	case GREATER_EQUAL:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum >= rightNum, nil
	case LESS:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum < rightNum, nil
	case LESS_EQUAL:
		leftNum, rightNum, err := i.assertNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum <= rightNum, nil
	case BANG_EQUAL:
		return left != right, nil
	case EQUAL_EQUAL:
		return left == right, nil
	default:
		return nil, NewRuntimeError(e.operator, "Undefined binary operator.")
	}
}

func (i *Interpreter) assertNumber(val any) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, errors.New("NaN")
	}
}

func (i *Interpreter) assertString(val any) (string, error) {
	strVal, ok := val.(string)
	if ok {
		return strVal, nil
	}
	numVal, err := i.assertNumber(val)
	if err == nil {
		return strconv.FormatFloat(numVal, 'g', 'g', 64), nil
	}
	return "", errors.New("not a string")
}

func (i *Interpreter) assertNumberOperands(operator token, left, right any) (leftNum, rightNum float64, err error) {
	err = NewRuntimeError(operator, "Operands must be numbers.")
	leftNum, nErr := i.assertNumber(left)
	if nErr != nil {
		return 0, 0, err
	}
	rightNum, nErr = i.assertNumber(right)
	if nErr != nil {
		return 0, 0, err
	}
	return leftNum, rightNum, nil
}

func (i *Interpreter) assertStringOperands(operator token, left, right any) (leftStr, rightStr string, err error) {
	err = NewRuntimeError(operator, "Operands must be strings.")
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

func (i *Interpreter) visitGetExpr(e getExpr) (any, error) {
	object, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}
	instance, ok := object.(instance)
	if !ok {
		return nil, NewRuntimeError(e.name, "Only instances have properties.")
	}
	val, ok := instance.fields[e.name.lexeme]
	if ok {
		return val, nil
	}
	method, ok := instance.class.methods[e.name.lexeme]
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
	instance, ok := object.(instance)
	if !ok {
		return nil, NewRuntimeError(e.name, "Only instances have fields.")
	}
	instance.fields[e.name.lexeme] = val
	return val, nil
}

func (i *Interpreter) visitThisExpr(e thisExpr) (any, error) {
	return i.lookUpVariable(e.keyword, e)
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
	i.env.define(s.name.lexeme, newFunction(s, i.env))
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
	}
	return nil
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
	return &returnValue{value: val}
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
			return err
		}
	}
}

func (i *Interpreter) visitBlockStmt(s blockStmt) error {
	return i.executeBlock(s, newEnvironment(i.env))
}

func (i *Interpreter) visitClassStmt(s classStmt) error {
	// two-stage variable binding process allows references to the class
	// inside its own methods
	i.env.define(s.name.lexeme, nil)
	methods := make(map[string]function, len(s.methods))
	for _, m := range s.methods {
		methods[m.name.lexeme] = newFunction(m, i.env)
	}
	i.env.assign(s.name, newClass(s.name.lexeme, methods))
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
