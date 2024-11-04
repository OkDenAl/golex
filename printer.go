package main

import "fmt"

const (
	mainIdent  = " "
	blockIdent = " "
)

type AstPrinter interface {
	Print(indent string)
}

func (r *RegExpr) Print(indent string) {
	if r.union != nil {
		r.union.Print(indent + mainIdent)
	}

	if r.simple != nil {
		r.simple.Print(indent + mainIdent)
	}
}

func (u *Union) Print(indent string) {
	u.simple.Print(indent + mainIdent)
	fmt.Println(indent + "|")
	u.regex.Print(indent + mainIdent)
}

func (s *SimpleExpr) Print(indent string) {
	if s.concatenation != nil {
		s.concatenation.Print(indent + mainIdent)
	}

	if s.basic != nil {
		s.basic.Print(indent + mainIdent)
	}
}

func (c *Concatenation) Print(indent string) {
	c.basic.Print(indent + mainIdent)
	c.simple.Print(indent + mainIdent)
}

func (be *BasicExpr) Print(indent string) {
	if be.op != nil {
		be.op.Print(indent + mainIdent)
	}
	be.element.Print(indent + mainIdent)
}

func (e *Element) Print(indent string) {
	if e.group != nil {
		e.group.Print(indent + mainIdent)
	}

	if e.set != nil {
		e.set.Print(indent + mainIdent)
	}

	if e.character != nil {
		e.character.Print(indent + mainIdent)
	}

	if e.escape != nil {
		e.escape.Print(indent + mainIdent)
	}
}

func (g *Group) Print(indent string) {
	g.regExpr.Print(indent + mainIdent)
}

func (e *Escape) Print(indent string) {
	e.base.Print(indent + mainIdent)
}

func (s *Set) Print(indent string) {
	if s.positive != nil {
		s.positive.Print(indent + mainIdent)
	}

	if s.negative != nil {
		s.negative.Print(indent + mainIdent)
	}
}

func (s *SetItems) Print(indent string) {
	s.item.Print(indent + mainIdent)

	if s.items != nil {
		s.items.Print(indent + blockIdent)
	}
}

func (s *SetItem) Print(indent string) {
	if s.rnge != nil {
		s.rnge.Print(indent + mainIdent)
	}

	if s.base != nil {
		s.base.Print(indent + blockIdent)
	}

	if s.escape != nil {
		s.escape.Print(indent + mainIdent)
	}
}

func (r *Range) Print(indent string) {
	if r.startToken != nil {
		r.startToken.Print(indent)
	}
	if r.startEscape != nil {
		r.startEscape.Print(indent + mainIdent)
	}

	r.end.Print(indent + mainIdent)
}

func (t *Token) Print(indent string) {
	fmt.Println(indent + t.String())
}
