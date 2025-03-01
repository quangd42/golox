package lox

import "errors"

type interpreter struct{}

func NewInterpreter() *interpreter {
	return &interpreter{}
}

func (i interpreter) Interpret(e expr) (any, error) {
	return e.accept(i)
}

func (i interpreter) visitLiteralExpr(e literalExpr) (any, error) {
	return e.value, nil
}

func (i interpreter) visitUnaryExpr(e unaryExpr) (any, error) {
	val, err := e.right.accept(i)
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
		return nil, NewRuntimeError(e.operator, "Undefined Operator.")
	}
}

func (i interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func (i interpreter) visitBinaryExpr(e binaryExpr) (any, error) {
	left, err := e.left.accept(i)
	if err != nil {
		return nil, err
	}
	right, err := e.right.accept(i)
	if err != nil {
		return nil, err
	}
	switch e.operator.tokenType {
	case SLASH:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum / rightNum, nil
	case STAR:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum * rightNum, nil
	case MINUS:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum - rightNum, nil
	case PLUS:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err == nil {
			return leftNum + rightNum, nil
		}
		leftStr, rightStr, err := i.checkStringOperands(e.operator, left, right)
		if err == nil {
			return leftStr + rightStr, nil
		}
		return nil, NewRuntimeError(e.operator, "Operands must be either numbers or strings.")
	case GREATER:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum > rightNum, nil
	case GREATER_EQUAL:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum >= rightNum, nil
	case LESS:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum < rightNum, nil
	case LESS_EQUAL:
		leftNum, rightNum, err := i.checkNumberOperands(e.operator, left, right)
		if err != nil {
			return nil, err
		}
		return leftNum <= rightNum, nil
	case BANG_EQUAL:
		return left != right, nil
	case EQUAL_EQUAL:
		return left == right, nil
	default:
		return nil, NewRuntimeError(e.operator, "Undefined Operator.")
	}
}

func (i interpreter) assertNumber(val any) (float64, error) {
	out, ok := val.(float64)
	if ok {
		return out, nil
	}
	outInt, ok := val.(int)
	if ok {
		return float64(outInt), nil
	}
	return 0, errors.New("NaN")
}

func (i interpreter) checkNumberOperands(operator token, left, right any) (leftNum, rightNum float64, err error) {
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

func (i interpreter) checkStringOperands(operator token, left, right any) (leftStr, rightStr string, err error) {
	err = NewRuntimeError(operator, "Operands must be strings.")
	leftStr, ok := left.(string)
	if !ok {
		return "", "", err
	}
	rightStr, ok = right.(string)
	if !ok {
		return "", "", err
	}
	return leftStr, rightStr, nil
}

func (i interpreter) visitGroupingExpr(e groupingExpr) (any, error) {
	return e.expr.accept(i)
}
