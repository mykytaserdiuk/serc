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

type Block struct{
	Statements []Statement
}

type If struct {
	Conditions Binary
	Then Block
	Else Block
}

func (If) statement() {}
