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

type Binary struct{
	Left Expression
	Op *Token
	Right Expression
}

type FuncResult struct{
	Value Expression
}

func (Variable) expression()      {}
func (NumberLiteral) expression() {}
func (StringLiteral) expression() {}
func (Binary) expression(){}
func (FuncResult) expression(){}

// //func (Variable) expression(){}
