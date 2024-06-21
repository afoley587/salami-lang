package main

import (
	"fmt"
	"os"

	"github.com/afoley/salami-lang/interpreter"
	"github.com/afoley/salami-lang/lexer"
	"github.com/afoley/salami-lang/parser"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		panic("No file given...") // do some better cli/usage stuff here
	}
	file, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}

	lexer := lexer.NewLexer(file)
	p := parser.New(lexer)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		fmt.Println("parser errors:")
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
		return
	}

	fmt.Printf("Parsed Program: %+v\n", program)

	interp := interpreter.New()
	result := interp.Interpret(program)

	if interp.Exited {
		fmt.Printf("Program exited with value: %v\n", interp.ExitCode)
	} else {
		fmt.Printf("Result: %v\n", result)
	}

}
