package lox

type callable interface {
	call(i *Interpreter, args []any) (any, error)
	arity() int
}
