package main

type Program struct {
	namedRegExprs []NamedRegExpr
	state         *State
	rules         Rules
}

type NamedRegExpr struct {
	name Token
	expr *RegExpr
}

type State struct {
	names []Token
}

type Rules struct {
	ruleArr []Rule
}

type Rule struct {
	startCondition  *StartCondition
	name            Token
	expr            RegExpr
	contin          *Token
	edit            *Token
	switchCondition *SwitchCondition
}

type StartCondition struct {
	condition Token
}

type SwitchCondition struct {
	nextCondition Token
}

type RegExpr struct {
	union *Union
}

type Union struct {
	concatenations []Concatenation
}

type Concatenation struct {
	basic []BasicExpr
}

type BasicExpr struct {
	op      *Token
	element *Element
}

type Element struct {
	// Value     rune
	character *Character
	group     *Group
	set       *Set
	escape    *Escape
}

type Group struct {
	regExpr *RegExpr
}

type Escape struct {
	// character *Character
	base *Character
}

type Set struct {
	positive *SetItems
	negative *SetItems
	pos      int
}

type SetItems struct {
	item  *SetItem
	items *SetItems
}

type SetItem struct {
	rnge   *Range
	base   *Character
	escape *Escape
}

type Range struct {
	startToken  *Token
	startEscape *Escape
	end         *Token
	pos         int
}

type Character struct {
	tok *Token
	pos int
}
