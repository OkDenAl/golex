package main

import (
	"fmt"
	"slices"
	"strconv"
)

type Pos struct {
	followPos []int
}

var (
	flPos   = make(map[int]Pos)
	letters = make(map[rune][]int)
)

func (f *FiniteState) CompileV2() *FiniteState {
	res := NewAutomata()
	res.TerminalStates = []TerminalState{}
	i := 0
	states := make(map[string]int)
	states[arrToString(f.firstpos)] = i
	var nextState []int
	var stStack [][]int
	stStack = append(stStack, f.firstpos)
	fmt.Println(flPos)
	fmt.Println(letters)
	for len(stStack) != 0 {
		curState := stStack[0]
		if slices.Contains(curState, letters['#'][0]) {
			res.TerminalStates = append(res.TerminalStates, TerminalState{
				State: states[arrToString(curState)],
			})
		}
		stStack = stStack[1:]
		for key := range letters {
			if key == '#' {
				continue
			}
			for _, el := range curState {
				if slices.Contains(letters[key], el) {
					nextState = mergeUnique(nextState, flPos[el].followPos)
				}
			}
			if len(nextState) == 0 {
				continue
			}
			if _, ok := states[arrToString(nextState)]; !ok {
				i++
				states[arrToString(nextState)] = i
				stStack = append(stStack, nextState)
			}
			if _, ok := res.Transitions[states[arrToString(curState)]]; !ok {
				res.Transitions[states[arrToString(curState)]] = make(map[rune]int)
			}
			res.Transitions[states[arrToString(curState)]][key] = states[arrToString(nextState)]
			if slices.Contains(nextState, letters['#'][0]) {
				res.TerminalStates = append(res.TerminalStates, TerminalState{
					State: states[arrToString(nextState)],
				})
			}
			nextState = []int{}
		}
	}
	return res
}

func arrToString(arr []int) string {
	s := ""
	for _, el := range arr {
		s += strconv.Itoa(el) + "-"
	}

	return s
}
