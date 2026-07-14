package main

type Environment struct {
	values map[string]Value
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]Value),
	}
}

func (e *Environment) Set(name string, value Value) {
	e.values[name] = value
}

func (e *Environment) Get(name string) (Value, bool) {
	value, ok := e.values[name]
	return value, ok
}
