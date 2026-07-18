package core

import (
	"os"
	"path"

	"github.com/mykytaserdiuk/serc/ast"
	"github.com/mykytaserdiuk/serc/buildin"
)

type ModuleLoader struct {
	loaded map[string]*ast.Program
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

var (
	buildInModules = map[string]map[string]ast.BuiltinFunc{
		"http": buildin.LoadHttp(),
	}
)

func (l *ModuleLoader) TryLoadBuildin(name string) (map[string]ast.BuiltinFunc, bool) {
	if funcs, ok := buildInModules[name]; ok {
		return funcs, true
	}

	return nil, false
}
