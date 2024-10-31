package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	inputFile  = flag.String("i", "./examples/lab1.2/example.gl", "input file location")
	outputFile = flag.String("o", "./examples/lab1.2/golexgen/lexer.go", "output file location")
	printTree  = flag.Bool("t", false, "if true -> print tree")
)

func main() {
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	scn := NewScanner(bufio.NewReader(file), NewCompiler())
	tokens := GetTokens(scn)
	parser := New(tokens)

	parse, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}

	var automatas []*FiniteState
	for _, rule := range parse.rules {
		rule.expr.Print("")
		automatas = append(automatas, rule.expr.Compile())
	}
	automatas[0].ToGraph(os.Stdout)
	if automatas[0].Execute("iuq2JgJR{75}") == false {
		automatas[0].Execute("iuq2JgJR{75}")
	}
	fmt.Println(automatas[0].Execute("iuq2JgJR{75}"))

	fmt.Println(parse)

	gen := parse.Process()

	generateFile("templates/lexer.tmpl", *outputFile, gen)
}
