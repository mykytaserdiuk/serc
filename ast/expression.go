package ast

import "fmt"

type Expression interface {
	Node
	expression()
}

type StringLiteral struct {
	Value string
}
type NumberLiteral struct {
	Value int
}

type Variable struct {
	Name string
}

type Object struct {
	Type   *Structure
	Fields map[string]Value
}

func (o Object) String() string {
	result := o.Type.Name + "{"

	first := true

	for name, value := range o.Fields {
		if !first {
			result += ", "
		}

		result += name + ": " + fmt.Sprint(value.Data)

		first = false
	}

	result += "}"

	return result
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
	ObjectValue
)

type FuncResult struct {
	Value Expression
}

type FieldAccess struct {
	Name  string
	Value Expression
}

func (Variable) expression()      {}
func (NumberLiteral) expression() {}
func (StringLiteral) expression() {}
func (Binary) expression()        {}
func (FuncResult) expression()    {}
func (FieldAccess) expression()   {}

// //func (Variable) expression(){}

// get values
func GetIntValue(v int) Value {
	return Value{
		Type: IntValue,
		Data: v,
	}
}

func GetStringValue(v string) Value {
	return Value{
		Type: StringValue,
		Data: v,
	}
}

func GetBoolValue(v bool) Value {
	return Value{
		Type: BoolValue,
		Data: v,
	}
}

func GetNullValue() Value {
	return Value{
		Type: NullValue,
		Data: "null",
	}
}
