package lox

import "fmt"

type instance struct {
	class class
}

func newInstance(c class) instance {
	return instance{c}
}

func (i instance) String() string {
	return fmt.Sprintf("%s instance", i.class.name)
}
