package main

import (
	"fmt"
	"io"
)

// FiniteState struct
type FiniteState struct {
	nextState      int
	CurrentState   int
	TerminalStates []TerminalState
	Transitions    map[int]map[rune]int
	equivalents    map[int]int

	nullable     bool
	firstpos     []int
	lastpos      []int
	flPos        map[int]Pos
	letters      map[rune]map[int]struct{}
	lettersCount int
}

type TerminalState struct {
	State     int
	LexemName string
}

func NewAutomata() *FiniteState {
	return &FiniteState{
		nextState:      1,
		CurrentState:   0,
		TerminalStates: []TerminalState{{State: 0}},
		Transitions:    make(map[int]map[rune]int),
		equivalents:    map[int]int{},
		nullable:       true,
		firstpos:       make([]int, 0),
		lastpos:        make([]int, 0),
		lettersCount:   0,
	}
}

func Create(chars []rune, pos int) *FiniteState {
	f := NewAutomata()
	f.AddTransition(0, 1, chars)
	f.TerminalStates = []TerminalState{{State: 1}}
	f.firstpos = append(f.firstpos, pos)
	f.lastpos = append(f.lastpos, pos)
	f.nullable = false
	flPos[pos] = Pos{
		followPos: make([]int, 0),
	}
	for _, ch := range chars {
		if _, ok := letters[ch]; !ok {
			letters[ch] = make(map[int]struct{})
		}
		letters[ch][pos] = struct{}{}
	}
	f.lettersCount = 1
	return f
}

func (f *FiniteState) setLexemName(name string) {
	for i := range f.TerminalStates {
		f.TerminalStates[i].LexemName = name
	}
}

func (f *FiniteState) copy() *FiniteState {
	copyState := &FiniteState{
		nextState:      f.nextState,
		CurrentState:   f.CurrentState,
		TerminalStates: make([]TerminalState, len(f.TerminalStates)),
		Transitions:    make(map[int]map[rune]int),
		equivalents:    make(map[int]int),
		nullable:       f.nullable,
		firstpos:       f.firstpos,
		lastpos:        f.lastpos,
		lettersCount:   f.lettersCount,
		flPos:          copyMap(f.flPos),
		letters:        copyMap(f.letters),
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
		if val.State == terminal {
			return
		}
	}
	f.TerminalStates = append(f.TerminalStates, TerminalState{State: terminal})
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
		TerminalStates: make([]TerminalState, 0),
		equivalents:    map[int]int{},
		nullable:       f.nullable && other.nullable,
		firstpos:       f.firstpos,
		lastpos:        other.lastpos,
		lettersCount:   f.lettersCount + other.lettersCount,
	}

	if f.nullable {
		result.firstpos = mergeUnique(f.firstpos, other.firstpos)
	}
	if other.nullable {
		result.lastpos = mergeUnique(f.lastpos, other.lastpos)
	}

	for _, p := range f.lastpos {
		r := mergeUnique(flPos[p].followPos, other.firstpos)
		flPos[p] = Pos{followPos: r}
	}

	newFinalStates := make(map[TerminalState]struct{})
	for _, finalState := range other.TerminalStates {
		if finalState.State == 0 {
			for _, fs := range f.TerminalStates {
				newFinalStates[fs] = struct{}{}
			}
		}
		newFinalStates[TerminalState{State: finalState.State + statesNumF}] = struct{}{}
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
				result.AddTransition(finalState.State, nextState+statesNumF, []rune{symbol})
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

var naming = map[int]string{}

func (f *FiniteState) UnionNext(other *FiniteState) {
	for i := range other.firstpos {
		other.firstpos[i] += f.lettersCount
	}

	for i := range other.lastpos {
		other.lastpos[i] += f.lettersCount
	}

	for k, v := range other.flPos {
		f.flPos[k+f.lettersCount] = Pos{
			followPos: copyWithAdd(v.followPos, f.lettersCount),
		}
	}

	for k, v := range other.letters {
		for vv := range v {
			if _, ok := f.letters[k]; !ok {
				f.letters[k] = make(map[int]struct{})
			}
			f.letters[k][vv+f.lettersCount] = struct{}{}
		}

		if k == endSymbol {
			for vv := range v {
				naming[vv+f.lettersCount] = other.TerminalStates[0].LexemName
			}
		}
	}

	newFSM := &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{},
		Transitions:    make(map[int]map[rune]int),
		equivalents:    map[int]int{},
		nullable:       f.nullable || other.nullable,
		firstpos:       mergeUnique(f.firstpos, other.firstpos),
		lastpos:        mergeUnique(f.lastpos, other.lastpos),
		lettersCount:   f.lettersCount + other.lettersCount,
		flPos:          copyMap(f.flPos),
		letters:        copyMap(f.letters),
	}

	*f = *newFSM
}

func (f *FiniteState) Union(other *FiniteState) {
	statesNumF := f.nextState
	statesNumOther := other.nextState
	newFSM := &FiniteState{
		nextState:      statesNumF + statesNumOther,
		CurrentState:   0,
		TerminalStates: []TerminalState{},
		Transitions:    make(map[int]map[rune]int),
		equivalents:    map[int]int{},
		nullable:       f.nullable || other.nullable,
		firstpos:       mergeUnique(f.firstpos, other.firstpos),
		lastpos:        mergeUnique(f.lastpos, other.lastpos),
		lettersCount:   f.lettersCount + other.lettersCount,
	}

	newFSM.TerminalStates = append(newFSM.TerminalStates, f.TerminalStates...)
	for _, tState := range other.TerminalStates {
		if tState.State == 0 {
			newFSM.TerminalStates = append(newFSM.TerminalStates, tState)
		} else {
			newFSM.TerminalStates = append(newFSM.TerminalStates, TerminalState{State: tState.State + statesNumF, LexemName: tState.LexemName})
		}
	}

	uniqueTerminalStates := make(map[TerminalState]struct{})
	for _, tState := range newFSM.TerminalStates {
		uniqueTerminalStates[tState] = struct{}{}
	}
	newFSM.TerminalStates = []TerminalState{}
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
	f.nullable = true
	for from, transition := range f.Transitions {
		if from == 0 {
			for _, state := range f.TerminalStates {
				for ch, to := range transition {
					f.AddTransition(state.State, to, []rune{ch})
				}
			}
		}
	}
	for _, p := range f.lastpos {
		r := mergeUnique(flPos[p].followPos, f.firstpos)
		flPos[p] = Pos{followPos: r}
	}

	f.addTerminal(0)
}

func (f *FiniteState) Negate() {
	var terminals []TerminalState

	for from, transition := range f.Transitions {
		if !f.isTerminal(from) {
			terminals = append(terminals, TerminalState{State: from})
		}

		for _, to := range transition {
			if !f.isTerminal(to) {
				terminals = append(terminals, TerminalState{State: to})
			}
		}
	}

	f.TerminalStates = terminals
}

func (f *FiniteState) ToGraph(out io.Writer) {
	fmt.Fprintln(out, "digraph {")
	fmt.Fprintln(out, "0 [color=\"green\"]")

	for _, state := range f.TerminalStates {
		fmt.Fprintf(out, "%d [peripheries = 2]\n", state.State)
	}

	for from, trans := range f.Transitions {
		for label, to := range trans {
			if label == '\n' {
				label = 'n'
			}
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
		if state == val.State {
			return true
		}
	}
	return false
}
