// Code generated by golex utility;
// This code is present a default handling of tokens;
// YOU CAN EDIT IT IF YOU NEED.
package main

import (
	"fmt"
	"golex/examples/lab1.6/golexgen"
	"os"
)

type Tag string

const TagString = Tag("STRING")

func (t Tag) GetTag() string {
	return string(t)
}

type Handler struct {
	str   string
	start golexgen.Position
	golexgen.ErrHandlerBase
}

func (h *Handler) Skip(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	return golexgen.Token{}, true
}

func (h *Handler) RegularStart(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	switchCond(golexgen.ConditionREGULAR)
	h.start = start

	return golexgen.Token{}, true
}

func (h *Handler) RegularEnd(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	switchCond(golexgen.ConditionINIT)

	return golexgen.NewToken(
		TagString,
		h.start, end,
		h.str,
	), false
}

func (h *Handler) RegularNewLine(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	errFunc("ERROR unknown symbol", start, "\\n")
	switchCond(golexgen.ConditionINIT)

	return golexgen.Token{}, true
}

func (h *Handler) RegularEscapeNewLine(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])

	return golexgen.Token{}, true
}

func (h *Handler) RegularEscapeTab(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])
	return golexgen.Token{}, true
}

func (h *Handler) RegularEscapeQota(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])
	return golexgen.Token{}, true
}

func (h *Handler) RegularSymb(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])
	return golexgen.Token{}, true
}

func (h *Handler) StartLiteral(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	switchCond(golexgen.ConditionLITERAL)
	h.start = start

	return golexgen.Token{}, true
}

func (h *Handler) Literal1(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])
	return golexgen.Token{}, true
}

func (h *Handler) LiteralEnd(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	switchCond(golexgen.ConditionINIT)

	return golexgen.NewToken(
		TagString,
		h.start, end,
		h.str,
	), false
}

func (h *Handler) LiteralNewLine(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])
	return golexgen.Token{}, true
}

func (h *Handler) LiterlaChar(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	h.str += string(text[start.Index():end.Index()])
	return golexgen.Token{}, true
}

func (h *Handler) Num(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	return golexgen.NewToken(
		golexgen.TagNum,
		start, end,
		string(text[start.Index():end.Index()]),
	), false
}

func (h *Handler) Any(
	text []rune,
	start, end golexgen.Position,
	errFunc golexgen.ErrFunc,
	switchCond golexgen.SwitchConditionFunc,
) (golexgen.Token, golexgen.Continued) {
	errFunc("ERROR unknown symbol", start, string(text[start.Index():end.Index()]))
	return golexgen.Token{}, true
}

func main() {
	filePath := "./examples/lab1.6/test.txt"

	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	scn := golexgen.NewScanner([]rune(string(content)), &Handler{})

	t := scn.NextToken()
	for t.Tag() != golexgen.EOP {
		fmt.Println(t.String())
		t = scn.NextToken()
	}
}
