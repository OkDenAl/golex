package main

import "fmt"

type Parser struct {
	cursor int
	tokens []Token

	startConditionNames map[string]struct{}
	ruleNames           map[string]Token
}

func New(tokens []Token) *Parser {
	return &Parser{
		cursor:              0,
		tokens:              tokens,
		startConditionNames: map[string]struct{}{},
		ruleNames:           map[string]Token{},
	}
}

// Parse a list of tokens into an AST Tree
func (p *Parser) Parse() (Program, error) {
	return p.program()
}

// Program         ::= (NamedRegExpr)* (NewLine)* (State)? (NewLine)* Rules
func (p *Parser) program() (Program, error) {
	var namedExpr []NamedRegExpr
	for p.tokens[p.cursor].Tag() == TagName {
		expr, err := p.namedRegExpr()
		if err != nil {
			return Program{}, err
		}
		namedExpr = append(namedExpr, expr)
	}

	for p.tokens[p.cursor].tag == TagNL {
		p.mustExpectTag(TagNL)
	}

	var (
		st  *State
		err error
	)
	if p.tokens[p.cursor].Tag() == TagStateMarker {
		st, err = p.state()
		if err != nil {
			return Program{}, err
		}
	}

	for p.tokens[p.cursor].tag == TagNL {
		p.mustExpectTag(TagNL)
	}

	rules, err := p.rules()
	if err != nil {
		return Program{}, err
	}

	return Program{namedRegExprs: namedExpr, state: st, rules: rules}, nil
}

// NamedRegExpr    ::= RegularName "/" RegExpr "/" NewLine
func (p *Parser) namedRegExpr() (NamedRegExpr, error) {
	name := p.mustExpectTag(TagName)
	p.mustExpectTag(TagRegularMarker)
	reset := p.reset()
	expr, ok := p.regExpr()
	if !ok {
		reset()
		return NamedRegExpr{}, fmt.Errorf("parse error: failed to parse regular expr starts with %s", p.tokens[p.cursor].String())
	}
	p.mustExpectTag(TagRegularMarker)
	p.mustExpectTag(TagNL)

	return NamedRegExpr{name: name, expr: expr}, nil
}

// State           ::= "%x" (Name)+ NewLine
func (p *Parser) state() (*State, error) {
	p.mustExpectTag(TagStateMarker)

	var names []Token
	name := p.mustExpectTag(TagName)
	names = append(names, name)
	p.startConditionNames[name.val] = struct{}{}
	for p.tokens[p.cursor].tag == TagName {
		name = p.mustExpectTag(TagName)
		names = append(names, name)
		p.startConditionNames[name.val] = struct{}{}
	}

	p.mustExpectTag(TagNL)

	return &State{names: names}, nil
}

// Rules           ::= "%%" (NewLine)+ (Rule)+ "%%"
func (p *Parser) rules() (Rules, error) {
	p.mustExpectTag(TagRulesMarker)
	p.mustExpectTag(TagNL)
	for p.tokens[p.cursor].tag == TagNL {
		p.mustExpectTag(TagNL)
	}

	var rules []Rule
	for p.tokens[p.cursor].tag == TagRegularMarker || p.tokens[p.cursor].tag == TagOpenStartCondition {
		rule, err := p.rule()
		if err != nil {
			return Rules{}, err
		}
		rules = append(rules, rule)
	}

	p.mustExpectTag(TagRulesMarker)

	return Rules{ruleArr: rules}, nil
}

// Rule            ::= (StartCondition)? "/" RegExpr "/"  Name (SwitchCondition)? (NewLine)+
func (p *Parser) rule() (Rule, error) {
	var startCond *StartCondition
	if p.tokens[p.cursor].Tag() == TagOpenStartCondition {
		startCond = p.startCondition()
	}

	p.mustExpectTag(TagRegularMarker)
	reset := p.reset()
	expr, ok := p.regExpr()
	if !ok {
		reset()
		return Rule{}, fmt.Errorf("parse error: failed to parse regular expr starts with %s", p.tokens[p.cursor].String())
	}
	p.mustExpectTag(TagRegularMarker)

	token := p.mustExpectTag(TagName)
	if v, ok := p.ruleNames[token.val]; ok {
		panic("Error: duplicate: " + token.String() + "\n      first declaration was here: " + v.String())
	}
	p.ruleNames[token.val] = token

	var switchCond *SwitchCondition
	if p.tokens[p.cursor].Tag() == TagBegin {
		switchCond = p.switchCondition()
	}

	p.mustExpectTag(TagNL)
	for p.tokens[p.cursor].tag == TagNL {
		p.mustExpectTag(TagNL)
	}

	return Rule{startCondition: startCond, expr: *expr, name: token, switchCondition: switchCond}, nil
}

