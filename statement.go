package main

type Statement interface {
	Node
	statement()
}

type Argument struct {
	Name  string
	Value Expression
}

type Call struct {
	name string
	args []Argument
}

func (Call) statement()  {}
func (Call) expression() {}

type NewAssign struct {
	VarName string
	Value   Expression
}

func (NewAssign) statement() {}

type Assign struct {
	Target Expression
	Value  Expression
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

type Structure struct {
	Name   string
	Fields []string
}

type StructureCall struct {
	Name  string
	Value Expression
}

func (StructureCall) statement() {}
