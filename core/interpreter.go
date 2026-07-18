package core

import (
	"fmt"
	"strconv"

	"github.com/mykytaserdiuk/serc/ast"
)

type Interpreter struct {
	functions  map[string]*ast.Func
	structures map[string]*ast.Structure
	builtinFns map[string]ast.BuiltinFunc
	env        *Environment
}

func NewInterpreter(content string) *Interpreter {
	mainParser := &Parser{
		l: NewLexer(content),
	}

	moduleLoader := ModuleLoader{
		loaded: make(map[string]*ast.Program),
	}
	i := &Interpreter{}
	moduleLoader.Init(i)

	parsedProgram := mainParser.parseProgram()
	for _, imp := range parsedProgram.Imports {
		if funcs, ok := moduleLoader.TryLoadBuildin(imp.Name); ok {
			for name, fn := range funcs {
				fnName := imp.Alias + "." + name
				fmt.Println(fnName, "imported")
				parsedProgram.BuildinFns[fnName] = fn
			}
			continue
		}

		module := moduleLoader.Load(imp.Name)
		for _, mStr := range module.Structs {
			stName := imp.Alias + "." + mStr.Name
			fmt.Println(stName, "improted")
			parsedProgram.Structs[stName] = mStr
		}
		for _, mFn := range module.Functions {
			fnName := imp.Alias + "." + mFn.Name
			fmt.Println(fnName, "improted")
			parsedProgram.Functions[fnName] = mFn
		}
	}

	i.functions = parsedProgram.Functions
	i.structures = parsedProgram.Structs
	i.builtinFns = parsedProgram.BuildinFns
	i.env = NewEnvironment()

	for name, fn := range parsedProgram.Functions {
		i.env.Set(name, ast.Value{
			Type: ast.FuncValue,
			Data: ast.FunctionValue{
				Func: fn,
			},
		})
	}

	return i
}

func (i *Interpreter) execute(fn *ast.Func) ast.FuncResult {
	return i.exressionExecute(fn.Body)
}

func (i *Interpreter) Main() {
	mainFn := i.findFunc("main")
	if mainFn == nil {
		fmt.Printf("ERROR: cant find main ast.func")
		return
	}

	i.execute(mainFn)
}

