package main

import (
	"fmt"
	"math"
	"os"
)

type AutomataCompiler interface {
	Compile() *FiniteState
}

func (r *RegExpr) Compile() *FiniteState {
	if r.union != nil {
		return r.union.Compile()
	}

	if r.simple != nil {
		return r.simple.Compile()
	}

	// impossible case
	return nil
}

func (u *Union) Compile() *FiniteState {
	a := u.simple.Compile()
	b := u.regex.Compile()
	a.Union(b)
	return a
}

func (s *SimpleExpr) Compile() *FiniteState {
	if s.concatenation != nil {
		return s.concatenation.Compile()
	}

	if s.basic != nil {
		return s.basic.Compile()
	}

	// impossible case
	return nil
}

func (c *Concatenation) Compile() *FiniteState {
	a := c.basic.Compile()

	var rec func(simple *SimpleExpr)
	rec = func(simple *SimpleExpr) {
		if simple.concatenation != nil {
			b := simple.concatenation.basic.Compile()
			a.Append(b)
			simple.lastpos = a.lastpos
			simple.firstpos = a.firstpos
			simple.nullable = a.nullable
			rec(simple.concatenation.simple)
		} else {
			b := simple.basic.Compile()
			a.Append(b)
			simple.lastpos = a.lastpos
			simple.firstpos = a.firstpos
			simple.nullable = a.nullable
		}
	}
	rec(c.simple)

	c.lastpos = a.lastpos
	c.firstpos = a.firstpos
	c.nullable = a.nullable
	return a
}

func (be *BasicExpr) Compile() *FiniteState {
	a := be.element.Compile()
	if be.op != nil {
		switch be.op.tag {
		case TagStar:
			a.ToGraph(os.Stdout)
			a.Loop()
			a.ToGraph(os.Stdout)
		case TagPlus:
			b := a.copy()
			b.Loop()
			a.Append(b)
		case TagQuestion:
			a.TerminalStates = append(a.TerminalStates, TerminalState{State: 0})
			a.nullable = true
		}
	}

	if be.repetition != nil {
		if be.repetition.max < be.repetition.min {
			panic("invalid min and max repetition")
		}
		init := a.copy()
		for i := 0; i < be.repetition.min-2; i++ {
			b := init.copy()
			a.Append(b)
		}
		if be.repetition.max == math.MaxInt {
			b := init.copy()
			a.Append(b)
			b.Loop()
			a.Append(b)
		} else if be.repetition.max == be.repetition.min && be.repetition.min != 1 && be.repetition.min != 0 {
			b := init.copy()
			a.Append(b)
		} else if be.repetition.min != 1 && be.repetition.min != 0 {
			b := init.copy()
			a.Append(b)
			b.Append(init)
			for i := 0; i < be.repetition.max-2; i++ {
				b.Append(init)
				a.Union(b)
			}
		}

		if be.repetition.min == 0 && be.repetition.max == be.repetition.min {
			return NewAutomata()
		}
	}

	be.lastpos = a.lastpos
	be.firstpos = a.firstpos
	be.nullable = a.nullable

	return a
}

func (e *Element) Compile() *FiniteState {
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

func (g *Group) Compile() *FiniteState {
	return g.regExpr.Compile()
}

func (e *Escape) Compile() *FiniteState {
	switch e.base.tok.val {
	case "t":
		return Create([]rune{'\t'}, e.base.pos)
	case "n":
		return Create([]rune{'\n'}, e.base.pos)
		//case "r":
		//	return Create([]rune{'\r'})
		//case "f":
		//	return Create([]rune{'\f'})
		//case "d":
		//	return Create(genRuneInRange('0', '9'))
		//case "s":
		//	return Create([]rune{'\n', '\r', '\f', '\t', ' '})
		//case "w":
		//	lowerLetters := genRuneInRange('a', 'z')
		//	letters := append(lowerLetters, genRuneInRange('A', 'Z')...)
		//	lettersAndNums := append(letters, genRuneInRange('0', '9')...)
		//	return Create(append(lettersAndNums, '_'))
	}

	return Create([]rune(e.base.tok.val), e.base.pos)
}

func (s *Set) Compile() *FiniteState {
	if s.positive != nil {
		return s.positive.Compile()
	}

	if s.negative != nil {
		a := s.negative.Compile()
		a.Negate()
		return a
	}

	// impossible case
	return nil
}

func (s *SetItems) Compile() *FiniteState {
	a := s.item.Compile()

	if s.items != nil {
		b := s.items.Compile()
		a.Union(b)
	}

	s.lastpos = a.lastpos
	s.firstpos = a.firstpos
	s.nullable = a.nullable

	return a
}

func (s *SetItem) Compile() *FiniteState {
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

func (r *Range) Compile() *FiniteState {
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
	r.lastpos = a.lastpos
	r.firstpos = a.firstpos
	r.nullable = a.nullable
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

//func (t *Token) Compile() *FiniteState {
//	if t.tag == TagAnyCharacter {
//		res := NewAutomata()
//		for _, i := range anyRuneNotNL {
//			res.Union(Create([]rune{i}, 0))
//		}
//
//		return res
//	}
//
//	return Create([]rune(t.val), 0)
//}

func (t *Character) Compile() *FiniteState {
	if t.tok.tag == TagAnyCharacter {
		res := NewAutomata()
		for _, i := range anyRuneNotNL {
			res.Union(Create([]rune{i}, t.pos))
		}

		t.lastpos = res.lastpos
		t.firstpos = res.firstpos
		t.nullable = res.nullable

		return res
	}
	a := Create([]rune(t.tok.val), t.pos)
	t.lastpos = a.lastpos
	t.firstpos = a.firstpos
	t.nullable = a.nullable

	return a
}
