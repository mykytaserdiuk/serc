package main

type FuncCall struct {
	name string
	args []Expression
}

func (FuncCall) statement() {}

type Assignment struct {
	VarName string
	Value   Expression
}

func (Assignment) statement() {}
