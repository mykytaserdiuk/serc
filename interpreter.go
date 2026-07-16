package main

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	functions  map[string]*Func
	structures map[string]*Structure
	env        *Environment
}

func NewInterpreter(parser *Parser) *Interpreter {
	parsedProgram := parser.parseProgram()
	return &Interpreter{
		functions:  parsedProgram.Functions,
		structures: parsedProgram.Structs,
		env:        NewEnvironment(),
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
		case Call:
			switch s.name {
			case "printf":
				format, ok := s.args[0].Value.(StringLiteral)
				if !ok {
					panic("RUNTIME ERROR: format isnt string")
				}
				args := s.args[1:]
				i.printf(format.value, args...)
			case "print":
				vals := make([]any, 0)
				for _, a := range s.args {
					vals = append(vals, i.eval(a.Value).Data)
				}
				fmt.Print(vals...)
			default:
				exFn := i.findFunc(s.name)
				if exFn != nil {
					argsValues := make([]Value, len(s.args))

					for idx, argExpr := range s.args {
						argsValues[idx] = i.eval(argExpr.Value)
					}

					for idx, argName := range exFn.params {
						if idx < len(argsValues) {
							i.env.Set(argName, argsValues[idx])
						}
					}

					result := i.execute(exFn)
					if result.Value != nil {
						return result
					}
				}
			}
		case NewAssign:
			computedValue := i.eval(s.Value)
			i.env.Set(s.VarName, computedValue)
		case Assign:
			value := i.eval(s.Value)
			switch t := s.Target.(type) {
			case Variable:
				fmt.Println(s.Value)
				v, ok := i.env.Get(t.name)
				if !ok {
					panic("RUNTIME ERROR: variable '" + t.name + "' not found")
				}
				v.Data = value
				i.env.Set(t.name, v)
			case FieldAccess:
				obj := i.eval(t.Value)

				switch data := obj.Data.(type) {
				case Object:
					_, ok := data.Fields[t.Name]
					if !ok {
						panic("RUNTIME ERROR: field '" + t.Name + "' not found")
					}
					data.Fields[t.Name] = value
				default:
					panic("RUNTIME ERROR: object has no fields")
				}
			}
		case Return:
			return FuncResult{
				Value: s.Value,
			}
		case If:
			conditions := i.eval(s.Conditions)
			binaryConditionResult := conditions.Data.(bool)
			if binaryConditionResult {
				i.exressionExecute(s.Then.Statements)
			} else if !binaryConditionResult && len(s.Then.Statements) > 0 {
				i.exressionExecute(s.Else.Statements)
			}
		case StructureCall:
			fmt.Printf(s.Name, s.Value)
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
	case Call:
		return i.evalCall(e)
	case nil:
		return nullValue()
	case Binary:
		return i.evalBinary(e)
	case FuncResult:
		return i.eval(e.Value)
	case FieldAccess:
		obj := i.eval(e.Value)
		name := e.Name
		switch v := obj.Data.(type) {
		case Object:
			field, ok := v.Fields[name]
			if !ok {
				panic("field not found: " + name)
			}

			return field

		default:
			panic("cannot access field " + name)
		}
	case Variable:
		val, ok := i.env.Get(e.name)
		if !ok {
			panic(fmt.Sprintf("Runtime Error: variable '%s' is not defined", e.name))
		}
		return val
	}

	return nullValue()
}

func (i *Interpreter) evalCall(c Call) Value {
	st, ok := i.structures[c.name]
	if ok {
		return i.createStruct(st, c.args)
	}

	fn := i.findFunc(c.name)
	if fn != nil {
		args := make([]Value, len(c.args))

		for idx, arg := range c.args {
			args[idx] = i.eval(arg.Value)
		}

		oldEnv := i.env
		i.env = NewEnvironment()
		for idx, param := range fn.params {
			if idx < len(args) {
				i.env.Set(param, args[idx])
			}
		}

		res := i.execute(fn)
		i.env = oldEnv
		if res.Value == nil {
			return nullValue()
		}
		return i.eval(res.Value)
	}
	panic("unknown call: " + c.name)
}
func (i *Interpreter) createStruct(str *Structure, args []Argument) Value {
	obj := Object{
		Type:   str,
		Fields: make(map[string]Value),
	}

	for _, arg := range args {
		found := false
		if arg.Name == "" {
			panic("struct fields must be named")
		}
		for _, sArg := range str.Fields {
			if arg.Name == sArg {
				found = true
			}
		}
		if !found {
			panic("unknown field " + arg.Name +
				" in struct " + str.Name)
		}
		obj.Fields[arg.Name] = i.eval(arg.Value)

	}
	return Value{
		Type: ObjectValue,
		Data: obj,
	}
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
		case NotEqTokenType:
			return boolValue(lint != rint)
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
			return boolValue(lStr == rStr)
		case NotEqTokenType:
			return boolValue(lStr != rStr)
		case PlusTokenType, MinusTokenType:
			panic("RUNTIME ERROR: unexpected operator to string: '" + bin.Op.type_ + "'")
		}
	}
	panic("Row " + strconv.Itoa(bin.Op.line) + ": RUNTIME ERROR: Unexpected type to calculate binary")
}

// macros
func (i *Interpreter) printf(format string, exps ...Argument) {
	values := make([]any, len(exps))
	for idx, e := range exps {
		values[idx] = i.eval(e.Value).Data
	}
	fmt.Printf(format, values...)
}