// StartCondition ::= "<" Name ">"
func (p *Parser) startCondition() *StartCondition {
	p.mustExpectTag(TagOpenStartCondition)
	token := p.mustExpectTag(TagName)
	p.mustExpectTag(TagCloseStartCondition)

	return &StartCondition{condition: token}
}

// SwitchCondition ::= "BEGIN" "(" Name ")"
func (p *Parser) switchCondition() *SwitchCondition {
	p.mustExpectTag(TagBegin)
	p.mustExpectTag(TagDefaultOpenBracket)
	token := p.mustExpectTag(TagName)
	p.mustExpectTag(TagDefaultCloseBracket)

	return &SwitchCondition{nextCondition: token}
}

// RegExpr         ::= Union | SimpleExpr
func (p *Parser) regExpr() (*RegExpr, bool) {
	reset := p.reset()
	union, ok := p.union()
	if ok {
		return &RegExpr{union: union}, true
	}

	reset()
	simple, ok := p.simpleExpr()
	if ok {
		return &RegExpr{simple: simple}, true
	}

	return nil, false
}

// Union           ::= SimpleExpr "|" RegExpr
func (p *Parser) union() (*Union, bool) {
	simple, ok := p.simpleExpr()
	if !ok {
		return nil, false
	}

	if _, ok = p.expectTags(TagPipe); !ok {
		return nil, false
	}

	regex, ok := p.regExpr()
	if !ok {
		return nil, false
	}

	return &Union{regex, simple}, true
}

// SimpleExpr      ::= Concatenation | BasicExpr
func (p *Parser) simpleExpr() (*SimpleExpr, bool) {
	concatenation, ok := p.concatenation()
	if ok {
		return &SimpleExpr{concatenation: concatenation}, true
	}

	basic, ok := p.basicExpr()
	if ok {
		return &SimpleExpr{basic: basic}, true
	}

	return nil, false
}

// Concatenation   ::= BasicExpr SimpleExpr
func (p *Parser) concatenation() (*Concatenation, bool) {
	reset := p.reset()

	basic, ok := p.basicExpr()

	if !ok {
		reset()
		return nil, false
	}

	simple, ok := p.simpleExpr()

	if !ok {
		reset()
		return nil, false
	}

	return &Concatenation{simple, basic}, true

}

// BasicExpr       ::= Element ("*"|"+"|"?")?
func (p *Parser) basicExpr() (*BasicExpr, bool) {
	base, ok := p.element()
	if !ok {
		return nil, false
	}

	token, ok := p.expectTags(TagStar, TagPlus, TagQuestion)
	if ok {
		return &BasicExpr{element: base, op: &token}, true
	}

	return &BasicExpr{element: base}, true
}

// Element         ::= Group | Set | Escape | AnyCharacter
func (p *Parser) element() (*Element, bool) {
	group, ok := p.group()
	if ok {
		return &Element{group: group}, true
	}

	set, ok := p.set()
	if ok {
		return &Element{set: set}, true
	}

	escape, ok := p.escape()
	if ok {
		return &Element{escape: escape}, true
	}

	character, ok := p.validIndependentCharacter()
	if ok {
		return &Element{character: character}, true
	}

	return nil, false
}

// Group           ::= "(" RegExpr ")"
func (p *Parser) group() (*Group, bool) {
	if _, ok := p.expectTags(TagOpenParen); !ok {
		return nil, false
	}
	reset := p.reset()
	defer p.mustExpectTag(TagCloseParen)
	regex, ok := p.regExpr()
	if !ok {
		reset()
		return nil, false
	}
	return &Group{regex}, true
}

// Escape          ::= "\" EscapeCharacter
func (p *Parser) escape() (*Escape, bool) {
	if _, ok := p.expectTags(TagEscape); !ok {
		return nil, false
	}

	base, ok := p.escapeCharacter()
	if !ok {
		panic("Escape: no character")
		return nil, false
	}

	return &Escape{base}, true

}

// Set             ::= "[" ("^")? SetItems "]"
func (p *Parser) set() (*Set, bool) {
	reset := p.reset()
	if _, ok := p.expectTags(TagOpenBracket); !ok {
		return nil, false
	}

	var set *Set
	if _, ok := p.expectTags(TagCaret); !ok {
		positive, ok := p.setItems(true)
		if ok {
			set = &Set{positive: positive}
		}
	} else {
		negative, ok := p.setItems(true)
		if ok {
			set = &Set{negative: negative}
		}
	}

	if set != nil {
		p.mustExpectTag(TagCloseBracket)
		return set, true
	}

	reset()
	return nil, false

}

