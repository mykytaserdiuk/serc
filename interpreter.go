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
	for _, b := range fn.body {
		//		fmt.Printf("kind: %s \n", reflect.TypeOf(b))
		switch s:= b.(type) {
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
		}
	}
}

func (i *Interpreter) Main() {
	mainFn := i.findFunc("main")
	if mainFn == nil {
		fmt.Printf("ERROR: cant find main func")
		return
	}
	i.execute(mainFn)
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


// macros

func (i *Interpreter) printf(format string, exps ...Expression) {
	values := make([]any, len(exps))
	for idx, e := range exps{
		values[idx] = i.eval(e)
	}
	fmt.Printf(format, values...)
}
