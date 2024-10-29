package main

import (
	"fmt"
	"log"
	"os"

	"golex/examples/lab1.2/golexgen"
)

type Tag int

const (
	ErrTag Tag = iota
	IdentTag
	NumTag
	AssemblyCommandTag
)

func (t Tag) GetTag() string {
	var tagToString = map[Tag]string{
		IdentTag:           "IDENT",
		NumTag:             "NUM",
		AssemblyCommandTag: "ASSEMBLY",
		ErrTag:             "ERR",
	}

	return tagToString[t]
}

type MyHandler struct {
	golexgen.ErrHandlerBase
}

func (m MyHandler) Skip(text []rune, start, end golexgen.Position) *golexgen.Token {
	return nil
}

func (m MyHandler) Ident(text []rune, start, end golexgen.Position) *golexgen.Token {
	token := golexgen.NewToken(IdentTag, start, end, string(text[start.Index():end.Index()]))
	return &token
}

func (m MyHandler) Num(text []rune, start, end golexgen.Position) *golexgen.Token {
	token := golexgen.NewToken(NumTag, start, end, string(text[start.Index():end.Index()]))
	return &token
}

func (m MyHandler) Assembly(text []rune, start, end golexgen.Position) *golexgen.Token {
	token := golexgen.NewToken(AssemblyCommandTag, start, end, string(text[start.Index():end.Index()]))
	return &token
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage must be: go run main.go <fileTag.txt>\n")
	}
	filePath := os.Args[1]

	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	scn := golexgen.NewScanner([]rune(string(content)), &MyHandler{})

	t := scn.NextToken()
	for t.Tag() != golexgen.EOP {
		fmt.Println(t.String())
		t = scn.NextToken()
	}
}
