package main

import (
	"bufio"
	"fmt"
)

type comment struct {
	Fragment
	value string
}

func newComment(starting, following Position, value string) comment {
	return comment{NewFragment(starting, following), value}
}

func (c comment) String() string {
	return fmt.Sprintf("COMMENT %s-%s: %s", c.starting.String(), c.following.String(), c.value)
}

type Scanner struct {
	programReader *bufio.Reader
	compiler      *Compiler
	curPos        Position
	comments      []comment

	regularMode bool
	setMode     bool
}

func NewScanner(programFile *bufio.Reader, compiler *Compiler) Scanner {
	return Scanner{programReader: programFile, compiler: compiler, curPos: NewPosition(programFile)}
}

func (s *Scanner) printComments() {
	for _, comm := range s.comments {
		fmt.Println(comm)
	}
}

func (s *Scanner) NextToken() Token {
	if s.regularMode && s.curPos.Cp() != '/' {
		return s.nextTokenRegular()
	}
	for s.curPos.Cp() != -1 {
		for s.curPos.IsWhiteSpace() {
			s.curPos.Next()
		}
		start := s.curPos
		curWord := ""

		switch s.curPos.Cp() {
		case '\n':
			s.curPos.Next()
			return NewToken(TagNL, start, start, "NEW_LINE")
		case '%':
			curWord += string(rune(s.curPos.Cp()))
			s.curPos.Next()

			if s.curPos.Cp() == -1 || (s.curPos.GetSymbol() != 'x' && s.curPos.GetSymbol() != '%') {
				s.compiler.AddMessage(true, start, "invalid syntax")
				s.curPos.SkipErrors()

				return NewToken(TagErr, s.curPos, s.curPos, "")
			}

			curWord += string(rune(s.curPos.Cp()))
			pos := s.curPos
			s.curPos.Next()

			var token Token
			switch curWord {
			case "%%":
				token = NewToken(TagRulesMarker, start, pos, curWord)
			case "%x":
				token = NewToken(TagStateMarker, start, pos, curWord)
			}
			return token
		case '/':
			if s.regularMode {
				s.regularMode = false
				s.curPos.Next()
				return NewToken(TagRegularMarker, start, start, "/")
			}
			s.regularMode = true
			s.curPos.Next()

			return NewToken(TagRegularMarker, start, start, "/")
		default:
			var pos Position
			for s.curPos.Cp() != -1 && s.curPos.IsLetter() {
				curWord += string(s.curPos.GetSymbol())
				pos = s.curPos
				s.curPos.Next()
			}

			if curWord == "" {
				s.compiler.AddMessage(true, start, "invalid syntax")
				s.curPos.SkipErrors()

				return NewToken(TagErr, s.curPos, s.curPos, "")
			}

			return NewToken(TagRuleName, start, pos, curWord)
		}
	}

	return NewToken(TagEOP, s.curPos, s.curPos, "")
}

func (s *Scanner) nextTokenRegular() Token {
	defer func() {
		s.curPos.Next()
	}()
	if s.curPos.Cp() != -1 {
		start := s.curPos
		switch s.curPos.Cp() {
		case '(':
			return NewToken(TagOpenParen, start, start, string(rune(s.curPos.Cp())))
		case ')':
			return NewToken(TagCloseParen, start, start, string(rune(s.curPos.Cp())))
		case '[':
			s.setMode = true
			return NewToken(TagOpenBracket, start, start, string(rune(s.curPos.Cp())))
		case ']':
			s.setMode = false
			return NewToken(TagCloseBracket, start, start, string(rune(s.curPos.Cp())))
		case '{':
			return NewToken(TagOpenBrace, start, start, string(rune(s.curPos.Cp())))
		case '}':
			return NewToken(TagCloseBrace, start, start, string(rune(s.curPos.Cp())))
		case '*':
			return NewToken(TagStar, start, start, string(rune(s.curPos.Cp())))
		case '+':
			return NewToken(TagPlus, start, start, string(rune(s.curPos.Cp())))
		case '?':
			return NewToken(TagQuestion, start, start, string(rune(s.curPos.Cp())))
		case '^':
			return NewToken(TagCaret, start, start, string(rune(s.curPos.Cp())))
		case '\\':
			return NewToken(TagEscape, start, start, string(rune(s.curPos.Cp())))
		case '|':
			return NewToken(TagPipe, start, start, string(rune(s.curPos.Cp())))
		case '-':
			return NewToken(TagDash, start, start, string(rune(s.curPos.Cp())))
		case '.':
			return NewToken(TagAnyCharacter, start, start, string(rune(s.curPos.Cp())))
		case '\n':
			s.regularMode = false
			s.compiler.AddMessage(true, start, "invalid syntax: expected /")
			return NewToken(TagErr, start, start, string(rune(s.curPos.Cp())))
		default:
			return NewToken(TagCharacter, start, start, string(rune(s.curPos.Cp())))
		}
	}

	return NewToken(TagEOP, s.curPos, s.curPos, "")
}

func GetTokens(scn Scanner) []Token {
	t := scn.NextToken()
	var tokens []Token
	for t.Tag() != TagEOP {
		if t.Tag() != TagErr {
			tokens = append(tokens, t)
		}
		t = scn.NextToken()
	}

	if len(scn.compiler.messages) != 0 {
		scn.compiler.OutputMessages()
		panic("error while tokenizing")
	}

	return tokens
}
