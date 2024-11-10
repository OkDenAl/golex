package main

import (
	"fmt"
	"io"
	"sort"
)

// FiniteState struct
type FiniteState struct {
	nextState      int
	CurrentState   int
	TerminalStates []int
	Transitions    map[int]map[rune]int
	equivalents    map[int]int
}

func NewAutomata() *FiniteState {
	return &FiniteState{
		nextState:      1,
		CurrentState:   0,
		TerminalStates: []int{0},
		Transitions:    make(map[int]map[rune]int),
		equivalents:    map[int]int{},
	}
}

func Create(chars []rune) *FiniteState {
	f := NewAutomata()
	f.AddTransition(0, 1, chars)
	f.TerminalStates = []int{1}
	return f
}

func (f *FiniteState) copy() *FiniteState {
	copyState := &FiniteState{
		nextState:      f.nextState,
		CurrentState:   f.CurrentState,
		TerminalStates: make([]int, len(f.TerminalStates)),
		Transitions:    make(map[int]map[rune]int),
		equivalents:    make(map[int]int),
	}

	copy(copyState.TerminalStates, f.TerminalStates)
	for key, value := range f.Transitions {
		copyState.Transitions[key] = make(map[rune]int)
		for k, v := range value {
			copyState.Transitions[key][k] = v
		}
	}
	for key, value := range f.equivalents {
		copyState.equivalents[key] = value
	}

	return copyState
}

func (f *FiniteState) addTerminal(terminal int) {
	for _, val := range f.TerminalStates {
		if val == terminal {
			return
		}
	}
	f.TerminalStates = append(f.TerminalStates, terminal)
}

func (f *FiniteState) AddTransition(from, to int, chars []rune) {
	if from > f.nextState {
		f.nextState = from + 1
	}
	if to > f.nextState {
		f.nextState = to + 1
	}

	if v, ok := f.equivalents[from]; ok {
		from = v
	}
	if v, ok := f.equivalents[to]; ok {
		to = v
	}

	if transitionsFrom, ok := f.Transitions[from]; ok {
		for _, ch := range chars {
			if f.isTerminal(to) {
				if v, ok := transitionsFrom[ch]; ok {
					f.addTerminal(v)
				}
			}
			if v, ok := transitionsFrom[ch]; ok {
				f.equivalents[to] = v
			} else {
				f.Transitions[from][ch] = to
			}
		}
	} else {
		transitionsFrom = make(map[rune]int)
		for _, ch := range chars {
			transitionsFrom[ch] = to
		}
		f.Transitions[from] = transitionsFrom
	}
}

func (f *FiniteState) Append(other *FiniteState) {
	statesNumF := f.nextState
	statesNumOther := other.nextState
	if (len(f.TerminalStates) == 0 && statesNumF == 0) ||
		(len(other.TerminalStates) == 0 && statesNumOther == 0) {
		*f = FiniteState{}
		return
	}
	if statesNumF == 0 && len(f.TerminalStates) > 0 {
		*f = *other
		return
	}
	if statesNumOther == 0 && len(other.TerminalStates) > 0 {
		return
	}

	result := &FiniteState{
		nextState:      statesNumF + statesNumOther,
		CurrentState:   0,
		Transitions:    make(map[int]map[rune]int),
		TerminalStates: make([]int, 0),
		equivalents:    map[int]int{},
	}

	newFinalStates := make(map[int]struct{})
	for _, finalState := range other.TerminalStates {
		if finalState == 0 {
			for _, fs := range f.TerminalStates {
				newFinalStates[fs] = struct{}{}
			}
		}
		newFinalStates[finalState+statesNumF] = struct{}{}
	}

	for finalState := range newFinalStates {
		result.TerminalStates = append(result.TerminalStates, finalState)
	}

	for state, transitions := range f.Transitions {
		for symbol, nextState := range transitions {
			result.AddTransition(state, nextState, []rune{symbol})
		}
	}

	for _, finalState := range f.TerminalStates {
		if transitions, ok := other.Transitions[0]; ok {
			for symbol, nextState := range transitions {
				result.AddTransition(finalState, nextState+statesNumF, []rune{symbol})
			}
		}
	}

	keys := getSortedIntKeys(other.Transitions)
	for _, state := range keys {
		if state == 0 {
			continue
		}
		for symbol, nextState := range other.Transitions[state] {
			result.AddTransition(state+statesNumF, nextState+statesNumF, []rune{symbol})
		}
	}

	*f = *result
}

