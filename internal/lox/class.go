package lox

type class struct {
	name string
}

func newClass(name string) class {
	return class{name}
}

func (c class) call(i *Interpreter, args []any) (any, error) {
	return newInstance(c), nil
}

func (c class) arity() int {
	return 0
}

func (c class) String() string {
	return c.name
}
