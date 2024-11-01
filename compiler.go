package main

import (
	"fmt"
	"reflect"
	"sort"
)

func SortedMapKeys(m map[Position]Message) (keyList []Position) {
	keys := reflect.ValueOf(m).MapKeys()

	for _, key := range keys {
		keyList = append(keyList, key.Interface().(Position))
	}
	sort.Slice(keyList, func(i, j int) bool {
		return keyList[i].line < keyList[j].line ||
			(keyList[i].line == keyList[j].line && keyList[i].pos < keyList[j].pos)
	})
	return
}

type Compiler struct {
	messages    map[Position]Message
	namesTokens map[string][]Token
}

func NewCompiler() *Compiler {
	return &Compiler{namesTokens: map[string][]Token{}, messages: make(map[Position]Message)}
}

func (c *Compiler) AddToken(name string, token Token) {
	if _, ok := c.namesTokens[name]; !ok {
		c.namesTokens[name] = make([]Token, 0)
	}
	c.namesTokens[name] = append(c.namesTokens[name], token)
}

func (c *Compiler) HasName(name string) bool {
	_, ok := c.namesTokens[name]
	return ok
}

func (c *Compiler) AddMessage(isErr bool, p Position, text string) {
	c.messages[p] = NewMessage(isErr, text)
}

func (c *Compiler) OutputMessages() {
	list := SortedMapKeys(c.messages)
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
