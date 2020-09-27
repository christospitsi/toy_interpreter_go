package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"toy_interpreter_go/evaluator"
	"toy_interpreter_go/lexer"
	"toy_interpreter_go/object"
	"toy_interpreter_go/parser"
)

func main() {
	interpret("example1.cmm", "output1.txt")
	interpret("example2.cmm", "output2.txt")
	interpret("example3.cmm", "output3.txt")
	interpret("example4.cmm", "output4.txt")
	interpret("example5.cmm", "output5.txt")
	interpret("example6.cmm", "output6.txt")
}

func interpret(src, dst string) (int, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}

	content, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	code := string(content)

	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()

	env := object.NewEnvironment()
	lex := lexer.LexConstructor(code)
	pars := parser.ParsConstructor(lex)
	program := pars.ParseProgram()
	fmt.Println("Parsed program: \n", program.Statements)
	evaluated := evaluator.Eval(program, env)

	nBytes, err := io.WriteString(destination, evaluated.Inspect())

	return nBytes, err
}
