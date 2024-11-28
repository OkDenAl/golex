package main

import (
	"strconv"
	"unicode"
)

const endSymbol = unicode.MaxRune

type Pos struct {
	followPos []int
}

var (
	flPos   = make(map[int]Pos)
	letters = make(map[rune]map[int]struct{})
)

func (f *FiniteAutomata) CompileV2() *FiniteAutomata {
	res := NewAutomata()
	usedTerminalStates := make(map[int]struct{})
	res.TerminalStates = []TerminalState{}
	i := 0
	states := make(map[string]int)
	states[arrToString(f.firstpos)] = i
	var nextState []int
	var stStack [][]int
	stStack = append(stStack, f.firstpos)
	for len(stStack) != 0 {
		curState := stStack[0]
		if el, ok := containsEndSymbol(curState); ok {
			if _, ok := usedTerminalStates[states[arrToString(curState)]]; !ok {
				res.TerminalStates = append(res.TerminalStates, TerminalState{
					State:     states[arrToString(curState)],
					LexemName: naming[el],
				})
				usedTerminalStates[states[arrToString(curState)]] = struct{}{}
			}
		}

		stStack = stStack[1:]
		for key := range letters {
			if key == endSymbol {
				continue
			}
			for _, el := range curState {
				if _, ok := letters[key][el]; ok {
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

			if el, ok := containsEndSymbol(nextState); ok {
				if _, ok := usedTerminalStates[states[arrToString(nextState)]]; !ok {
					res.TerminalStates = append(res.TerminalStates, TerminalState{
						State:     states[arrToString(nextState)],
						LexemName: naming[el],
					})
					usedTerminalStates[states[arrToString(nextState)]] = struct{}{}
				}
			}

			nextState = []int{}
		}
	}
	res.flPos = copyMap(flPos)
	res.letters = copyMap(letters)
	res.nullable = f.nullable
	res.firstpos = f.firstpos
	res.lastpos = f.lastpos
	res.lettersCount = f.lettersCount

	letters = make(map[rune]map[int]struct{})
	flPos = make(map[int]Pos)
	return res
}

func containsEndSymbol(arr []int) (int, bool) {
	for _, el := range arr {
		if _, ok := letters[endSymbol][el]; ok {
			return el, true
		}
	}

	return 0, false
}

func arrToString(arr []int) string {
	s := ""
	for _, el := range arr {
		s += strconv.Itoa(el) + "-"
	}

	return s
}
