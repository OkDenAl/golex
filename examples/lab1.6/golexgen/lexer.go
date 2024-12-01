// Code generated by golex utility; DO NOT EDIT.
package golexgen

import (
	"fmt"
)

type ErrHandler interface {
	Error(msg string, pos Position, symbol string)
}

type (
	ErrFunc             func(msg string, pos Position, symbol string)
	SwitchConditionFunc func(cond Condition)
)

type Continued bool

type LexemHandler interface {
	ErrHandler
	Skip(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularStart(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEnd(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularNewLine(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEscapeNewLine(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEscapeTab(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEscapeQota(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularSymb(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	StartLiteral(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	Literal1(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	LiteralEnd(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	LiteralNewLine(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	LiteralChar(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	Num(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	Any(text string, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
}

type Tag interface {
	GetTag() string
}

const EOP = "EOP"

type Condition int

const (
	dummy = Condition(iota)

	ConditionINIT
	ConditionLITERAL
	ConditionREGULAR
)

type DefaultTag int

const (
	TagErr DefaultTag = iota
	TagSkip
	TagRegularStart
	TagRegularEnd
	TagRegularNewLine
	TagRegularEscapeNewLine
	TagRegularEscapeTab
	TagRegularEscapeQota
	TagRegularSymb
	TagStartLiteral
	TagLiteral1
	TagLiteralEnd
	TagLiteralNewLine
	TagLiteralChar
	TagNum
	TagAny
	TagINIT
	TagLITERAL
	TagREGULAR
)

func (t DefaultTag) GetTag() string {
	var tagToString = map[DefaultTag]string{
		TagSkip:                 "Skip",
		TagRegularStart:         "RegularStart",
		TagRegularEnd:           "RegularEnd",
		TagRegularNewLine:       "RegularNewLine",
		TagRegularEscapeNewLine: "RegularEscapeNewLine",
		TagRegularEscapeTab:     "RegularEscapeTab",
		TagRegularEscapeQota:    "RegularEscapeQota",
		TagRegularSymb:          "RegularSymb",
		TagStartLiteral:         "StartLiteral",
		TagLiteral1:             "Literal1",
		TagLiteralEnd:           "LiteralEnd",
		TagLiteralNewLine:       "LiteralNewLine",
		TagLiteralChar:          "LiteralChar",
		TagNum:                  "Num",
		TagAny:                  "Any",
		TagINIT:                 "INIT",
		TagLITERAL:              "LITERAL",
		TagREGULAR:              "REGULAR",
	}

	return tagToString[t]
}

type FiniteState struct {
	NextState      int
	CurrentState   int
	TerminalStates []TerminalState
	Transitions    map[int]map[rune]int
}

type TerminalState struct {
	state     int
	lexemName string
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

func (f *FiniteState) FindMatchEndIndexOneAutomata(input string) (int, string) {
	f.CurrentState = 0
	i := 0
	prevI := i
	prevStr := ""
	for _, ch := range input {
		if !f.canMoveBy(ch) {
			break
		}
		i++

		if val, ok := f.isTerminalOneAutoamta(f.CurrentState); ok {
			prevI = i
			prevStr = val.lexemName
		}
	}

	return prevI, prevStr
}

func (f *FiniteState) canMoveBy(ch rune) bool {
	from := f.CurrentState
	if to, ok := f.Transitions[from][ch]; ok {
		f.CurrentState = to
		return true
	}

	return false
}

func (f *FiniteState) isTerminalOneAutoamta(state int) (TerminalState, bool) {
	for _, val := range f.TerminalStates {
		if state == val.state {
			return val, true
		}
	}
	return TerminalState{}, false
}

func (f *FiniteState) isTerminal(state int) bool {
	for _, val := range f.TerminalStates {
		if state == val.state {
			return true
		}
	}
	return false
}

var (
	automataSkip *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "Skip"}},
		Transitions: map[int]map[rune]int{
			0: {9: 1, 10: 1, 32: 1},
			1: {9: 1, 10: 1, 32: 1},
		},
	}
	automataRegularStart *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "RegularStart"}},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
		},
	}
	automataRegularEnd *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "RegularEnd"}},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
		},
	}
	automataRegularNewLine *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "RegularNewLine"}},
		Transitions: map[int]map[rune]int{
			0: {10: 1},
		},
	}
	automataRegularEscapeNewLine *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 2, lexemName: "RegularEscapeNewLine"}},
		Transitions: map[int]map[rune]int{
			0: {92: 1},
			1: {110: 2},
		},
	}
	automataRegularEscapeTab *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 2, lexemName: "RegularEscapeTab"}},
		Transitions: map[int]map[rune]int{
			0: {92: 1},
			1: {116: 2},
		},
	}
	automataRegularEscapeQota *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 2, lexemName: "RegularEscapeQota"}},
		Transitions: map[int]map[rune]int{
			0: {92: 1},
			1: {34: 2},
		},
	}
	automataRegularSymb *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 0, lexemName: "RegularSymb"}, {state: 1, lexemName: "RegularSymb"}},
		Transitions: map[int]map[rune]int{
			0: {9: 1, 32: 1, 33: 1, 34: 1, 35: 1, 36: 1, 37: 1, 38: 1, 39: 1, 40: 1, 41: 1, 42: 1, 43: 1, 44: 1, 45: 1, 46: 1, 47: 1, 48: 1, 49: 1, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1, 55: 1, 56: 1, 57: 1, 58: 1, 59: 1, 60: 1, 61: 1, 62: 1, 63: 1, 64: 1, 65: 1, 66: 1, 67: 1, 68: 1, 69: 1, 70: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 1, 76: 1, 77: 1, 78: 1, 79: 1, 80: 1, 81: 1, 82: 1, 83: 1, 84: 1, 85: 1, 86: 1, 87: 1, 88: 1, 89: 1, 90: 1, 91: 1, 92: 1, 93: 1, 94: 1, 95: 1, 96: 1, 97: 1, 98: 1, 99: 1, 100: 1, 101: 1, 102: 1, 103: 1, 104: 1, 105: 1, 106: 1, 107: 1, 108: 1, 109: 1, 110: 1, 111: 1, 112: 1, 113: 1, 114: 1, 115: 1, 116: 1, 117: 1, 118: 1, 119: 1, 120: 1, 121: 1, 122: 1, 123: 1, 124: 1, 125: 1, 126: 1},
		},
	}
	automataStartLiteral *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 2, lexemName: "StartLiteral"}},
		Transitions: map[int]map[rune]int{
			0: {64: 1},
			1: {34: 2},
		},
	}
	automataLiteral1 *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 2, lexemName: "Literal1"}},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
			1: {34: 2},
		},
	}
	automataLiteralEnd *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "LiteralEnd"}},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
		},
	}
	automataLiteralNewLine *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "LiteralNewLine"}},
		Transitions: map[int]map[rune]int{
			0: {10: 1},
		},
	}
	automataLiteralChar *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 0, lexemName: "LiteralChar"}, {state: 1, lexemName: "LiteralChar"}},
		Transitions: map[int]map[rune]int{
			0: {9: 1, 32: 1, 33: 1, 34: 1, 35: 1, 36: 1, 37: 1, 38: 1, 39: 1, 40: 1, 41: 1, 42: 1, 43: 1, 44: 1, 45: 1, 46: 1, 47: 1, 48: 1, 49: 1, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1, 55: 1, 56: 1, 57: 1, 58: 1, 59: 1, 60: 1, 61: 1, 62: 1, 63: 1, 64: 1, 65: 1, 66: 1, 67: 1, 68: 1, 69: 1, 70: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 1, 76: 1, 77: 1, 78: 1, 79: 1, 80: 1, 81: 1, 82: 1, 83: 1, 84: 1, 85: 1, 86: 1, 87: 1, 88: 1, 89: 1, 90: 1, 91: 1, 92: 1, 93: 1, 94: 1, 95: 1, 96: 1, 97: 1, 98: 1, 99: 1, 100: 1, 101: 1, 102: 1, 103: 1, 104: 1, 105: 1, 106: 1, 107: 1, 108: 1, 109: 1, 110: 1, 111: 1, 112: 1, 113: 1, 114: 1, 115: 1, 116: 1, 117: 1, 118: 1, 119: 1, 120: 1, 121: 1, 122: 1, 123: 1, 124: 1, 125: 1, 126: 1},
		},
	}
	automataNum *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 1, lexemName: "Num"}, {state: 2, lexemName: "Num"}},
		Transitions: map[int]map[rune]int{
			0: {48: 1, 49: 2},
			2: {49: 2},
		},
	}
	automataAny *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 0, lexemName: "Any"}, {state: 1, lexemName: "Any"}},
		Transitions: map[int]map[rune]int{
			0: {9: 1, 32: 1, 33: 1, 34: 1, 35: 1, 36: 1, 37: 1, 38: 1, 39: 1, 40: 1, 41: 1, 42: 1, 43: 1, 44: 1, 45: 1, 46: 1, 47: 1, 48: 1, 49: 1, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1, 55: 1, 56: 1, 57: 1, 58: 1, 59: 1, 60: 1, 61: 1, 62: 1, 63: 1, 64: 1, 65: 1, 66: 1, 67: 1, 68: 1, 69: 1, 70: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 1, 76: 1, 77: 1, 78: 1, 79: 1, 80: 1, 81: 1, 82: 1, 83: 1, 84: 1, 85: 1, 86: 1, 87: 1, 88: 1, 89: 1, 90: 1, 91: 1, 92: 1, 93: 1, 94: 1, 95: 1, 96: 1, 97: 1, 98: 1, 99: 1, 100: 1, 101: 1, 102: 1, 103: 1, 104: 1, 105: 1, 106: 1, 107: 1, 108: 1, 109: 1, 110: 1, 111: 1, 112: 1, 113: 1, 114: 1, 115: 1, 116: 1, 117: 1, 118: 1, 119: 1, 120: 1, 121: 1, 122: 1, 123: 1, 124: 1, 125: 1, 126: 1},
		},
	}

	unionAutomataINIT *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 0, lexemName: "Any"}, {state: 1, lexemName: "Any"}, {state: 2, lexemName: "RegularStart"}, {state: 3, lexemName: "Any"}, {state: 4, lexemName: ""}, {state: 5, lexemName: "Num"}, {state: 6, lexemName: ""}, {state: 7, lexemName: "Num"}, {state: 8, lexemName: "StartLiteral"}, {state: 9, lexemName: "Num"}},
		Transitions: map[int]map[rune]int{
			0: {9: 6, 10: 4, 32: 6, 33: 1, 34: 2, 35: 1, 36: 1, 37: 1, 38: 1, 39: 1, 40: 1, 41: 1, 42: 1, 43: 1, 44: 1, 45: 1, 46: 1, 47: 1, 48: 5, 49: 7, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1, 55: 1, 56: 1, 57: 1, 58: 1, 59: 1, 60: 1, 61: 1, 62: 1, 63: 1, 64: 3, 65: 1, 66: 1, 67: 1, 68: 1, 69: 1, 70: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 1, 76: 1, 77: 1, 78: 1, 79: 1, 80: 1, 81: 1, 82: 1, 83: 1, 84: 1, 85: 1, 86: 1, 87: 1, 88: 1, 89: 1, 90: 1, 91: 1, 92: 1, 93: 1, 94: 1, 95: 1, 96: 1, 97: 1, 98: 1, 99: 1, 100: 1, 101: 1, 102: 1, 103: 1, 104: 1, 105: 1, 106: 1, 107: 1, 108: 1, 109: 1, 110: 1, 111: 1, 112: 1, 113: 1, 114: 1, 115: 1, 116: 1, 117: 1, 118: 1, 119: 1, 120: 1, 121: 1, 122: 1, 123: 1, 124: 1, 125: 1, 126: 1},
			3: {34: 8},
			4: {9: 4, 10: 4, 32: 4},
			6: {9: 4, 10: 4, 32: 4},
			7: {49: 9},
			9: {49: 9},
		},
	}
	unionAutomataREGULAR *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 0, lexemName: "RegularSymb"}, {state: 1, lexemName: "RegularSymb"}, {state: 2, lexemName: "RegularNewLine"}, {state: 3, lexemName: "RegularEnd"}, {state: 4, lexemName: "RegularSymb"}, {state: 5, lexemName: "RegularEscapeTab"}, {state: 6, lexemName: "RegularEscapeNewLine"}, {state: 7, lexemName: "RegularEscapeQota"}},
		Transitions: map[int]map[rune]int{
			0: {9: 1, 10: 2, 32: 1, 33: 1, 34: 3, 35: 1, 36: 1, 37: 1, 38: 1, 39: 1, 40: 1, 41: 1, 42: 1, 43: 1, 44: 1, 45: 1, 46: 1, 47: 1, 48: 1, 49: 1, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1, 55: 1, 56: 1, 57: 1, 58: 1, 59: 1, 60: 1, 61: 1, 62: 1, 63: 1, 64: 1, 65: 1, 66: 1, 67: 1, 68: 1, 69: 1, 70: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 1, 76: 1, 77: 1, 78: 1, 79: 1, 80: 1, 81: 1, 82: 1, 83: 1, 84: 1, 85: 1, 86: 1, 87: 1, 88: 1, 89: 1, 90: 1, 91: 1, 92: 4, 93: 1, 94: 1, 95: 1, 96: 1, 97: 1, 98: 1, 99: 1, 100: 1, 101: 1, 102: 1, 103: 1, 104: 1, 105: 1, 106: 1, 107: 1, 108: 1, 109: 1, 110: 1, 111: 1, 112: 1, 113: 1, 114: 1, 115: 1, 116: 1, 117: 1, 118: 1, 119: 1, 120: 1, 121: 1, 122: 1, 123: 1, 124: 1, 125: 1, 126: 1},
			4: {34: 7, 110: 6, 116: 5},
		},
	}
	unionAutomataLITERAL *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []TerminalState{{state: 0, lexemName: "LiteralChar"}, {state: 1, lexemName: "LiteralChar"}, {state: 2, lexemName: "LiteralEnd"}, {state: 3, lexemName: "LiteralNewLine"}, {state: 4, lexemName: "Literal1"}},
		Transitions: map[int]map[rune]int{
			0: {9: 1, 10: 3, 32: 1, 33: 1, 34: 2, 35: 1, 36: 1, 37: 1, 38: 1, 39: 1, 40: 1, 41: 1, 42: 1, 43: 1, 44: 1, 45: 1, 46: 1, 47: 1, 48: 1, 49: 1, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1, 55: 1, 56: 1, 57: 1, 58: 1, 59: 1, 60: 1, 61: 1, 62: 1, 63: 1, 64: 1, 65: 1, 66: 1, 67: 1, 68: 1, 69: 1, 70: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 1, 76: 1, 77: 1, 78: 1, 79: 1, 80: 1, 81: 1, 82: 1, 83: 1, 84: 1, 85: 1, 86: 1, 87: 1, 88: 1, 89: 1, 90: 1, 91: 1, 92: 1, 93: 1, 94: 1, 95: 1, 96: 1, 97: 1, 98: 1, 99: 1, 100: 1, 101: 1, 102: 1, 103: 1, 104: 1, 105: 1, 106: 1, 107: 1, 108: 1, 109: 1, 110: 1, 111: 1, 112: 1, 113: 1, 114: 1, 115: 1, 116: 1, 117: 1, 118: 1, 119: 1, 120: 1, 121: 1, 122: 1, 123: 1, 124: 1, 125: 1, 126: 1},
			2: {34: 4},
		},
	}
)

