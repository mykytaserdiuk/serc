package ast

type Program struct {
	Functions map[string]*Func
	Structs   map[string]*Structure
	Imports   map[string]Import
}

type Func struct {
	Name   string
	Params []string
	Body   []Statement
}
