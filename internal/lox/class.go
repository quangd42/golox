package lox

type classType string

const (
	classTypeNONE  classType = "none"
	classTypeCLASS classType = "class"
)

type class struct {
	name    string
	methods map[string]function
}

func newClass(name string, methods map[string]function) class {
	return class{
		name:    name,
		methods: methods,
	}
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
