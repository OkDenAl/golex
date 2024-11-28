package main

import (
	"fmt"
	"io"
)

type FiniteAutomata struct {
	CurrentState   int
	TerminalStates []TerminalState
	Transitions    map[int]map[rune]int

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

func NewAutomata() *FiniteAutomata {
	return &FiniteAutomata{
		TerminalStates: []TerminalState{{State: 0}},
		Transitions:    make(map[int]map[rune]int),
		nullable:       true,
		firstpos:       make([]int, 0),
		lastpos:        make([]int, 0),
		lettersCount:   0,
	}
}

func Create(chars []rune, pos int) *FiniteAutomata {
	f := NewAutomata()
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

func (f *FiniteAutomata) setLexemName(name string) {
	for i := range f.TerminalStates {
		f.TerminalStates[i].LexemName = name
	}
}

func (f *FiniteAutomata) copy() *FiniteAutomata {
	copyState := &FiniteAutomata{
		CurrentState:   f.CurrentState,
		TerminalStates: make([]TerminalState, len(f.TerminalStates)),
		Transitions:    make(map[int]map[rune]int),
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

	return copyState
}

func (f *FiniteAutomata) Concat(other *FiniteAutomata) {
	result := &FiniteAutomata{
		CurrentState:   0,
		Transitions:    make(map[int]map[rune]int),
		TerminalStates: make([]TerminalState, 0),
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

	*f = *result
}

func (f *FiniteAutomata) UnionNext(other *FiniteAutomata) *FiniteAutomata {
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

	newFSM := &FiniteAutomata{
		CurrentState:   0,
		TerminalStates: []TerminalState{},
		Transitions:    make(map[int]map[rune]int),
		nullable:       f.nullable || other.nullable,
		firstpos:       mergeUnique(f.firstpos, other.firstpos),
		lastpos:        mergeUnique(f.lastpos, other.lastpos),
		lettersCount:   f.lettersCount + other.lettersCount,
		flPos:          copyMap(f.flPos),
		letters:        copyMap(f.letters),
	}

	return newFSM
}

func (f *FiniteAutomata) Union(other *FiniteAutomata) {
	newFSM := &FiniteAutomata{
		CurrentState:   0,
		TerminalStates: []TerminalState{},
		Transitions:    make(map[int]map[rune]int),
		nullable:       f.nullable || other.nullable,
		firstpos:       mergeUnique(f.firstpos, other.firstpos),
		lastpos:        mergeUnique(f.lastpos, other.lastpos),
		lettersCount:   f.lettersCount + other.lettersCount,
	}

	*f = *newFSM
}

func (f *FiniteAutomata) Loop() {
	f.nullable = true
	for _, p := range f.lastpos {
		r := mergeUnique(flPos[p].followPos, f.firstpos)
		flPos[p] = Pos{followPos: r}
	}
}

func Negate(f *FiniteAutomata, pos int) *FiniteAutomata {
	var symb []rune
	for _, transition := range f.Transitions {
		for r, _ := range transition {
			symb = append(symb, r)
		}
	}

	return Create(genAnyRuneExcept(symb), pos)
}

func (f *FiniteAutomata) ToGraph(out io.Writer) {
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

func (f *FiniteAutomata) MatchString(input string) bool {
	f.CurrentState = 0
	for _, ch := range input {
		if !f.canMoveBy(ch) {
			return false
		}
	}
	return f.isTerminal(f.CurrentState)
}

func (f *FiniteAutomata) FindMatchEndIndex(input string) int {
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

func (f *FiniteAutomata) canMoveBy(ch rune) bool {
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

func (f *FiniteAutomata) isTerminal(state int) bool {
	for _, val := range f.TerminalStates {
		if state == val.State {
			return true
		}
	}
	return false
}
