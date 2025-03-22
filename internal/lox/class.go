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
	instance := newInstance(c)
	if initializer, ok := c.methods["init"]; ok {
		// discard returned values from initializer when creating new a instance
		initializer.bind(instance).call(i, args)
	}
	return instance, nil
}

func (c class) arity() int {
	if initializer, ok := c.methods["init"]; ok {
		return initializer.arity()
	}
	return 0
}

func (c class) String() string {
	return c.name
}
