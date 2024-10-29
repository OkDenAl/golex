package main

import "fmt"

//go:generate stringer -type=DomainTag
type DomainTag int

const (
	TagCode DomainTag = iota
	TagRulesMarker
	TagStateMarker
	TagRuleName
	Tag
	TagRegexp
	TagNL
	TagCharacter
	TagOpenParen
	TagCloseParen
	TagOpenBracket
	TagCloseBracket
	TagOpenBrace
	TagCloseBrace
	TagStar
	TagPlus
	TagQuestion
	TagCaret
	TagEscape
	TagPipe
	TagDash
	TagComma
	TagAnyCharacter
	TagErr
	TagEOP
	TagRegularMarker
)

type Token struct {
	tag    DomainTag
	coords Fragment
	val    string
}

func NewToken(tag DomainTag, starting, following Position, val string) Token {
	return Token{tag: tag, coords: NewFragment(starting, following), val: val}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s: %s", t.tag, t.coords, t.val)
}

func (t Token) Tag() DomainTag {
	return t.tag
}

var tagToString = map[DomainTag]string{}
