package ast

type Statement interface {
	Node
	statement()
}

type Argument struct {
	Name  string
	Value Expression
}

type Call struct {
	Target Expression
	Args   []Argument
}

func (Call) statement()  {}
func (Call) expression() {}
func (c Call) Name() string {
	return expressionName(c.Target)
}

func expressionName(expr Expression) string {
	switch t := expr.(type) {
	case Variable:
		return t.Name
	case FieldAccess:
		left := expressionName(t.Value)
		if left == "" {
			return t.Name
		}
		return left + "." + t.Name
	}

	return ""
}

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

type Import struct {
	Name  string
	Alias string
}

func (Import) statement() {}

type Loop struct {
	Conditions Expression
	Body       Block
}

func (Loop) statement() {}
