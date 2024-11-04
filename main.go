package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	inputFile       = flag.String("i", "./examples/testing/example.gl", "input file location")
	lexerOutputFile = flag.String("lo", "./examples/testing/golexgen/lexer.go", "lexer output file location")
	mainOutputFile  = flag.String("mo", "./examples/testing/main.go", "main output file location")
	regenerateMain  = flag.Bool("rg", false, "if true -> regenerate main")
	printTree       = flag.Bool("t", false, "if true -> print tree")
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

	var automatas []*FiniteState
	for _, rule := range parse.rules.ruleArr {
		rule.expr.Print("")
		automatas = append(automatas, rule.expr.Compile())
	}

	gen := parse.Process()

	generateFile("templates/lexer.tmpl", *lexerOutputFile, gen, true)
	generateFile("templates/main.tmpl", *mainOutputFile, gen, *regenerateMain)
}
