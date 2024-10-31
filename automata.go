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
}

func NewAutomata() *FiniteState {
	return &FiniteState{
		nextState:      1,
		CurrentState:   0,
		TerminalStates: []int{0},
		Transitions:    make(map[int]map[rune]int),
	}
}

func Create(chars []rune) *FiniteState {
	f := NewAutomata()
	f.AddTransition(0, 1, chars)
	f.TerminalStates = []int{1}
	return f
}

func (f *FiniteState) GetTerminalStates() []int {
	return f.TerminalStates
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

	if transitionsFrom, ok := f.Transitions[from]; ok {
		for _, ch := range chars {
			if f.isTerminal(transitionsFrom[ch]) {
				f.addTerminal(to)
			}
			transitionsFrom[ch] = to
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
	if (len(f.TerminalStates) == 0 && len(f.Transitions) == 0) ||
		(len(other.TerminalStates) == 0 && len(other.Transitions) == 0) {
		*f = FiniteState{}
		return
	}
	if len(f.Transitions) == 0 && len(f.TerminalStates) > 0 {
		*f = *other
		return
	}
	if len(other.Transitions) == 0 && len(other.TerminalStates) > 0 {
		return
	}

	result := &FiniteState{
		nextState:      len(f.Transitions) + len(other.Transitions),
		CurrentState:   0,
		Transitions:    make(map[int]map[rune]int),
		TerminalStates: make([]int, 0),
	}

	for _, finalState := range f.TerminalStates {
		result.Transitions[finalState] = make(map[rune]int)
	}

	for state, transitions := range f.Transitions {
		result.Transitions[state] = make(map[rune]int)
		for symbol, nextState := range transitions {
			result.Transitions[state][symbol] = nextState
		}
	}

	for _, finalState := range f.TerminalStates {
		if transitions, ok := other.Transitions[0]; ok {
			for symbol, nextState := range transitions {
				result.Transitions[finalState][symbol] = nextState + len(f.Transitions)
			}
		}
	}

	for state, transitions := range other.Transitions {
		if state == 0 {
			continue
		}
		result.Transitions[state+len(f.Transitions)] = make(map[rune]int)
		for symbol, nextState := range transitions {
			result.Transitions[state+len(f.Transitions)][symbol] = nextState + len(f.Transitions)
		}
	}

	newFinalStates := make(map[int]struct{})
	for _, finalState := range other.TerminalStates {
		if finalState == 0 {
			for _, fFinalState := range f.TerminalStates {
				newFinalStates[fFinalState] = struct{}{}
			}
		}
		newFinalStates[finalState+len(f.Transitions)] = struct{}{}
	}

	for finalState := range newFinalStates {
		result.TerminalStates = append(result.TerminalStates, finalState)
	}

	*f = *result
}

func (f *FiniteState) Union(other *FiniteState) {
	numStates := len(f.Transitions) + len(other.Transitions)
	newFSM := &FiniteState{
		nextState:      numStates,
		CurrentState:   0,
		TerminalStates: []int{},
		Transitions:    make(map[int]map[rune]int),
	}

	for state, transitions := range f.Transitions {
		newFSM.Transitions[state] = make(map[rune]int)
		for symbol, nextState := range transitions {
			newFSM.Transitions[state][symbol] = nextState
		}
	}

	for state, transitions := range other.Transitions {
		if state == 0 {
			for symbol, nextState := range transitions {
				newFSM.Transitions[state][symbol] = nextState + len(f.Transitions)
			}
		} else {
			newFSM.Transitions[state+len(f.Transitions)] = make(map[rune]int)
			for symbol, nextState := range transitions {
				newFSM.Transitions[state+len(f.Transitions)][symbol] = nextState + len(f.Transitions)
			}
		}
	}

	newFSM.TerminalStates = append(newFSM.TerminalStates, f.TerminalStates...)
	for _, tState := range other.TerminalStates {
		if tState == 0 {
			newFSM.TerminalStates = append(newFSM.TerminalStates, tState)
		} else {
			newFSM.TerminalStates = append(newFSM.TerminalStates, tState+len(f.Transitions))
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

func (f *FiniteState) Execute(input string) bool {
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
