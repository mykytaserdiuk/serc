package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("ERROR: not enought args")
	}
	path := os.Args[1]

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("ERROR: Open file: %s", err)
		return
	}

	content := string(file)
	lexer := NewLexer(content)

	parser := &Parser{
		l: lexer,
	}

	interpreter := NewInterpreter(parser)
	interpreter.Main()
}
