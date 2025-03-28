package lox

import (
	"errors"
	"fmt"
)

type fnType string

const (
	fnTypeNONE        fnType = "none"
	fnTypeFUNCTION    fnType = "function"
	fnTypeMETHOD      fnType = "method"
	fnTypeINITIALIZER fnType = "initializer"
	fnTypeANONYMOUS   fnType = "anonymous"
)

type function struct {
	name          token
	literal       functionExpr
	closure       *environment
	isInitializer bool
}

func newFunction(name token, literal functionExpr, closure *environment, isInitializer bool) *function {
	return &function{
		name:          name,
		literal:       literal,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func newAnonymousFunction(literal functionExpr, closure *environment) *function {
	return &function{
		literal: literal,
		closure: closure,
	}
}

func (f *function) call(i *Interpreter, args []any) (any, error) {
	env := newEnvironment(f.closure)
	for idx, param := range f.literal.params {
		env.define(param.lexeme, args[idx])
	}
	err := i.executeBlock(blockStmt{f.literal.body}, env)
	if err != nil {
		var fnRet *functionReturn
		if errors.As(err, &fnRet) {
			if f.isInitializer { // Returned value in initializer is overiden to 'this'
				return f.closure.getAt(0, "this")
			}
			return fnRet.value, nil
		}
		return nil, err
	}
	if f.isInitializer { // Early return in initializer should return 'this'
		return f.closure.getAt(0, "this")
	}
	return nil, nil
}

func (f *function) arity() int {
	return len(f.literal.params)
}

func (f *function) bind(i *instance) *function {
	env := newEnvironment(f.closure)
	env.define("this", i)
	return newFunction(f.name, f.literal, env, f.isInitializer)
}

func (f *function) String() string {
	if f.name.lexeme == "" {
		return "<anonymous fn>"
	}
	return fmt.Sprintf("<fn %s>", f.name.lexeme)
}
