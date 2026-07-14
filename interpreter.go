package main

import (
	"fmt"
	"strconv"
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

func (i *Interpreter) execute(fn *Func) FuncResult {
	return i.exressionExecute(fn.body)
}

func (i *Interpreter) Main() {
	mainFn := i.findFunc("main")
	if mainFn == nil {
		fmt.Printf("ERROR: cant find main func")
		return
	}
	i.execute(mainFn)
}

func (i *Interpreter) exressionExecute(statements []Statement) FuncResult {
	for _, state := range statements {
		switch s := state.(type) {
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
				vals := make([]any, 0)
				for _, a := range s.args {
					vals = append(vals, i.eval(a).Data)
				}
				fmt.Println(vals...)
			default:
				exFn := i.findFunc(s.name)
				if exFn != nil {
					argsValues := make([]Value, len(s.args))

					for idx, argExpr := range s.args {
						argsValues[idx] = i.eval(argExpr)
					}

					for idx, argName := range exFn.params {
						if idx < len(argsValues) {
							i.env.Set(argName, argsValues[idx])
						}
					}
					return i.execute(exFn)
				}
			}
		case NewAssign:
			computedValue := i.eval(s.Value)
			i.env.Set(s.VarName, computedValue)
		case Assign:
			if _, ok := i.env.Get(s.VarName); !ok {
				panic("RUNTIME ERROR: " + s.VarName + " is not defined")
			}
			computedValue := i.eval(s.Value)
			//fmt.Printf("%T %+v", computedValue, computedValue)
			i.env.Set(s.VarName, computedValue)
		case Return:
			return FuncResult{
				Value: s.Value,
			}
		case If:
			conditions := i.eval(s.Conditions)
			binaryConditionResult := conditions.Data.(bool)
			//			if binaryConditionResult.Type != BoolValue {
			//				panic("RUNTIME ERROR: Unexpected binary condition result: " + string(binaryConditionResult.Type))
			//			}
			//res := binaryConditionResult.Data.(bool)
			if binaryConditionResult {
				i.exressionExecute(s.Then.Statements)
			} else if !binaryConditionResult && len(s.Then.Statements) > 0 {
				i.exressionExecute(s.Else.Statements)
			}
		}
	}
	return FuncResult{
		Value: nil,
	}
}

func (i *Interpreter) findFunc(name string) *Func {
	if fn, ok := i.functions[name]; ok {
		return fn
	}
	return nil
}

func (i *Interpreter) eval(expr Expression) Value {
	switch e := expr.(type) {
	case StringLiteral:
		return stringValue(e.value)
	case NumberLiteral:
		return intValue(e.value)
	case FuncCall:
		result := i.execute(i.findFunc(e.name))
		return i.eval(result)
	case nil:
		return nullValue()
	case Binary:
		return i.evalBinary(e)
	case FuncResult:
		return i.eval(e.Value)
	case Variable:
		val, ok := i.env.Get(e.name)
		if !ok {
			panic(fmt.Sprintf("Runtime Error: variable '%s' is not defined", e.name))
		}
		return val
	}

	return nullValue()
}

func (i *Interpreter) evalBinary(bin Binary) Value {
	leftValue := i.eval(bin.Left)
	rightVal := i.eval(bin.Right)
	switch leftValue.Type {
	case IntValue:
		lint := leftValue.Data.(int)
		rint, ok := rightVal.Data.(int)
		if !ok {
			panic("RUNTIME ERROR: expected 'int' at right part of binary")
		}
		switch bin.Op.type_ {
		case LessTokenType:
			return boolValue(lint < rint)
		case MoreTokenType:
			return boolValue(lint > rint)
		case EqLessTokenType:
			return boolValue(lint <= rint)
		case EqMoreTokenType:
			return boolValue(lint >= rint)
		case EqEqTokenType:
			return boolValue(lint == rint)
		case PlusTokenType:
			return intValue(lint + rint)
		case MinusTokenType:
			return intValue(lint - rint)
		}
	case StringValue:
		lStr := leftValue.Data.(string)
		rStr, ok := rightVal.Data.(string)
		if !ok {
			panic("RUNTIME ERROR: expected 'string' at left part of binary")
		}

		llen := len(lStr)
		rlen := len(rStr)
		switch bin.Op.type_ {
		case LessTokenType:
			return boolValue(llen < rlen)
		case MoreTokenType:
			return boolValue(llen > rlen)
		case EqLessTokenType:
			return boolValue(llen <= rlen)
		case EqMoreTokenType:
			return boolValue(llen >= rlen)
		case EqEqTokenType:
			return boolValue(llen == rlen)
		case PlusTokenType, MinusTokenType:
			panic("RUNTIME ERROR: unexpected operator to string: '" + bin.Op.type_ + "'")
		}
	}
	panic("Row " + strconv.Itoa(bin.Op.line) + ": RUNTIME ERROR: Unexpected type to calculate binary")
}

// macros
func (i *Interpreter) printf(format string, exps ...Expression) {
	values := make([]any, len(exps))
	for idx, e := range exps {
		values[idx] = i.eval(e).Data
	}
	fmt.Printf(format, values...)
}
