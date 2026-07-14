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

type Binary struct {
	Left  Expression
	Op    *Token
	Right Expression
}

type Value struct {
	Type ValueType
	Data any
}

type ValueType int

const (
	IntValue ValueType = iota
	StringValue
	BoolValue
	NullValue
)

type FuncResult struct {
	Value Expression
}

func (Variable) expression()      {}
func (NumberLiteral) expression() {}
func (StringLiteral) expression() {}
func (Binary) expression()        {}
func (FuncResult) expression()    {}

// //func (Variable) expression(){}

// get values
func intValue(v int) Value {
	return Value{
		Type: IntValue,
		Data: v,
	}
}

func stringValue(v string) Value {
	return Value{
		Type: StringValue,
		Data: v,
	}
}

func boolValue(v bool) Value {
	return Value{
		Type: BoolValue,
		Data: v,
	}
}

func nullValue() Value {
	return Value{
		Type: NullValue,
		Data: "null",
	}
}