type HandlerBase struct{}

func (e *HandlerBase) Error(msg string, pos Position, symbol string) {
	fmt.Printf("ERROR%s: %s %s\n", pos.String(), msg, symbol)
}

func (h *HandlerBase) Skip(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return Token{}, true
}

func (h *HandlerBase) RegularStart(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switchCond(ConditionREGULAR)

	return NewToken(TagRegularStart, start, end, text), false
}

func (h *HandlerBase) RegularEnd(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switchCond(ConditionINIT)

	return NewToken(TagRegularEnd, start, end, text), false
}

func (h *HandlerBase) RegularNewLine(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switchCond(ConditionINIT)

	return NewToken(TagRegularNewLine, start, end, text), false
}

func (h *HandlerBase) RegularEscapeNewLine(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagRegularEscapeNewLine, start, end, text), false
}

func (h *HandlerBase) RegularEscapeTab(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagRegularEscapeTab, start, end, text), false
}

func (h *HandlerBase) RegularEscapeQota(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagRegularEscapeQota, start, end, text), false
}

func (h *HandlerBase) RegularSymb(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagRegularSymb, start, end, text), false
}

func (h *HandlerBase) StartLiteral(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switchCond(ConditionLITERAL)

	return NewToken(TagStartLiteral, start, end, text), false
}