func (f *FiniteState) Union(other *FiniteState) {
	statesNumF := f.nextState
	statesNumOther := other.nextState
	newFSM := &FiniteState{
		nextState:      statesNumF + statesNumOther,
		CurrentState:   0,
		TerminalStates: []int{},
		Transitions:    make(map[int]map[rune]int),
		equivalents:    map[int]int{},
	}

	newFSM.TerminalStates = append(newFSM.TerminalStates, f.TerminalStates...)
	for _, tState := range other.TerminalStates {
		if tState == 0 {
			newFSM.TerminalStates = append(newFSM.TerminalStates, tState)
		} else {
			newFSM.TerminalStates = append(newFSM.TerminalStates, tState+statesNumF)
		}
	}

	uniqueTerminalStates := make(map[int]struct{})
	for _, tState := range newFSM.TerminalStates {
		uniqueTerminalStates[tState] = struct{}{}
	}
	newFSM.TerminalStates = []int{}
	for tState := range uniqueTerminalStates {
		newFSM.TerminalStates = append(newFSM.TerminalStates, tState)
	}

	for state, transitions := range f.Transitions {
		for symbol, nextState := range transitions {
			newFSM.AddTransition(state, nextState, []rune{symbol})
		}
	}

	keys := getSortedIntKeys(other.Transitions)
	for _, state := range keys {
		if state == 0 {
			for symbol, nextState := range other.Transitions[state] {
				newFSM.AddTransition(state, nextState+statesNumF, []rune{symbol})
			}
		} else {
			for symbol, nextState := range other.Transitions[state] {
				newFSM.AddTransition(state+statesNumF, nextState+statesNumF, []rune{symbol})
			}
		}
	}

	*f = *newFSM
}

func (f *FiniteState) Loop() {
	for from, transition := range f.Transitions {
		if from == 0 {
			for _, state := range f.TerminalStates {
				for ch, to := range transition {
					f.AddTransition(state, to, []rune{ch})
				}
			}
		}
	}

	f.addTerminal(0)
}

func (f *FiniteState) Negate() {
	var terminals []int

	for from, transition := range f.Transitions {
		if !f.isTerminal(from) {
			terminals = append(terminals, from)
		}

		for _, to := range transition {
			if !f.isTerminal(to) {
				terminals = append(terminals, to)
			}
		}
	}

	f.TerminalStates = terminals
}

func (f *FiniteState) String() string {
	sort.Ints(f.TerminalStates)
	str := fmt.Sprintf("Terminals: %v\n", f.TerminalStates)
	for from, transition := range f.Transitions {
		tran := ""
		for ch, to := range transition {
			tran += fmt.Sprintf("\n    %c => %d", ch, to)
		}
		str += fmt.Sprintf("%d: %s\n", from, tran)
	}
	return str
}

func (f *FiniteState) ToGraph(out io.Writer) {
	fmt.Fprintln(out, "digraph {")
	fmt.Fprintln(out, "0 [color=\"green\"]")

	for _, state := range f.TerminalStates {
		fmt.Fprintf(out, "%d [peripheries = 2]\n", state)
	}

	for from, trans := range f.Transitions {
		for label, to := range trans {
			fmt.Fprintf(out, "%d -> %d [label=\"%c\"]\n", from, to, label)
		}
	}

	fmt.Fprintln(out, "}")
}

func (f *FiniteState) MatchString(input string) bool {
	f.CurrentState = 0
	for _, ch := range input {
		if !f.canMoveBy(ch) {
			return false
		}
	}
	return f.isTerminal(f.CurrentState)
}

func (f *FiniteState) FindMatchEndIndex(input string) int {
	f.CurrentState = 0
	i := 0
	for _, ch := range input {
		if !f.canMoveBy(ch) {
			break
		}
		i++
	}

	if f.isTerminal(f.CurrentState) {
		return i
	}

	return 0
}

func (f *FiniteState) canMoveBy(ch rune) bool {
	from := f.CurrentState
	if to, ok := f.Transitions[from][ch]; ok {
		f.CurrentState = to
		return true
	}

	for k, to := range f.Transitions[from] {
		if k == -1 {
			f.CurrentState = to
			return true
		}
	}

	return false
}

func (f *FiniteState) isTerminal(state int) bool {
	for _, val := range f.TerminalStates {
		if state == val {
			return true
		}
	}
	return false
}