func (i *Interpreter) exressionExecute(statements []ast.Statement) ast.FuncResult {
	for _, state := range statements {
		switch s := state.(type) {
		case ast.Call:
			switch s.Name() {
			case "printf":
				format, ok := s.Args[0].Value.(ast.StringLiteral)
				if !ok {
					panic("RUNTIME ERROR: format isnt string")
				}
				args := s.Args[1:]
				i.printf(format.Value, args...)
			case "print":
				vals := make([]any, 0)
				for _, a := range s.Args {
					vals = append(vals, i.eval(a.Value).Data)
				}
				fmt.Print(vals...)
			default:
				buildinFn := i.findBuiltin(s.Name())
				if buildinFn != nil {
					argsValues := make([]ast.Value, len(s.Args))
					for idx, argExpr := range s.Args {
						argsValues[idx] = i.eval(argExpr.Value)
					}

					for idx, argName := range s.Args {
						if idx < len(argsValues) {
							i.env.Set(argName.Name, argsValues[idx])
						}
					}

					result := buildinFn(argsValues)
					if result.Value != nil {
						return result
					}
				}
				exFn := i.findFunc(s.Name())
				if exFn != nil {
					argsValues := make([]ast.Value, len(s.Args))

					for idx, argExpr := range s.Args {
						argsValues[idx] = i.eval(argExpr.Value)
					}

					for idx, argName := range exFn.Params {
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
		case ast.NewAssign:
			computedValue := i.eval(s.Value)
			i.env.Set(s.VarName, computedValue)
		case ast.Assign:
			value := i.eval(s.Value)
			switch t := s.Target.(type) {
			case ast.Variable:
				v, ok := i.env.Get(t.Name)
				if !ok {
					panic("RUNTIME ERROR: variable '" + t.Name + "' not found")
				}
				v.Data = value.Data
				i.env.Set(t.Name, v)
			case ast.FieldAccess:
				obj := i.eval(t.Value)

				switch data := obj.Data.(type) {
				case ast.Object:
					_, ok := data.Fields[t.Name]
					if !ok {
						panic("RUNTIME ERROR: field '" + t.Name + "' not found")
					}
					data.Fields[t.Name] = value
				default:
					panic("RUNTIME ERROR: object has no fields")
				}
			}
		case ast.Return:
			return ast.FuncResult{
				Value: s.Value,
			}
		case ast.If:
			conditions := i.eval(s.Conditions)
			binaryConditionResult := conditions.Data.(bool)
			if binaryConditionResult {
				i.exressionExecute(s.Then.Statements)
			} else if !binaryConditionResult && len(s.Then.Statements) > 0 {
				i.exressionExecute(s.Else.Statements)
			}
		case ast.Loop:
			conditions := i.eval(s.Conditions)
			cond, ok := conditions.Data.(bool)
			if !ok {
				panic("loop condition must be bool")
			}
			for cond {
				i.exressionExecute(s.Body.Statements)

				conditions = i.eval(s.Conditions)
				cond, ok = conditions.Data.(bool)
				if !ok {
					panic("loop condition must be bool")
				}
			}
		}
	}
	return ast.FuncResult{
		Value: nil,
	}
}

func (i *Interpreter) executeWithArgs(fn *ast.Func, args []ast.Value) ast.Value {
	oldEnv := i.env
	env := NewEnvironmentWithParent(oldEnv)
	for idx, param := range fn.Params {
		if idx >= len(args) {
			break
		}
		env.Set(param, args[idx])
	}

	i.env = env
	defer func() {
		i.env = oldEnv
	}()

	var result ast.Value
	result = i.eval(i.exressionExecute(fn.Body))

	return result
}

func (i *Interpreter) findFunc(name string) *ast.Func {
	if fn, ok := i.functions[name]; ok {
		return fn
	}
	return nil
}
func (i *Interpreter) findBuiltin(name string) ast.BuiltinFunc {
	if fn, ok := i.builtinFns[name]; ok {
		return fn
	}

	return nil
}

func (i *Interpreter) eval(expr ast.Expression) ast.Value {
	switch e := expr.(type) {
	case ast.StringLiteral:
		return ast.GetStringValue(e.Value)
	case ast.NumberLiteral:
		return ast.GetIntValue(e.Value)
	case ast.Call:
		return i.evalCall(e)
	case nil:
		return ast.GetNullValue()
	case ast.Binary:
		return i.evaBinary(e)
	case ast.FuncResult:
		return i.eval(e.Value)
	case ast.FieldAccess:
		obj := i.eval(e.Value)
		name := e.Name
		switch v := obj.Data.(type) {
		case ast.Object:
			field, ok := v.Fields[name]
			if !ok {
				panic("field not found: " + name)
			}

			return field

		default:
			panic("cannot access field " + name)
		}
	case ast.Variable:
		val, ok := i.env.Get(e.Name)
		if !ok {
			panic(fmt.Sprintf("Runtime Error: variable '%s' is not defined", e.Name))
		}
		return val
	}

	return ast.GetNullValue()
}

func (i *Interpreter) evalCall(c ast.Call) ast.Value {
	st, ok := i.structures[c.Name()]
	if ok {
		return i.createStruct(st, c.Args)
	}
	buildinFn := i.findBuiltin(c.Name())
	if buildinFn != nil {
		argsValues := make([]ast.Value, len(c.Args))
		for idx, argExpr := range c.Args {
			argsValues[idx] = i.eval(argExpr.Value)
		}
		oldEnv := i.env
		i.env = NewEnvironment()

		for idx, argName := range c.Args {
			if idx < len(argsValues) {
				i.env.Set(argName.Name, argsValues[idx])
			}
		}
		var result ast.Value
		res := buildinFn(argsValues)

		if res.Value != nil {
			result = i.eval(res)
		} else {
			result = ast.GetNullValue()
		}
		i.env = oldEnv
		return result
	}
	fn := i.findFunc(c.Name())
	if fn != nil {
		args := make([]ast.Value, len(c.Args))

		for idx, arg := range c.Args {
			args[idx] = i.eval(arg.Value)
		}

		result := i.executeWithArgs(fn, args)
		return result
	}
	panic("unknown call: " + c.Name())
}

func (i *Interpreter) callValue(fn ast.Value, args []ast.Value) ast.Value {
	if fn.Type != ast.FuncValue {
		panic("not callable")
	}

	f := fn.Data.(ast.FunctionValue).Func

	return i.executeWithArgs(f, args)
}

func (i *Interpreter) createStruct(str *ast.Structure, args []ast.Argument) ast.Value {
	obj := ast.Object{
		Type:   str,
		Fields: make(map[string]ast.Value),
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
	return ast.Value{
		Type: ast.ObjectValue,
		Data: obj,
	}
}

func (i *Interpreter) evaBinary(bin ast.Binary) ast.Value {
	leftValue := i.eval(bin.Left)
	rightVal := i.eval(bin.Right)
	switch leftValue.Type {
	case ast.IntValue:
		lint := leftValue.Data.(int)
		rint, ok := rightVal.Data.(int)
		if !ok {
			panic("RUNTIME ERROR: expected 'int' at right part of ast.binary")
		}
		switch bin.Op.Type_ {
		case ast.LessTokenType:
			return ast.GetBoolValue(lint < rint)
		case ast.MoreTokenType:
			return ast.GetBoolValue(lint > rint)
		case ast.EqLessTokenType:
			return ast.GetBoolValue(lint <= rint)
		case ast.EqMoreTokenType:
			return ast.GetBoolValue(lint >= rint)
		case ast.EqEqTokenType:
			return ast.GetBoolValue(lint == rint)
		case ast.NotEqTokenType:
			return ast.GetBoolValue(lint != rint)
		case ast.PlusTokenType:
			return ast.GetIntValue(lint + rint)
		case ast.MinusTokenType:
			return ast.GetIntValue(lint - rint)
		}
	case ast.StringValue:
		lStr := leftValue.Data.(string)
		rStr, ok := rightVal.Data.(string)
		if !ok {
			panic("RUNTIME ERROR: expected 'string' at left part of ast.binary")
		}

		llen := len(lStr)
		rlen := len(rStr)
		switch bin.Op.Type_ {
		case ast.LessTokenType:
			return ast.GetBoolValue(llen < rlen)
		case ast.MoreTokenType:
			return ast.GetBoolValue(llen > rlen)
		case ast.EqLessTokenType:
			return ast.GetBoolValue(llen <= rlen)
		case ast.EqMoreTokenType:
			return ast.GetBoolValue(llen >= rlen)
		case ast.EqEqTokenType:
			return ast.GetBoolValue(lStr == rStr)
		case ast.NotEqTokenType:
			return ast.GetBoolValue(lStr != rStr)
		case ast.PlusTokenType, ast.MinusTokenType:
			panic("RUNTIME ERROR: unexpected operator to string: '" + bin.Op.Type_ + "'")
		}
	}
	panic("Row " + strconv.Itoa(bin.Op.Line) + ": RUNTIME ERROR: Unexpected type to calculate ast.binary")
}

func (i *Interpreter) Call(fn ast.Value, args []ast.Value) ast.Value {
	return i.callValue(fn, args)
}

// macros
func (i *Interpreter) printf(format string, exps ...ast.Argument) {
	values := make([]any, len(exps))
	for idx, e := range exps {
		values[idx] = i.eval(e.Value).Data
	}
	fmt.Printf(format, values...)
}
