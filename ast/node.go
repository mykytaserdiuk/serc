package ast

type Node interface{}

type BuiltinFunc func(args []Value) FuncResult
