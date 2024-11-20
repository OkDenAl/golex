// Code generated by golex utility;
// This code is present a default handling of tokens;
// YOU CAN EDIT IT IF YOU NEED.
package main

import (
	"fmt"
	"log"
	"os"
)

type Handler struct {
	golexgen.HandlerBase
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
	scn := golexgen.NewScanner([]rune(string(content)), &Handler{})

	t := scn.NextToken()
	for t.Tag() != golexgen.EOP {
		fmt.Println(t.String())
		t = scn.NextToken()
	}
}
