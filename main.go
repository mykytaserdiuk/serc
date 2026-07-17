package main

import (
	"fmt"
	"os"

	"github.com/mykytaserdiuk/serc/core"
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

	interpreter := core.NewInterpreter(content)
	interpreter.Main()
}
