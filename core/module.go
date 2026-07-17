package core

import (
	"os"
	"path"
)

type ModuleLoader struct {
    loaded map[string]*Program
}

func (l *ModuleLoader) Load(name string) *Program {
    if p, ok := l.loaded[name]; ok {
        return p
    }
    data, err := os.ReadFile(path.Join("modules",name+".ser"))
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