func (h *HandlerBase) Literal1(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagLiteral1, start, end, text), false
}

func (h *HandlerBase) LiteralEnd(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switchCond(ConditionINIT)

	return NewToken(TagLiteralEnd, start, end, text), false
}

func (h *HandlerBase) LiteralNewLine(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagLiteralNewLine, start, end, text), false
}

func (h *HandlerBase) LiteralChar(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagLiteralChar, start, end, text), false
}

func (h *HandlerBase) Num(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagNum, start, end, text), false
}

func (h *HandlerBase) Any(
	text string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	return NewToken(TagAny, start, end, text), false
}

type EOPTag struct{}

func (e EOPTag) GetTag() string {
	return EOP
}

type Token struct {
	tag    Tag
	coords fragment
	val    string
}

func NewToken(tag Tag, starting, following Position, val string) Token {
	return Token{tag: tag, coords: newFragment(starting, following), val: val}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s: %s", t.tag.GetTag(), t.coords, t.val)
}

func (t Token) Tag() string {
	return t.tag.GetTag()
}

type fragment struct {
	starting  Position
	following Position
}

func newFragment(starting, following Position) fragment {
	return fragment{starting: starting, following: following}
}

func (f fragment) String() string {
	return fmt.Sprintf("%s-%s", f.starting.String(), f.following.String())
}

