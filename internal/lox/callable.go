package lox

type callable interface {
	call(i *Interpreter, args []any) (any, error)
	arity() int
}

type nativeFn struct {
	arityFn  func() int
	callFn   func(i *Interpreter, args []any) (any, error)
	stringFn func() string
}

func (f nativeFn) arity() int {
	return f.arityFn()
}

func (f nativeFn) call(i *Interpreter, args []any) (any, error) {
	return f.callFn(i, args)
}

func (f nativeFn) String() string {
	return f.stringFn()
}
