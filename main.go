package main

import (
	"fmt"
	"os"
)

const pathToScript = "./script.ser"

func main(){
	file, err := os.ReadFile(pathToScript)
	if err != nil {
		fmt.Printf("ERROR: Open file: %s", err)
		return
	}

	content := string(file)
	lexer := NewLexer(content)

	//	for {
		//	if token, ok := lexer.NextToken(); ok{
			//	fmt.Printf("type: %s, val: %s\n", token.type_, token.value)
			//	}
		//}


	parser := &Parser{
		l: lexer,
	}
	fn:=parser.parseFunc()
	for _,  b := range fn.body{
		if call, ok := b.(FuncCall); ok{
			if call.name == "print"{
				fmt.Printf("%s \n", call.args[0])
			}
		}
	}
}
