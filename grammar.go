package main

type NamedRegExpr struct {
	name *Token
	expr *RegExpr
	nl   *NewLine
}

type State struct {
	names []*Token
}

type NewLine struct {
	base Token
}

type Rules struct {
	rules []*Rule
}

type Rule struct {
	name *Token
	expr *RegExpr
}

// RegExpr ::= Union | SimpleExpr
type RegExpr struct {
	union  *Union
	simple *SimpleExpr
}

// Union ::= RegExpr "|" SimpleExpr
type Union struct {
	regex  *RegExpr
	simple *SimpleExpr
}

// SimpleExpr ::= Concatenation | BasicExpr
type SimpleExpr struct {
	concatenation *Concatenation
	basic         *BasicExpr
}

// Concatenation ::= SimpleExpr BasicExpr
type Concatenation struct {
	simple *SimpleExpr
	basic  *BasicExpr
}

// BasicExpr ::= Element ("*"|"+"|"?")?
type BasicExpr struct {
	op      *Token
	element *Element
}

// Element ::= Character | Group | Set
type Element struct {
	// Value     rune
	character *Token
	group     *Group
	set       *Set
	escape    *Escape
}

// Group ::= (RegExpr)
type Group struct {
	regExpr *RegExpr
}

// Escape ::= "\" Character
type Escape struct {
	// character *Character
	base *Token
}

// Set ::= "[" ("^")? SetItems "]"
type Set struct {
	positive *SetItems
	negative *SetItems
}

// SetItems ::= SetItem SetItems
type SetItems struct {
	item  *SetItem
	items *SetItems
}

// SetItem ::= Range | Character
type SetItem struct {
	rnge   *Range
	base   *Token
	escape *Escape
}

// Range ::= Character "-" Character
type Range struct {
	start *Token
	end   *Character
}

// Character ::= literal character
type Character struct {
	// Value rune
	base *Token
}
