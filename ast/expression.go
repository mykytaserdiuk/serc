package ast

import (
	"fmt"
	"reflect"
)

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

type NativeObject struct {
	Name   string
	Data   any
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

type FunctionValue struct {
	Func *Func
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
	NativeObjectValue
	FuncValue
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

// func (NativeObject) expression()  {}

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

func GetNativeObjectValue(name string, data any) Value {
	fields := make(map[string]Value)
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	t := v.Type()
	for idx := 0; idx < v.NumField(); idx++ {
		fieldName := t.Field(idx).Name
		fieldValue := v.Field(idx)

		switch fieldValue.Kind() {
		case reflect.String:
			fields[fieldName] = GetStringValue(
				fieldValue.String(),
			)
		case reflect.Int:
			fields[fieldName] = GetIntValue(
				int(fieldValue.Int()),
			)
		case reflect.Bool:
			fields[fieldName] = GetBoolValue(
				fieldValue.Bool(),
			)
		default:
			fields[fieldName] = GetNullValue()
		}
	}
	return Value{
		Type: NativeObjectValue,
		Data: &NativeObject{
			Name:   name,
			Data:   data,
			Fields: fields,
		},
	}
}
