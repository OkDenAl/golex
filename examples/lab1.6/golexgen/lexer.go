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
	Skip(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularStart(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEnd(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularNewLine(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEscapeNewLine(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEscapeTab(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularEscapeQota(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	RegularSymb(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	StartLiteral(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	Literal1(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	LiteralEnd(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	LiteralNewLine(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	LiterlaChar(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	Num(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
	Any(text []rune, start, end Position, errFunc ErrFunc, switchCond SwitchConditionFunc) (Token, Continued)
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
	TagLiterlaChar
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
		TagLiterlaChar:          "LiterlaChar",
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
	TerminalStates []int
	Transitions    map[int]map[rune]int
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

var (
	automataSkip *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{4, 3, 1, 2, 5, 6},
		Transitions: map[int]map[rune]int{
			0: {9: 2, 10: 1, 32: 3},
			1: {9: 5, 10: 4, 32: 6},
			2: {9: 5, 10: 4, 32: 6},
			3: {9: 5, 10: 4, 32: 6},
			4: {9: 5, 10: 4, 32: 6},
			5: {9: 5, 10: 4, 32: 6},
			6: {9: 5, 10: 4, 32: 6},
		},
	}
	automataRegularStart *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{1},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
		},
	}
	automataRegularEnd *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{1},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
		},
	}
	automataRegularNewLine *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{1},
		Transitions: map[int]map[rune]int{
			0: {10: 1},
		},
	}
	automataRegularEscapeNewLine *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{2},
		Transitions: map[int]map[rune]int{
			0: {92: 1},
			1: {110: 2},
		},
	}
	automataRegularEscapeTab *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{2},
		Transitions: map[int]map[rune]int{
			0: {92: 1},
			1: {116: 2},
		},
	}
	automataRegularEscapeQota *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{2},
		Transitions: map[int]map[rune]int{
			0: {92: 1},
			1: {34: 2},
		},
	}
	automataRegularSymb *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{17, 92, 82, 7, 79, 74, 65, 9, 32, 28, 84, 81, 0, 29, 88, 62, 10, 34, 86, 2, 30, 76, 38, 33, 55, 51, 27, 68, 24, 72, 96, 21, 20, 25, 71, 91, 93, 73, 52, 85, 87, 43, 56, 67, 57, 58, 6, 22, 59, 8, 53, 75, 97, 44, 31, 77, 95, 19, 18, 66, 78, 11, 5, 3, 48, 40, 13, 15, 94, 41, 35, 23, 70, 42, 45, 14, 4, 61, 39, 60, 63, 26, 69, 12, 47, 54, 49, 46, 37, 90, 80, 50, 36, 16, 83, 64, 89},
		Transitions: map[int]map[rune]int{
			0: {9: 97, 32: 96, 33: 64, 34: 65, 35: 66, 36: 67, 37: 68, 38: 69, 39: 70, 40: 71, 41: 72, 42: 73, 43: 74, 44: 75, 45: 76, 46: 77, 47: 78, 48: 2, 49: 3, 50: 4, 51: 5, 52: 6, 53: 7, 54: 8, 55: 9, 56: 10, 57: 11, 58: 79, 59: 80, 60: 81, 61: 82, 62: 83, 63: 84, 64: 85, 65: 38, 66: 39, 67: 40, 68: 41, 69: 42, 70: 43, 71: 44, 72: 45, 73: 46, 74: 47, 75: 48, 76: 49, 77: 50, 78: 51, 79: 52, 80: 53, 81: 54, 82: 55, 83: 56, 84: 57, 85: 58, 86: 59, 87: 60, 88: 61, 89: 62, 90: 63, 91: 86, 92: 87, 93: 88, 94: 89, 95: 90, 96: 91, 97: 12, 98: 13, 99: 14, 100: 15, 101: 16, 102: 17, 103: 18, 104: 19, 105: 20, 106: 21, 107: 22, 108: 23, 109: 24, 110: 25, 111: 26, 112: 27, 113: 28, 114: 29, 115: 30, 116: 31, 117: 32, 118: 33, 119: 34, 120: 35, 121: 36, 122: 37, 123: 92, 124: 93, 125: 94, 126: 95},
		},
	}
	automataStartLiteral *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{2},
		Transitions: map[int]map[rune]int{
			0: {64: 1},
			1: {34: 2},
		},
	}
	automataLiteral1 *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{2},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
			1: {34: 2},
		},
	}
	automataLiteralEnd *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{1},
		Transitions: map[int]map[rune]int{
			0: {34: 1},
		},
	}
	automataLiteralNewLine *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{1},
		Transitions: map[int]map[rune]int{
			0: {10: 1},
		},
	}
	automataLiterlaChar *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{61, 43, 76, 51, 90, 42, 63, 73, 64, 40, 57, 92, 32, 31, 2, 97, 81, 5, 75, 95, 30, 46, 47, 13, 11, 58, 52, 70, 54, 21, 7, 91, 10, 45, 19, 53, 79, 33, 14, 50, 26, 12, 65, 93, 84, 4, 89, 39, 8, 44, 17, 71, 18, 15, 0, 56, 62, 60, 88, 36, 6, 34, 55, 22, 20, 23, 59, 66, 72, 87, 74, 16, 28, 24, 29, 78, 41, 96, 83, 77, 94, 49, 86, 68, 37, 69, 9, 67, 3, 38, 25, 85, 27, 48, 80, 82, 35},
		Transitions: map[int]map[rune]int{
			0: {9: 97, 32: 96, 33: 64, 34: 65, 35: 66, 36: 67, 37: 68, 38: 69, 39: 70, 40: 71, 41: 72, 42: 73, 43: 74, 44: 75, 45: 76, 46: 77, 47: 78, 48: 2, 49: 3, 50: 4, 51: 5, 52: 6, 53: 7, 54: 8, 55: 9, 56: 10, 57: 11, 58: 79, 59: 80, 60: 81, 61: 82, 62: 83, 63: 84, 64: 85, 65: 38, 66: 39, 67: 40, 68: 41, 69: 42, 70: 43, 71: 44, 72: 45, 73: 46, 74: 47, 75: 48, 76: 49, 77: 50, 78: 51, 79: 52, 80: 53, 81: 54, 82: 55, 83: 56, 84: 57, 85: 58, 86: 59, 87: 60, 88: 61, 89: 62, 90: 63, 91: 86, 92: 87, 93: 88, 94: 89, 95: 90, 96: 91, 97: 12, 98: 13, 99: 14, 100: 15, 101: 16, 102: 17, 103: 18, 104: 19, 105: 20, 106: 21, 107: 22, 108: 23, 109: 24, 110: 25, 111: 26, 112: 27, 113: 28, 114: 29, 115: 30, 116: 31, 117: 32, 118: 33, 119: 34, 120: 35, 121: 36, 122: 37, 123: 92, 124: 93, 125: 94, 126: 95},
		},
	}
	automataNum *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{3, 1, 2},
		Transitions: map[int]map[rune]int{
			0: {48: 1, 49: 2},
			2: {49: 3},
			3: {49: 3},
		},
	}
	automataAny *FiniteState = &FiniteState{
		CurrentState:   0,
		TerminalStates: []int{80, 78, 28, 16, 41, 46, 73, 31, 94, 48, 30, 83, 54, 77, 72, 75, 76, 9, 63, 25, 26, 49, 11, 97, 58, 10, 47, 20, 23, 68, 18, 50, 12, 0, 59, 92, 34, 70, 87, 36, 62, 37, 51, 66, 56, 95, 74, 14, 39, 3, 44, 96, 81, 90, 35, 19, 93, 24, 6, 45, 40, 33, 13, 27, 89, 88, 8, 32, 17, 2, 69, 79, 5, 7, 67, 15, 21, 82, 52, 53, 84, 61, 42, 60, 86, 55, 29, 57, 91, 43, 22, 64, 85, 38, 65, 71, 4},
		Transitions: map[int]map[rune]int{
			0: {9: 97, 32: 96, 33: 64, 34: 65, 35: 66, 36: 67, 37: 68, 38: 69, 39: 70, 40: 71, 41: 72, 42: 73, 43: 74, 44: 75, 45: 76, 46: 77, 47: 78, 48: 2, 49: 3, 50: 4, 51: 5, 52: 6, 53: 7, 54: 8, 55: 9, 56: 10, 57: 11, 58: 79, 59: 80, 60: 81, 61: 82, 62: 83, 63: 84, 64: 85, 65: 38, 66: 39, 67: 40, 68: 41, 69: 42, 70: 43, 71: 44, 72: 45, 73: 46, 74: 47, 75: 48, 76: 49, 77: 50, 78: 51, 79: 52, 80: 53, 81: 54, 82: 55, 83: 56, 84: 57, 85: 58, 86: 59, 87: 60, 88: 61, 89: 62, 90: 63, 91: 86, 92: 87, 93: 88, 94: 89, 95: 90, 96: 91, 97: 12, 98: 13, 99: 14, 100: 15, 101: 16, 102: 17, 103: 18, 104: 19, 105: 20, 106: 21, 107: 22, 108: 23, 109: 24, 110: 25, 111: 26, 112: 27, 113: 28, 114: 29, 115: 30, 116: 31, 117: 32, 118: 33, 119: 34, 120: 35, 121: 36, 122: 37, 123: 92, 124: 93, 125: 94, 126: 95},
		},
	}
)

