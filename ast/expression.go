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

// BE CAREFUL, its reflect
func GetNativeObjectValue(name string, data any) Value {
	fields := make(map[string]Value)

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {

	case reflect.Struct:
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldName := t.Field(i).Name

			fields[fieldName] = reflectToValue(field)
		}

	case reflect.Map:
		iter := v.MapRange()

		for iter.Next() {
			key := iter.Key()
			value := iter.Value()

			if key.Kind() != reflect.String {
				continue
			}

			fields[key.String()] = reflectToValue(value)
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

func reflectToValue(v reflect.Value) Value {
	switch v.Kind() {
	case reflect.String:
		return GetStringValue(v.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return GetIntValue(int(v.Int()))

	case reflect.Bool:
		return GetBoolValue(v.Bool())

	case reflect.Float32, reflect.Float64:
		return GetIntValue(int(v.Float())) // если есть

	case reflect.Map:
		obj := make(map[string]Value)

		iter := v.MapRange()
		for iter.Next() {
			if iter.Key().Kind() == reflect.String {
				obj[iter.Key().String()] = reflectToValue(iter.Value())
			}
		}

		return Value{
			Type: NativeObjectValue,
			Data: &NativeObject{
				Fields: obj,
			},
		}

	default:
		return GetNullValue()
	}
}
