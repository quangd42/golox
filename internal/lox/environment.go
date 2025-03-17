package lox

import (
	"errors"
	"fmt"
)

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

func (e *environment) assignAt(distance int, name token, value any) error {
	env, err := e.getOuterEnvAt(distance)
	if err != nil {
		return err
	}
	env.values[name.lexeme] = value
	return nil
}

func (e *environment) get(name token) (any, error) {
	out, ok := e.values[name.lexeme]
	if ok {
		return out, nil
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.lexeme))
}

func (e *environment) getAt(distance int, name token) (any, error) {
	targetEnv, err := e.getOuterEnvAt(distance)
	if err != nil {
		return nil, err
	}
	out, ok := targetEnv.values[name.lexeme]
	if !ok {
		return nil, fmt.Errorf("incorrect distance %d from resolver: could not find variable '%s' on line %d", distance, name.lexeme, name.line)
	}
	return out, nil
}

func (e *environment) getOuterEnvAt(distance int) (*environment, error) {
	cursor := e
	for range distance {
		if e.enclosing == nil {
			return nil, errors.New("could not find outer env: distance passed from resolver is too big")
		}
		cursor = e.enclosing
	}
	return cursor, nil
}