type Position struct {
	line  int
	pos   int
	index int
	text  []rune
}

func NewPosition(text []rune) Position {
	return Position{text: text, line: 1, pos: 1}
}

func (p *Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.line, p.pos)
}

func (p *Position) cp() int {
	if p.index == len(p.text) {
		return -1
	}
	return int(p.text[p.index])
}

func (p *Position) isNewLine() bool {
	return p.cp() == '\n'
}

func (p *Position) Index() int {
	return p.index
}

func (p *Position) next() Position {
	if p.index < len(p.text) {
		if p.isNewLine() {
			p.line++
			p.pos = 1
		} else {
			p.pos++
		}
		p.index++
	}

	return *p
}

type Scanner struct {
	program      []rune
	handler      LexemHandler
	regexps      map[Condition][]*FiniteState
	unionRegexps map[Condition]*FiniteState
	curPos       Position

	curCondition Condition
}

func NewScanner(program []rune, handler LexemHandler) Scanner {
	regexps := make(map[Condition][]*FiniteState)

	regexps[ConditionINIT] = make([]*FiniteState, 0, 5)
	regexps[ConditionINIT] = append(regexps[ConditionINIT], automataSkip)
	regexps[ConditionINIT] = append(regexps[ConditionINIT], automataRegularStart)
	regexps[ConditionINIT] = append(regexps[ConditionINIT], automataStartLiteral)
	regexps[ConditionINIT] = append(regexps[ConditionINIT], automataNum)
	regexps[ConditionINIT] = append(regexps[ConditionINIT], automataAny)

	regexps[ConditionLITERAL] = make([]*FiniteState, 0, 4)
	regexps[ConditionLITERAL] = append(regexps[ConditionLITERAL], automataLiteral1)
	regexps[ConditionLITERAL] = append(regexps[ConditionLITERAL], automataLiteralEnd)
	regexps[ConditionLITERAL] = append(regexps[ConditionLITERAL], automataLiteralNewLine)
	regexps[ConditionLITERAL] = append(regexps[ConditionLITERAL], automataLiteralChar)

	regexps[ConditionREGULAR] = make([]*FiniteState, 0, 6)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEnd)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularNewLine)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEscapeNewLine)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEscapeTab)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEscapeQota)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularSymb)

	unionRegexps := make(map[Condition]*FiniteState)

	unionRegexps[ConditionINIT] = unionAutomataINIT

	unionRegexps[ConditionLITERAL] = unionAutomataLITERAL

	unionRegexps[ConditionREGULAR] = unionAutomataREGULAR

	return Scanner{
		program:      program,
		handler:      handler,
		regexps:      regexps,
		unionRegexps: unionRegexps,
		curPos:       NewPosition(program),
		curCondition: ConditionINIT,
	}
}

