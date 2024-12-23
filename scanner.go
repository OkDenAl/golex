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
	compiler *Compiler
	curPos   Position
	comments []comment

	regularMode bool
	prevToken   Token
	wasEscape   bool
}

func NewScanner(programFile *bufio.Reader, compiler *Compiler) Scanner {
	return Scanner{compiler: compiler, curPos: NewPosition(programFile)}
}

func (scn *Scanner) printComments() {
	for _, comm := range scn.comments {
		fmt.Println(comm)
	}
}

func (scn *Scanner) addComment(start, end Position, val string) {
	scn.comments = append(scn.comments, newComment(start, end, val))
}

func (scn *Scanner) nextToken() Token {
	if scn.regularMode && (scn.wasEscape || scn.curPos.Cp() != '/') {
		return scn.nextTokenRegular()
	}
	for scn.curPos.Cp() != -1 {
		for scn.curPos.IsWhiteSpace() {
			scn.curPos.Next()
		}
		start := scn.curPos
		curWord := ""

		switch scn.curPos.Cp() {
		case '\r':
			fallthrough
		case '\n':
			scn.curPos.Next()
			scn.prevToken = NewToken(TagNL, start, start, "NEW_LINE")

			return NewToken(TagNL, start, start, "NEW_LINE")
		case '<':
			scn.curPos.Next()
			scn.prevToken = NewToken(TagOpenStartCondition, start, start, "<")

			return NewToken(TagOpenStartCondition, start, start, "<")
		case '>':
			scn.curPos.Next()
			scn.prevToken = NewToken(TagCloseStartCondition, start, start, ">")

			return NewToken(TagCloseStartCondition, start, start, ">")
		case '(':
			scn.curPos.Next()
			scn.prevToken = NewToken(TagDefaultOpenBracket, start, start, "(")

			return NewToken(TagDefaultOpenBracket, start, start, "(")
		case ')':
			scn.curPos.Next()
			scn.prevToken = NewToken(TagDefaultCloseBracket, start, start, ")")

			return NewToken(TagDefaultCloseBracket, start, start, ")")
		case '#':
			scn.curPos.Next()
			var pos Position
			for scn.curPos.Cp() != -1 && scn.curPos.Cp() != '\n' {
				curWord += string(scn.curPos.GetSymbol())
				pos = scn.curPos
				scn.curPos.Next()
			}
			scn.addComment(start, pos, curWord)
		case '%':
			curWord += string(rune(scn.curPos.Cp()))
			scn.curPos.Next()

			if scn.curPos.Cp() == -1 || (scn.curPos.GetSymbol() != 'x' && scn.curPos.GetSymbol() != '%') {
				scn.compiler.addMessage(true, start, "invalid syntax")
				scn.curPos.SkipErrors()
				scn.prevToken = NewToken(TagErr, scn.curPos, scn.curPos, "")

				return NewToken(TagErr, scn.curPos, scn.curPos, "")
			}

			curWord += string(rune(scn.curPos.Cp()))
			pos := scn.curPos
			scn.curPos.Next()

			var token Token
			switch curWord {
			case "%%":
				scn.prevToken = NewToken(TagRulesMarker, start, pos, curWord)
				token = NewToken(TagRulesMarker, start, pos, curWord)
			case "%x":
				scn.prevToken = NewToken(TagStateMarker, start, pos, curWord)
				token = NewToken(TagStateMarker, start, pos, curWord)
			}

			return token
		case '/':
			if scn.regularMode {
				scn.regularMode = false
				scn.curPos.Next()
				return NewToken(TagRegularMarker, start, start, "/")
			}
			scn.regularMode = true
			scn.curPos.Next()

			return NewToken(TagRegularMarker, start, start, "/")
		default:
			var pos Position
			if scn.curPos.IsLowerLetter() {
				for scn.curPos.IsLowerLetter() {
					curWord += string(scn.curPos.GetSymbol())
					pos = scn.curPos
					scn.curPos.Next()
				}
			} else if scn.curPos.IsUpperLetter() {
				curWord += string(scn.curPos.GetSymbol())
				pos = scn.curPos
				scn.curPos.Next()
				for scn.curPos.IsLetter() || scn.curPos.Cp() == '_' || scn.curPos.IsDigit() {
					curWord += string(scn.curPos.GetSymbol())
					pos = scn.curPos
					scn.curPos.Next()
				}
			}

			if curWord == "" {
				scn.compiler.addMessage(true, start, "invalid syntax")
				scn.curPos.SkipErrors()

				scn.prevToken = NewToken(TagErr, scn.curPos, scn.curPos, "")

				return NewToken(TagErr, scn.curPos, scn.curPos, "")
			}

			const (
				begin        = "begin"
				continueName = "continue"
				edit         = "edit"
			)

			switch curWord {
			case begin:
				scn.prevToken = NewToken(TagBegin, start, pos, begin)
				return NewToken(TagBegin, start, pos, begin)
			case continueName:
				scn.prevToken = NewToken(TagContinue, start, pos, continueName)
				return NewToken(TagContinue, start, pos, continueName)
			case edit:
				scn.prevToken = NewToken(TagEdit, start, pos, edit)
				return NewToken(TagEdit, start, pos, edit)
			}

			scn.prevToken = NewToken(TagName, scn.curPos, scn.curPos, curWord)
			return NewToken(TagName, start, pos, curWord)
		}
	}

	return NewToken(TagEOP, scn.curPos, scn.curPos, "")
}

