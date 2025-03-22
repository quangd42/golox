package lox

import (
	"errors"
	"fmt"
)

type fnType string

const (
	NONE     fnType = "none"
	FUNCTION fnType = "function"
	METHOD   fnType = "method"
)

type function struct {
	declaration functionStmt
	closure     *environment
}

func newFunction(stmt functionStmt, closure *environment) function {
	return function{
		declaration: stmt,
		closure:     closure,
	}
}

func (f function) call(i *Interpreter, args []any) (any, error) {
	env := newEnvironment(f.closure)
	for idx, param := range f.declaration.params {
		env.define(param.lexeme, args[idx])
	}
	err := i.executeBlock(blockStmt{f.declaration.body}, env)
	if err != nil {
		var retVal *returnValue
		if errors.As(err, &retVal) {
			return retVal.value, nil
		}
		return nil, err
	}
	return nil, nil
}

func (f function) arity() int {
	return len(f.declaration.params)
}

func (f function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}

type returnValue struct {
	value any
}

func (v *returnValue) Error() string {
	return fmt.Sprintf("%s", v.value)
}
