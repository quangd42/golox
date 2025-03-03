package lox

import "fmt"

type environment struct {
	values    map[string]any
	enclosing *environment
}

func NewEnvironment(enclosing *environment) *environment {
	return &environment{
		values:    make(map[string]any, 0),
		enclosing: enclosing,
	}
}

func NewGlobalEnvironment() *environment {
	return NewEnvironment(nil)
}

func (e *environment) define(varName string, value any) {
	e.values[varName] = value
}

func (e *environment) assign(name token, value any) error {
	if _, ok := e.values[name.lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}
		return NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}
	e.values[name.lexeme] = value
	return nil
}

func (e environment) get(name token) (any, error) {
	out, ok := e.values[name.lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}
	return out, nil
}
