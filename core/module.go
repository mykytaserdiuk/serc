package core

import (
	"os"
	"path"

	"github.com/mykytaserdiuk/serc/ast"
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

func (l *ModuleLoader) loadBuildin(name string) *ast.Program {
	if p, ok := l.loaded[name]; ok {
		return p
	}
	program := ast.Program{}
	return &program
}
