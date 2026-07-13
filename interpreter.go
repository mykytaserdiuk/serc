package main

import (
	"fmt"
)

type Interpreter struct {
	functions map[string]*Func
	env       *Environment
}

func NewInterpreter(parser *Parser) *Interpreter {
	funcs := make(map[string]*Func)
	for {
		if fn := parser.parseFunc(); fn != nil {
			if _, ok := funcs[fn.name]; !ok {
				funcs[fn.name] = fn
			}
		} else {
			break
		}
	}
	return &Interpreter{
		functions: funcs,
		env:       NewEnvironment(),
	}
}

func (i *Interpreter) execute(fn *Func) {
	i.exressionExecute(fn.body)
}

func (i *Interpreter) Main() {
	mainFn := i.findFunc("main")
	if mainFn == nil {
		fmt.Printf("ERROR: cant find main func")
		return
	}
	i.execute(mainFn)
}

func (i *Interpreter) exressionExecute(statements []Statement) {
	for _, state := range statements{
		switch s:= state.(type) {
		case FuncCall:
			switch s.name {
			case "printf":
				format, ok := s.args[0].(StringLiteral)
				if !ok {
					panic("RUNTIME ERROR: format isnt string")
				}
				args := s.args[1:]
				i.printf(format.value, args...)
			case "print":
				vals := make([]any,0)
				for _, a := range s.args{
					vals = append(vals, i.eval(a))
				}
				fmt.Println(vals...)
			default:
				exFn := i.findFunc(s.name)
				if exFn != nil {
					argsValues := make([]any, len(s.args))

					for idx, argExpr := range s.args {
						argsValues[idx] = i.eval(argExpr)
					}

					for idx, argName := range exFn.params {
						if idx < len(argsValues) {
							i.env.Set(argName, argsValues[idx])
						}
					}
					i.execute(exFn)
				}
			}
		case Assignment:
			computedValue := i.eval(s.Value)
			i.env.Set(s.VarName, computedValue)
		case If:
			binaryConditionResult := i.calculateBinary(s.Conditions)
			if binaryConditionResult{
				i.exressionExecute(s.Then.Statements)
			} else if !binaryConditionResult && len(s.Then.Statements) > 0 {
				i.exressionExecute(s.Else.Statements)
			}
		}
	}
}

func (i *Interpreter) findFunc(name string) *Func {
	if fn, ok := i.functions[name]; ok {
		return fn
	}
	return nil
}

func (i *Interpreter) eval(expr Expression) any {
	switch e := expr.(type) {
	case StringLiteral:
		return e.value
	case NumberLiteral:
		return e.value
	case Variable:
		val, ok := i.env.Get(e.name)
		if !ok {
			panic(fmt.Sprintf("Runtime Error: variable '%s' is not defined", e.name))
		}
		return val
	}

	return nil
}

func (i *Interpreter) calculateBinary(bin Binary) bool {
	lint := i.eval(bin.Left)
	rightVal := i.eval(bin.Right)
	switch v := lint.(type){
		case int:
		rint, ok := rightVal.(int)
		if !ok{
			panic("RUNTIME ERROR: expected 'int' at right part of binary")
		}
		switch bin.Op.type_{
			case LessTokenType:
			return  v < rint
			case MoreTokenType:
			return v>rint
			case EqLessTokenType:
			return v<=rint
			case EqMoreTokenType:
			return v>=rint
		}
		case string:
		rStr, ok := rightVal.(string)
		if !ok{
			panic("RUNTIME ERROR: expected 'string' at right part of binary")
		}
		llen := len(v)
		rlen := len(rStr)
		switch bin.Op.type_{
			case LessTokenType:
			return  llen < rlen
			case MoreTokenType:
			return llen>rlen
			case EqLessTokenType:
			return llen<=rlen
			case EqMoreTokenType:
			return llen>=rlen
		}
	}
	return false
}

// macros
func (i *Interpreter) printf(format string, exps ...Expression) {
	values := make([]any, len(exps))
	for idx, e := range exps{
		values[idx] = i.eval(e)
	}
	fmt.Printf(format, values...)
}