// SetItems        ::= SetItem (SetItem)*
func (p *Parser) setItems(isFirst bool) (*SetItems, bool) {
	item, ok := p.setItem(isFirst)
	if !ok {
		return nil, false
	}
	isFirst = false

	items, ok := p.setItems(isFirst)

	return &SetItems{item: item, items: items}, true

}

// SetItem         ::= Range | Escape | SetCharacter
func (p *Parser) setItem(isFirst bool) (*SetItem, bool) {
	reset := p.reset()
	rnge, ok := p.rangeExpr()
	if ok {
		return &SetItem{rnge: rnge}, true
	}

	reset()
	escape, ok := p.escape()
	if ok {
		return &SetItem{escape: escape}, true
	}

	reset()
	token, ok := p.setCharacter(isFirst)
	if ok {
		return &SetItem{base: token}, true
	}

	reset()
	return nil, false
}

// Range           ::=  (Escape | RangeStartCharacter) "-" RangeEndCharacter
func (p *Parser) rangeExpr() (*Range, bool) {
	var (
		startToken  *Token
		startEscape *Escape
	)
	reset := p.reset()
	escape, ok := p.escape()
	if ok {
		startEscape = escape
	} else {
		reset()
		startToken, ok = p.rangeStartCharacter()
		if !ok {
			return nil, false
		}
	}

	if _, ok = p.expectTags(TagDash); !ok {
		return nil, false
	}

	end, ok := p.rangeEndCharacter()
	if !ok {
		reset()
		panic("unexpected token " + p.tokens[p.cursor].String())
	}

	return &Range{startToken: startToken, startEscape: startEscape, end: end}, true
}

// ValidIndependentCharacter ::= [^()|/]
func (p *Parser) validIndependentCharacter() (*Token, bool) {
	token, ok := p.nextToken()
	if !ok {
		return nil, false
	}

	if token.Tag() == TagRegularMarker || token.Tag() == TagOpenParen || token.Tag() == TagCloseParen || token.Tag() == TagPipe {
		return nil, false
	}

	return &token, true
}

// EscapeCharacter ::= .
func (p *Parser) escapeCharacter() (*Token, bool) {
	token, ok := p.nextToken()
	if !ok {
		return nil, false
	}

	if token.Tag() == TagRegularMarker {
		return nil, false
	}

	if token.Tag() == TagAnyCharacter {
		token.tag = TagCharacter
	}

	return &token, true
}

// RangeStartCharacter ::= .
func (p *Parser) rangeStartCharacter() (*Token, bool) {
	token, ok := p.nextToken()
	if !ok {
		return nil, false
	}

	if token.Tag() == TagRegularMarker {
		return nil, false
	}

	if token.Tag() == TagAnyCharacter {
		token.tag = TagCharacter
	}

	return &token, true
}

// RangeStartCharacter ::= .
func (p *Parser) rangeEndCharacter() (*Token, bool) {
	token, ok := p.nextToken()
	if !ok {
		return nil, false
	}

	if token.Tag() == TagRegularMarker {
		return nil, false
	}

	if token.Tag() == TagAnyCharacter {
		token.tag = TagCharacter
	}

	return &token, true
}

// SetCharacter ::= .[^]]*
func (p *Parser) setCharacter(isFirst bool) (*Token, bool) {
	token, ok := p.nextToken()
	if !ok {
		return nil, false
	}

	if token.Tag() == TagRegularMarker {
		return nil, false
	}

	if token.Tag() == TagAnyCharacter {
		token.tag = TagCharacter
	}

	if token.Tag() != TagCloseBracket || isFirst {
		return &token, true
	}

	return nil, false
}

func (p *Parser) nextToken() (Token, bool) {
	if p.cursor == len(p.tokens) {
		if len(p.tokens) == 0 {
			return Token{tag: TagCharacter, val: " "}, false
		}
		return p.tokens[p.cursor-1], false
	}

	token := p.tokens[p.cursor]
	p.cursor++

	return token, true
}

func (p *Parser) mustExpectTag(tag DomainTag) Token {
	reset := p.reset()
	token, ok := p.nextToken()
	if !ok {
		panic(fmt.Sprintf("parse error: expected %s, but got EOF", tag))
	}

	if token.Tag() == tag {
		return token
	}

	reset()
	panic(fmt.Sprintf("parse error: expected %s, but got %s", tag, p.tokens[p.cursor].String()))
}

func (p *Parser) expectTags(tags ...DomainTag) (Token, bool) {
	reset := p.reset()
	token, ok := p.nextToken()
	if !ok {
		reset()
		return Token{}, false
	}

	for _, tag := range tags {
		if token.Tag() == tag {
			return token, true
		}
	}

	reset()
	return Token{}, false
}

func (p *Parser) reset() func() {
	cursor := p.cursor
	return func() {
		p.cursor = cursor
	}
}
