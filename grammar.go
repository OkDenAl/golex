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

	nullable bool
	firstpos []int
	lastpos  []int
}

// SimpleExpr ::= Concatenation | BasicExpr
type SimpleExpr struct {
	concatenation *Concatenation
	basic         *BasicExpr

	nullable bool
	firstpos []int
	lastpos  []int
}

// Concatenation ::= SimpleExpr BasicExpr
type Concatenation struct {
	simple *SimpleExpr
	basic  *BasicExpr

	nullable bool
	firstpos []int
	lastpos  []int
}

// BasicExpr ::= Element ("*"|"+"|"?"|Repetition)?
type BasicExpr struct {
	op         *Token
	repetition *Repetition
	element    *Element

	nullable bool
	firstpos []int
	lastpos  []int
}

// Repetition      ::= "{" Number ("}" | "," ("}" | Number "}"))
type Repetition struct {
	max int
	min int

	nullable bool
	firstpos []int
	lastpos  []int
}

// Element         ::= Group | Set | Escape | ValidIndependentCharacter
type Element struct {
	// Value     rune
	character *Character
	group     *Group
	set       *Set
	escape    *Escape

	nullable bool
	firstpos []int
	lastpos  []int
}

// Group ::= (RegExpr)
type Group struct {
	regExpr *RegExpr

	nullable bool
	firstpos []int
	lastpos  []int
}

// Escape ::= "\" Character
type Escape struct {
	// character *Character
	base *Character

	nullable bool
	firstpos []int
	lastpos  []int
}

// Set ::= "[" ("^")? SetItems "]"
type Set struct {
	positive *SetItems
	negative *SetItems

	nullable bool
	firstpos []int
	lastpos  []int
}

// SetItems ::= SetItem SetItems
type SetItems struct {
	item  *SetItem
	items *SetItems

	nullable bool
	firstpos []int
	lastpos  []int
}

// SetItem ::= Range | Character
type SetItem struct {
	rnge   *Range
	base   *Character
	escape *Escape

	nullable bool
	firstpos []int
	lastpos  []int
}

// Range ::= Character "-" Character
type Range struct {
	startToken  *Token
	startEscape *Escape
	end         *Token
	pos         int

	nullable bool
	firstpos []int
	lastpos  []int
}

type Character struct {
	tok *Token
	pos int

	nullable bool
	firstpos []int
	lastpos  []int
}
