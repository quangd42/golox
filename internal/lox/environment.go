package lox

import "fmt"

type environment struct {
	values map[string]any
}

func NewEnvironment() *environment {
	return &environment{values: make(map[string]any, 0)}
}

func (e *environment) define(varName string, value any) {
	e.values[varName] = value
}

func (e *environment) assign(name token, value any) (any, error) {
	if _, ok := e.values[name.lexeme]; !ok {
		return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}
	e.values[name.lexeme] = value
	return e.values[name.lexeme], nil
}

func (e environment) get(name token) (any, error) {
	out, ok := e.values[name.lexeme]
	if !ok {
		return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
	}
	return out, nil
}
