package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	inputFile       = flag.String("i", "./examples/lab1.2/example.gl", "input file location")
	lexerOutputFile = flag.String("lo", "./examples/lab1.2/golexgen/lexer.go", "lexer output file location")
	mainOutputFile  = flag.String("mo", "./examples/lab1.2/main.go", "main output file location")
	regenerateMain  = flag.Bool("rg", false, "if true -> regenerate main")
)

func main() {
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	scn := NewScanner(bufio.NewReader(file), NewCompiler())
	tokens := scn.GetTokens()
	parser := New(tokens)

	parse, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}

	gen := parse.ProcessOneAutomata()

	generateFile("templates/lexer.tmpl", *lexerOutputFile, gen, true)
	generateFile("templates/main.tmpl", *mainOutputFile, gen, *regenerateMain)
}
