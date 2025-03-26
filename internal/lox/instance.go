package lox

import "fmt"

type instance struct {
	class  *class
	fields map[string]any
}

func newInstance(c *class) *instance {
	return &instance{
		class:  c,
		fields: make(map[string]any),
	}
}

func (i *instance) String() string {
	return fmt.Sprintf("%s instance", i.class.name)
}