func (s *Scanner) switchCondition(cond Condition) {
	s.curCondition = cond
}

func (s *Scanner) findTokenOneAutomata(
	name string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch s.curCondition {
	case ConditionINIT:
		return s.findTokenOneAutomataINIT(name, start, end, errFunc, switchCond)
	case ConditionLITERAL:
		return s.findTokenOneAutomataLITERAL(name, start, end, errFunc, switchCond)
	case ConditionREGULAR:
		return s.findTokenOneAutomataREGULAR(name, start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findTokenOneAutomataINIT(
	name string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch name {
	case "Skip":
		return s.handler.Skip(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "RegularStart":
		return s.handler.RegularStart(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "StartLiteral":
		return s.handler.StartLiteral(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "Num":
		return s.handler.Num(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "Any":
		return s.handler.Any(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findTokenOneAutomataLITERAL(
	name string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch name {
	case "Literal1":
		return s.handler.Literal1(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "LiteralEnd":
		return s.handler.LiteralEnd(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "LiteralNewLine":
		return s.handler.LiteralNewLine(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "LiteralChar":
		return s.handler.LiteralChar(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findTokenOneAutomataREGULAR(
	name string,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch name {
	case "RegularEnd":
		return s.handler.RegularEnd(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "RegularNewLine":
		return s.handler.RegularNewLine(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "RegularEscapeNewLine":
		return s.handler.RegularEscapeNewLine(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "RegularEscapeTab":
		return s.handler.RegularEscapeTab(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "RegularEscapeQota":
		return s.handler.RegularEscapeQota(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case "RegularSymb":
		return s.handler.RegularSymb(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findToken(
	automata *FiniteState,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch s.curCondition {
	case ConditionINIT:
		return s.findTokenINIT(automata, start, end, errFunc, switchCond)
	case ConditionLITERAL:
		return s.findTokenLITERAL(automata, start, end, errFunc, switchCond)
	case ConditionREGULAR:
		return s.findTokenREGULAR(automata, start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findTokenINIT(
	automata *FiniteState,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch automata {
	case automataSkip:
		return s.handler.Skip(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataRegularStart:
		return s.handler.RegularStart(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataStartLiteral:
		return s.handler.StartLiteral(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataNum:
		return s.handler.Num(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataAny:
		return s.handler.Any(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findTokenLITERAL(
	automata *FiniteState,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch automata {
	case automataLiteral1:
		return s.handler.Literal1(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataLiteralEnd:
		return s.handler.LiteralEnd(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataLiteralNewLine:
		return s.handler.LiteralNewLine(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataLiteralChar:
		return s.handler.LiteralChar(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) findTokenREGULAR(
	automata *FiniteState,
	start, end Position,
	errFunc ErrFunc,
	switchCond SwitchConditionFunc,
) (Token, Continued) {
	switch automata {
	case automataRegularEnd:
		return s.handler.RegularEnd(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataRegularNewLine:
		return s.handler.RegularNewLine(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataRegularEscapeNewLine:
		return s.handler.RegularEscapeNewLine(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataRegularEscapeTab:
		return s.handler.RegularEscapeTab(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataRegularEscapeQota:
		return s.handler.RegularEscapeQota(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	case automataRegularSymb:
		return s.handler.RegularSymb(string(s.program[start.Index():end.Index()]), start, end, errFunc, switchCond)
	}

	return Token{}, true
}

func (s *Scanner) NextTokenOneAutomata() Token {
	for s.curPos.cp() != -1 {
		start := s.curPos.index

		maxRight := 0
		res, name := s.unionRegexps[s.curCondition].FindMatchEndIndexOneAutomata(string(s.program[s.curPos.index:]))
		if res > 0 {
			maxRight = res
		}

		startPos := s.curPos
		var pos Position
		for s.curPos.index != start+maxRight {
			pos = s.curPos
			s.curPos.next()
		}
		pos.index++

		if maxRight == 0 {
			if s.curPos.cp() != -1 {
				s.curPos.next()
			} else {
				break
			}
			s.handler.Error("ERROR: unknown symbol", startPos, string(s.program[start]))
		} else {
			tok, continued := s.findTokenOneAutomata(name, startPos, pos, s.handler.Error, s.switchCondition)
			if !continued {
				return tok
			}
		}
	}

	return NewToken(EOPTag{}, s.curPos, s.curPos, "")
}

func (s *Scanner) NextToken() Token {
	for s.curPos.cp() != -1 {
		start := s.curPos.index

		var maxRightReg *FiniteState
		maxRight := 0

		for _, r := range s.regexps[s.curCondition] {
			res := r.FindMatchEndIndex(string(s.program[s.curPos.index:]))
			if res > maxRight {
				maxRightReg = r
				maxRight = res
			}
		}
		startPos := s.curPos
		var pos Position
		for s.curPos.index != start+maxRight {
			pos = s.curPos
			s.curPos.next()
		}
		pos.index++

		if maxRight == 0 {
			if s.curPos.cp() != -1 {
				s.curPos.next()
			} else {
				break
			}
			s.handler.Error("ERROR: unknown symbol", startPos, string(s.program[start]))
		} else {
			tok, continued := s.findToken(maxRightReg, startPos, pos, s.handler.Error, s.switchCondition)
			if !continued {
				return tok
			}
		}
	}

	return NewToken(EOPTag{}, s.curPos, s.curPos, "")
}
