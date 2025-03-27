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
	fnTypeINITIALIZER fnType = "isInitializer"
)

type function struct {
	declaration   functionStmt
	closure       *environment
	isInitializer bool
}

func newFunction(stmt functionStmt, closure *environment, isInitializer bool) *function {
	return &function{
		declaration:   stmt,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *function) call(i *Interpreter, args []any) (any, error) {
	env := newEnvironment(f.closure)
	for idx, param := range f.declaration.params {
		env.define(param.lexeme, args[idx])
	}
	err := i.executeBlock(blockStmt{f.declaration.body}, env)
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
	return len(f.declaration.params)
}

func (f *function) bind(i *instance) *function {
	env := newEnvironment(f.closure)
	env.define("this", i)
	return newFunction(f.declaration, env, f.isInitializer)
}

func (f *function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}
