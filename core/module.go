package core

import (
	"os"
	"path"

	"github.com/mykytaserdiuk/serc/ast"
	"github.com/mykytaserdiuk/serc/builtin"
	"github.com/mykytaserdiuk/serc/runtime"
)

type ModuleLoader struct {
	loaded map[string]*ast.Program
}

var (
	buildInModules = make(map[string]map[string]ast.BuiltinFunc)
)

func (l *ModuleLoader) Init(rt runtime.Runtime) {
	buildInModules["http"] = builtin.LoadHttp(rt)
}
func (l *ModuleLoader) Load(name string) *ast.Program {
	if p, ok := l.loaded[name]; ok {
		return p
	}
	data, err := os.ReadFile(path.Join("modules", name+".ser"))
	if err != nil {
		panic(err)
	}
	parser := &Parser{
		l: NewLexer(string(data)),
	}
	program := parser.parseProgram()
	l.loaded[name] = &program
	return &program
}

func (l *ModuleLoader) TryLoadBuildin(name string) (map[string]ast.BuiltinFunc, bool) {
	if funcs, ok := buildInModules[name]; ok {
		return funcs, true
	}

	return nil, false
}
