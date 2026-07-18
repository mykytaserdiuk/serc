package core

import "github.com/mykytaserdiuk/serc/ast"

type Environment struct {
	values map[string]ast.Value
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]ast.Value),
	}
}

func (e *Environment) Set(name string, value ast.Value) {
	e.values[name] = value
}

func (e *Environment) Get(name string) (ast.Value, bool) {
	value, ok := e.values[name]
	return value, ok
}
