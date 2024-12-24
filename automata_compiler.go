package main

import (
	"fmt"
	"slices"
)

type AutomataCompiler interface {
	Compile() *FiniteAutomata
}

func (r *RegExpr) Compile() *FiniteAutomata {
	if r.union != nil {
		return r.union.Compile()
	}

	// impossible case
	return nil
}

func (u *Union) Compile() *FiniteAutomata {
	if len(u.concatenations) == 0 {
		return nil
	}
	a := u.concatenations[0].Compile()
	for _, c := range u.concatenations[1:] {
		a.Union(c.Compile())
	}

	return a
}

func (c *Concatenation) Compile() *FiniteAutomata {
	if len(c.basic) == 0 {
		return nil
	}

	a := c.basic[0].Compile()
	for _, b := range c.basic[1:] {
		a.Concat(b.Compile())
	}

	return a
}

func (be *BasicExpr) Compile() *FiniteAutomata {
	a := be.element.Compile()
	if be.op != nil {
		switch be.op.tag {
		case TagStar:
			a.Loop()
		case TagPlus:
			b := a.copy()
			b.Loop()
			a.Concat(b)
		case TagQuestion:
			a.TerminalStates = append(a.TerminalStates, TerminalState{State: 0})
			a.nullable = true
		}
	}

	return a
}

func (e *Element) Compile() *FiniteAutomata {
	if e.group != nil {
		return e.group.Compile()
	}

	if e.set != nil {
		return e.set.Compile()
	}

	if e.character != nil {
		return e.character.Compile()
	}

	if e.escape != nil {
		return e.escape.Compile()
	}

	// impossible case
	return nil
}

func (g *Group) Compile() *FiniteAutomata {
	return g.regExpr.Compile()
}

func (e *Escape) Compile() *FiniteAutomata {
	switch e.base.tok.val {
	case "t":
		return Create([]rune{'\t'}, e.base.pos)
	case "n":
		return Create([]rune{'\n'}, e.base.pos)
	case "r":
		return Create([]rune{'\r'}, e.base.pos)
	case "f":
		return Create([]rune{'\f'}, e.base.pos)
	case "d":
		return Create(genRuneInRange('0', '9'), e.base.pos)
	case "s":
		return Create([]rune{'\n', '\r', '\f', '\t', ' '}, e.base.pos)
	case "w":
		lowerLetters := genRuneInRange('a', 'z')
		lts := append(lowerLetters, genRuneInRange('A', 'Z')...)
		lettersAndNums := append(lts, genRuneInRange('0', '9')...)
		return Create(append(lettersAndNums, '_'), e.base.pos)
	}

	return Create([]rune(e.base.tok.val), e.base.pos)
}

func (s *Set) Compile() *FiniteAutomata {
	if s.positive != nil {
		a := s.positive.Compile()
		return a
	}

	if s.negative != nil {
		letCpy := copyMap(letters)
		flPosCpy := copyMap(flPos)

		a := s.negative.Compile()

		letters = letCpy
		flPos = flPosCpy
		Negate(a, s.pos)
		return a
	}

	// impossible case
	return nil
}

func (s *SetItems) Compile() *FiniteAutomata {
	a := s.item.Compile()

	if s.items != nil {
		b := s.items.Compile()
		a.Union(b)
	}

	return a
}

func (s *SetItem) Compile() *FiniteAutomata {
	if s.rnge != nil {
		return s.rnge.Compile()
	}

	if s.base != nil {
		return s.base.Compile()
	}

	if s.escape != nil {
		return s.escape.Compile()
	}

	// impossible case
	return nil
}

func (r *Range) Compile() *FiniteAutomata {
	startChar := ' '
	if r.startToken != nil {
		startChar = []rune(r.startToken.val)[0]
	}
	if r.startEscape != nil {
		startChar = []rune(r.startEscape.base.tok.val)[0]
	}

	if startChar > []rune(r.end.val)[0] {
		panic(fmt.Sprintf("character range is out of order: %s-%s: ASCII(%s)>ASCII(%s)",
			string(startChar), r.end.val, string(startChar), r.end.val))
	}

	chars := genRuneInRange(startChar, []rune(r.end.val)[0])

	a := Create(chars, r.pos)
	return a
}

func genRuneInRange(startChar, endChar rune) []rune {
	chars := make([]rune, 0, endChar-startChar+1)
	for i := startChar; i <= endChar; i++ {
		chars = append(chars, i)
	}

	return chars
}

var anyRuneNotNL = genAnyRuneNotNL()

//const (
//	minRune = 98
//	maxRune = 100
//)

const runeRangeEnd = 0x10ffff

const printableChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \t\n\r"

//const printableChars = "012abc\n\r"

var printableCharsNoNL = printableChars[:len(printableChars)-2]

func genAnyRuneNotNL() []rune {
	var res []rune
	//for i = minRune; i < '\n'; i++ {
	//	res = append(res, i)
	//}
	//var i rune
	//for i = minRune; i < maxRune; i++ {
	//	res = append(res, i)
	//}
	for _, i := range printableCharsNoNL {
		res = append(res, i)
	}

	return res
}

func genAnyRuneExcept(except []rune) []rune {
	var res []rune
	for _, i := range printableCharsNoNL {
		if slices.Contains(except, i) {
			continue
		}
		res = append(res, i)
	}

	return res
}

func (t *Character) Compile() *FiniteAutomata {
	if t.tok.tag == TagAnyCharacter {
		res := NewAutomata()
		for _, i := range anyRuneNotNL {
			res.Union(Create([]rune{i}, t.pos))
		}

		return res
	}
	a := Create([]rune(t.tok.val), t.pos)

	return a
}
