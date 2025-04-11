package lox

import (
	"time"
)

type builtinFn struct {
	arityFn  func() int
	callFn   func(i *Interpreter, args []any) (any, error)
	stringFn func() string
}

func (f builtinFn) arity() int {
	return f.arityFn()
}

func (f builtinFn) call(i *Interpreter, args []any) (any, error) {
	return f.callFn(i, args)
}

func (f builtinFn) String() string {
	return f.stringFn()
}

func defineNativeFns(env *environment) {
	defineClockFn(env)
	defineArrayFns(env)
}

func defineClockFn(env *environment) {
	env.define("clock", builtinFn{
		arityFn: func() int { return 0 },
		callFn: func(i *Interpreter, args []any) (any, error) {
			return time.Now().Unix(), nil
		},
		stringFn: func() string { return "<native fn clock>" },
	})
}

func defineArrayFns(env *environment) {
	env.define("len", builtinFn{
		arityFn: func() int { return 1 },
		callFn: func(i *Interpreter, args []any) (any, error) {
			arr, ok := args[0].(*array)
			if !ok {
				return nil, builtinErrMsg("Can only call 'len' on arrays.")
			}
			return arr.Len(), nil
		},
		stringFn: func() string { return "<native fn len>" },
	})

	env.define("append", builtinFn{
		arityFn: func() int { return -1 },
		callFn: func(i *Interpreter, args []any) (any, error) {
			arr, ok := args[0].(*array)
			if !ok {
				return nil, builtinErrMsg("Can only call 'append' on arrays.")
			}
			arr.Append(args[1:]...)
			return nil, nil
		},
	})
}

type builtinErrMsg string

func (em builtinErrMsg) Error() string {
	return string(em)
}