type ErrHandlerBase struct{}

func (e *ErrHandlerBase) Error(msg string, pos Position, symbol string) {
	fmt.Printf("ERROR%s: %s %s\n", pos.String(), msg, symbol)
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
	program []rune
	handler LexemHandler
	regexps map[Condition][]*FiniteState
	curPos  Position

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
	regexps[ConditionLITERAL] = append(regexps[ConditionLITERAL], automataLiterlaChar)

	regexps[ConditionREGULAR] = make([]*FiniteState, 0, 6)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEnd)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularNewLine)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEscapeNewLine)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEscapeTab)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularEscapeQota)
	regexps[ConditionREGULAR] = append(regexps[ConditionREGULAR], automataRegularSymb)

	return Scanner{program: program, handler: handler, regexps: regexps, curPos: NewPosition(program), curCondition: ConditionINIT}
}

func (s *Scanner) switchCondition(cond Condition) {
	s.curCondition = cond
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
		return s.handler.Skip(s.program, start, end, errFunc, switchCond)
	case automataRegularStart:
		return s.handler.RegularStart(s.program, start, end, errFunc, switchCond)
	case automataStartLiteral:
		return s.handler.StartLiteral(s.program, start, end, errFunc, switchCond)
	case automataNum:
		return s.handler.Num(s.program, start, end, errFunc, switchCond)
	case automataAny:
		return s.handler.Any(s.program, start, end, errFunc, switchCond)
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
		return s.handler.Literal1(s.program, start, end, errFunc, switchCond)
	case automataLiteralEnd:
		return s.handler.LiteralEnd(s.program, start, end, errFunc, switchCond)
	case automataLiteralNewLine:
		return s.handler.LiteralNewLine(s.program, start, end, errFunc, switchCond)
	case automataLiterlaChar:
		return s.handler.LiterlaChar(s.program, start, end, errFunc, switchCond)
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
		return s.handler.RegularEnd(s.program, start, end, errFunc, switchCond)
	case automataRegularNewLine:
		return s.handler.RegularNewLine(s.program, start, end, errFunc, switchCond)
	case automataRegularEscapeNewLine:
		return s.handler.RegularEscapeNewLine(s.program, start, end, errFunc, switchCond)
	case automataRegularEscapeTab:
		return s.handler.RegularEscapeTab(s.program, start, end, errFunc, switchCond)
	case automataRegularEscapeQota:
		return s.handler.RegularEscapeQota(s.program, start, end, errFunc, switchCond)
	case automataRegularSymb:
		return s.handler.RegularSymb(s.program, start, end, errFunc, switchCond)
	}

	return Token{}, true
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