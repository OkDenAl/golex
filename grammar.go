package main

// Program         ::= (NamedRegExpr)* (NewLine)* (State)? (NewLine)* Rules
type Program struct {
	namedRegExprs []NamedRegExpr
	state         *State
	rules         Rules
}

type NamedRegExpr struct {
	name Token
	expr *RegExpr
}

// State           ::= "%x" (StateName)+ NewLine
type State struct {
	names []Token
}

// Rules           ::= "%%" (NewLine)+ (Rule)+ (NewLine)+ "%%"
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

// Element         ::= Group | Set | Escape | ValidIndependentCharacter
type Element struct {
	// Value     rune
	character *Character
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
	base *Character
}

// Set ::= "[" ("^")? SetItems "]"
type Set struct {
	positive *SetItems
	negative *SetItems
	pos      int
}

// SetItems ::= SetItem SetItems
type SetItems struct {
	item  *SetItem
	items *SetItems
}

// SetItem ::= Range | Character
type SetItem struct {
	rnge   *Range
	base   *Character
	escape *Escape
}

// Range ::= Character "-" Character
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
