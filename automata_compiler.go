package main

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

	panic("invalid")
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

	panic("invalid")
}

func (c *Concatenation) Compile() *FiniteState {
	a := c.basic.Compile()
	b := c.simple.Compile()
	a.Append(b)
	return a
}

func (be *BasicExpr) Compile() *FiniteState {
	a := be.element.Compile()
	if be.op != nil {
		switch be.op.tag {
		case TagStar:
			a.Loop()
		case TagPlus:
			b := be.element.Compile()
			b.Loop()
			a.Append(b)
		case TagQuestion:
			a.TerminalStates = append(a.TerminalStates, 0)
		}
	}

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

	panic("invalid")
}

func (g *Group) Compile() *FiniteState {
	return g.regExpr.Compile()
}

func (e *Escape) Compile() *FiniteState {
	switch e.base.val {
	case "t":
		return Create([]rune{'\t'})
	case "n":
		return Create([]rune{'\n'})
	case "r":
		return Create([]rune{'\r'})
	case "f":
		return Create([]rune{'\f'})
	}

	return Create([]rune(e.base.val))
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

	panic("invalid")
}

func (s *SetItems) Compile() *FiniteState {
	a := s.item.Compile()

	if s.items != nil {
		b := s.items.Compile()
		a.Union(b)
	}

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

	panic("invalid")
}

func (r *Range) Compile() *FiniteState {
	var chars []rune
	for i := []rune(r.start.val)[0]; i <= []rune(r.end.base.val)[0]; i++ {
		chars = append(chars, i)
	}

	return Create(chars)
}

func (c *Character) Compile() *FiniteState {
	return c.base.Compile()
}

var anyRuneNotNL = genAnyRuneNotNL()

//const (
//	minRune = 98
//	maxRune = 100
//)

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

func (t *Token) Compile() *FiniteState {
	if t.tag == TagAnyCharacter {
		res := NewAutomata()
		for _, i := range anyRuneNotNL {
			res.Union(Create([]rune{i}))
		}

		return res
	}

	return Create([]rune(t.val))
}
