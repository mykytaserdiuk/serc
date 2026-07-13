package main

type StringLiteral struct {
	value string
}
type NumberLiteral struct {
	value int
}

type Variable struct {
	name string
}

func (Variable) expression()      {}
func (NumberLiteral) expression() {}
func (StringLiteral) expression() {}

//func (Variable) expression(){}
