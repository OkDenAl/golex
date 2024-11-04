package main

import (
	"fmt"
)

type Compiler struct {
	messages    map[Position]Message
	namesTokens map[string][]Token
}

func NewCompiler() *Compiler {
	return &Compiler{namesTokens: map[string][]Token{}, messages: make(map[Position]Message)}
}

func (c *Compiler) addToken(name string, token Token) {
	if _, ok := c.namesTokens[name]; !ok {
		c.namesTokens[name] = make([]Token, 0)
	}
	c.namesTokens[name] = append(c.namesTokens[name], token)
}

func (c *Compiler) hasName(name string) bool {
	_, ok := c.namesTokens[name]
	return ok
}

func (c *Compiler) addMessage(isErr bool, p Position, text string) {
	c.messages[p] = NewMessage(isErr, text)
}

func (c *Compiler) OutputMessages() {
	list := getSortedPositionKeys(c.messages)
	for _, key := range list {
		val := c.messages[key]
		if val.isError {
			fmt.Print("Error")
		} else {
			fmt.Print("Warning")
		}
		fmt.Print(" ", key.String(), ": ")
		fmt.Println(val.text)
	}
}