func (scn *Scanner) nextTokenRegular() Token {
	if scn.curPos.Cp() != -1 {
		start := scn.curPos
		var token Token
		switch scn.curPos.Cp() {
		case '(':
			token = NewToken(TagOpenParen, start, start, string(rune(scn.curPos.Cp())))
		case ')':
			token = NewToken(TagCloseParen, start, start, string(rune(scn.curPos.Cp())))
		case '[':
			token = NewToken(TagOpenBracket, start, start, string(rune(scn.curPos.Cp())))
		case ']':
			token = NewToken(TagCloseBracket, start, start, string(rune(scn.curPos.Cp())))
		case '{':
			token = NewToken(TagOpenBrace, start, start, string(rune(scn.curPos.Cp())))
		case '}':
			token = NewToken(TagCloseBrace, start, start, string(rune(scn.curPos.Cp())))
		case '*':
			token = NewToken(TagStar, start, start, string(rune(scn.curPos.Cp())))
		case '+':
			token = NewToken(TagPlus, start, start, string(rune(scn.curPos.Cp())))
		case '?':
			token = NewToken(TagQuestion, start, start, string(rune(scn.curPos.Cp())))
		case '^':
			token = NewToken(TagCaret, start, start, string(rune(scn.curPos.Cp())))
		case '\\':
			scn.wasEscape = !scn.wasEscape // если последний символ был обратный слэш, то считаем, не эскейп
			token = NewToken(TagEscape, start, start, string(rune(scn.curPos.Cp())))
		case '|':
			token = NewToken(TagPipe, start, start, string(rune(scn.curPos.Cp())))
		case '-':
			token = NewToken(TagDash, start, start, string(rune(scn.curPos.Cp())))
		case '.':
			token = NewToken(TagAnyCharacter, start, start, string(rune(scn.curPos.Cp())))
		case ',':
			token = NewToken(TagComma, start, start, string(rune(scn.curPos.Cp())))
		case '\n':
			scn.regularMode = false
			scn.compiler.addMessage(true, start, "invalid syntax: expected /")
			token = NewToken(TagErr, start, start, string(rune(scn.curPos.Cp())))
		default:
			token = NewToken(TagCharacter, start, start, string(rune(scn.curPos.Cp())))
		}
		scn.curPos.Next()
		if scn.prevToken.Tag() == TagName && token.Tag() != TagErr {
			scn.compiler.addToken(scn.prevToken.val, token)
		}

		if token.Tag() != TagEscape {
			scn.wasEscape = false
		}

		return token
	}

	return NewToken(TagEOP, scn.curPos, scn.curPos, "")
}

func (scn *Scanner) scanNameTokens(curToken Token) []Token {
	var name string
	var tokens []Token
	tokens = append(tokens, curToken)
	t := scn.nextToken()
	for t.Tag() != TagEOP && t.Tag() != TagCloseBrace {
		if t.Tag() != TagErr {
			name += t.val
			tokens = append(tokens, t)
		}
		t = scn.nextToken()
	}
	res, ok := scn.compiler.namesTokens[name]
	if !ok {
		if t.Tag() == TagCloseBrace {
			tokens = append(tokens, t)
		}
		return tokens
	}

	return res
}

func (scn *Scanner) GetTokens() []Token {
	t := scn.nextToken()
	i := 1
	var tokens []Token
	for t.Tag() != TagEOP {
		if t.Tag() == TagOpenBrace {
			tokens = append(tokens, scn.scanNameTokens(t)...)
		} else if t.Tag() != TagErr {
			if t.Tag() == TagRegularMarker {
				if i%2 == 0 {
					tokens = append(tokens, Token{
						tag: TagCloseParen,
						val: string([]rune{')'}),
					})
					tokens = append(tokens, Token{
						tag: TagCharacter,
						val: string([]rune{endSymbol}),
					})
				}
			}

			tokens = append(tokens, t)

			if t.Tag() == TagRegularMarker {
				if i%2 == 1 {
					tokens = append(tokens, Token{
						tag: TagOpenParen,
						val: string([]rune{'('}),
					})
				}
				i++
			}
		}

		t = scn.nextToken()
	}

	if len(scn.compiler.messages) != 0 {
		scn.compiler.OutputMessages()
		panic("error while tokenizing")
	}

	return tokens
}

func (scn *Scanner) NextToken() Token {
	return scn.nextToken()
}
