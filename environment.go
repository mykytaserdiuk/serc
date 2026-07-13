package main

type Environment struct {
	values map[string]any
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Set(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name string) (any, bool) {
	value, ok := e.values[name]
	return value, ok
}
