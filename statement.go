package main

type FuncCall struct {
	name string
	args []Expression
}

func (FuncCall) statement()  {}
func (FuncCall) expression() {}

type NewAssign struct {
	VarName string
	Value   Expression
}

func (NewAssign) statement() {}

type Assign struct {
	VarName string
	Value   Expression
}

func (Assign) statement() {}

type Block struct {
	Statements []Statement
}

type If struct {
	Conditions Expression
	Then       Block
	Else       Block
}

func (If) statement() {}

type Return struct {
	Value Expression
}

func (Return) statement() {}
