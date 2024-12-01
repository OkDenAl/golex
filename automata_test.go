package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"regexp/syntax"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type state struct {
	limit int
}

type Generator struct {
	re    *syntax.Regexp
	rand  *rand.Rand
	debug bool
}

func (g *Generator) generate(s *state, re *syntax.Regexp) string {
	//fmt.Println("re:", re, "sub:", re.Sub)
	op := re.Op
	switch op {
	case syntax.OpNoMatch:
	case syntax.OpEmptyMatch:
		return ""
	case syntax.OpLiteral:
		res := ""
		r := g.rand.Float64()
		if r > 0.2 {
			for _, r := range re.Rune {
				res += string(r)
				if g.rand.Float64() > 0.8 {
					c := printableChars[g.rand.Intn(len(printableChars))]
					res += string(c)
				} else if g.rand.Float64() > 0.7 {
					res += string(r)
				}
			}
		} else {
			for _, r := range re.Rune {
				res = string(r) + res
			}
		}
		return res
	case syntax.OpCharClass:
		// number of possible chars
		sum := 0
		for i := 0; i < len(re.Rune); i += 2 {
			if g.debug {
				fmt.Printf("Range: %#U-%#U\n", re.Rune[i], re.Rune[i+1])
			}
			sum += int(re.Rune[i+1]-re.Rune[i]) + 1
			if re.Rune[i+1] == runeRangeEnd {
				sum = -1
				break
			}
		}
		// pick random char in range (inverse match group)
		if sum == -1 {
			possibleChars := []uint8{}
			for j := 0; j < len(printableChars); j++ {
				c := printableChars[j]
				//fmt.Printf("Char %c %d\n", c, c)
				// Check c in range
				for i := 0; i < len(re.Rune); i += 2 {
					if rune(c) >= re.Rune[i] && rune(c) <= re.Rune[i+1] {
						possibleChars = append(possibleChars, c)
						break
					}
				}
			}
			//fmt.Println("Possible chars: ", possibleChars)
			if len(possibleChars) > 0 {
				c := possibleChars[g.rand.Intn(len(possibleChars))]
				if g.debug {
					fmt.Printf("Generated rune %c for inverse range %v\n", c, re)
				}
				return string([]byte{c})
			}
		}
		if g.debug {
			fmt.Println("Char range: ", sum)
		}
		r := g.rand.Intn(int(sum))
		var ru rune
		sum = 0
		for i := 0; i < len(re.Rune); i += 2 {
			gap := int(re.Rune[i+1]-re.Rune[i]) + 1
			if sum+gap > r {
				ru = re.Rune[i] + rune(r-sum)
				break
			}
			sum += gap
		}
		if g.debug {
			fmt.Printf("Generated rune %c for range %v\n", ru, re)
		}
		return string(ru)
	case syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		chars := printableChars
		if op == syntax.OpAnyCharNotNL {
			chars = printableCharsNoNL
		}
		c := chars[g.rand.Intn(len(chars))]
		return string([]byte{c})
	case syntax.OpBeginLine:
	case syntax.OpEndLine:
	case syntax.OpBeginText:
	case syntax.OpEndText:
	case syntax.OpWordBoundary:
	case syntax.OpNoWordBoundary:
	case syntax.OpCapture:
		if g.debug {
			fmt.Println("OpCapture", re.Sub, len(re.Sub))
		}
		return g.generate(s, re.Sub0[0])
	case syntax.OpStar:
		// Repeat zero or more times
		res := ""
		count := g.rand.Intn(s.limit + 1)
		for i := 0; i < count; i++ {
			for _, r := range re.Sub {
				res += g.generate(s, r)
			}
		}
		return res
	case syntax.OpPlus:
		// Repeat one or more times
		res := ""
		count := g.rand.Intn(s.limit) + 1
		for i := 0; i < count; i++ {
			for _, r := range re.Sub {
				res += g.generate(s, r)
			}
		}
		return res
	case syntax.OpQuest:
		// Zero or one instances
		res := ""
		count := g.rand.Intn(2)
		if g.debug {
			fmt.Println("Quest", count)
		}
		for i := 0; i < count; i++ {
			for _, r := range re.Sub {
				res += g.generate(s, r)
			}
		}
		return res
	case syntax.OpRepeat:
		// Repeat one or more times
		if g.debug {
			fmt.Println("OpRepeat", re.Min, re.Max)
		}
		res := ""
		count := 0
		re.Max = int(math.Min(float64(re.Max), float64(s.limit)))
		if re.Max > re.Min {
			count = g.rand.Intn(re.Max - re.Min + 1)
		}
		if g.debug {
			fmt.Println(re.Max, count)
		}
		for i := 0; i < re.Min || i < (re.Min+count); i++ {
			for _, r := range re.Sub {
				res += g.generate(s, r)
			}
		}
		return res
	case syntax.OpConcat:
		// Concatenate sub-regexes
		res := ""
		for _, r := range re.Sub {
			res += g.generate(s, r)
		}
		return res
	case syntax.OpAlternate:
		if g.debug {
			fmt.Println("OpAlternative", re.Sub, len(re.Sub))
		}
		i := g.rand.Intn(len(re.Sub))
		return g.generate(s, re.Sub[i])
	default:
		fmt.Fprintln(os.Stderr, "[reg-gen] Unhandled op: ", op)
	}
	return ""
}

// limit is the maximum number of times star, range or plus should repeat
// i.e. [0-9]+ will generate at most 10 characters if this is set to 10
func (g *Generator) Generate(limit int) string {
	return g.generate(&state{limit: limit}, g.re)
}

// create a new generator
func NewGenerator(regex string) (*Generator, error) {
	re, err := syntax.Parse(regex, syntax.Perl)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Compiled re ", re)
	return &Generator{
		re:   re,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (gen *Generator) SetSeed(seed int64) {
	gen.rand = rand.New(rand.NewSource(seed))
}

func Generate(regex string, limit int) (string, error) {
	g, err := NewGenerator(regex)
	if err != nil {
		return "", err
	}
	return g.Generate(limit), nil
}

func setup(t *testing.T, regexp string) *FiniteAutomata {
	scn := NewScanner(bufio.NewReader(strings.NewReader(regexp)), NewCompiler())
	scn.regularMode = true
	tokens := scn.GetTokens()
	tokens = append(tokens, Token{
		val: string([]rune{endSymbol}),
		tag: TagCharacter,
	})
	parser := New(tokens)
	expr, ok := parser.regExpr()
	require.True(t, ok)

	res := expr.Compile()
	r := res.CompileV2()

	r.ToGraph(os.Stdout)
	res.ToGraph(os.Stdout)
	return r
}

func TestFiniteState_Execute(t *testing.T) {
	type args struct {
		reg string
	}

	tests := []struct {
		name   string
		args   args
		count  int
		maxLen int
		accept bool
	}{
		{
			name:   "(((a|b)|(abc)*)|p)[0-9]*",
			args:   args{reg: "(((a|b)|(abc)*)|p)[0-9]*"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "[(((a|b)|(abc)*)|p)[0-9]*]",
			args:   args{reg: "[(((a|b)|(abc)*)|p)[0-9]*]"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "[\\n\\t ]*",
			args:   args{reg: "[\\n\\t ]*"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(([A-Za-z0-9]+\\{[0-9]+\\})|[0-9]+)",
			args:   args{reg: "(([A-Za-z0-9]+\\{[0-9]+\\})|[0-9]+)"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(\\'[A-Za-z ]*\\')",
			args:   args{reg: "(\\'[A-Za-z ]*\\')"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(\\()(a|b*)\\)",
			args:   args{reg: "(\\()(a|b*)\\)"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "[\\^1 ]+",
			args:   args{reg: "[\\^1 ]+"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(ab)*acge",
			args:   args{reg: "(ab)*ac"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "[A-Da-d]+a\\{",
			args:   args{reg: "[A-Da-d]ab\\{"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "\"[A-Za-z]?\"",
			args:   args{reg: "\"[A-Za-z]?\""},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "[0-9]*|([0-9]*h)",
			args:   args{reg: "([0-9]*|([0-9]*h))"},
			count:  100,
			maxLen: 100,
		},
		//{
		//	name:   "[^0-9]1(2)*",
		//	args:   args{reg: "a*[^0-9]1a[kek]lol"},
		//	count:  100,
		//	maxLen: 100,
		//},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := 0
			sut := setup(t, tt.args.reg)
			r := regexp.MustCompile("^" + tt.args.reg + "$")

			for i := 1; i < tt.maxLen; i++ {
				for j := 0; j < tt.count; j++ {
					str, err := Generate(tt.args.reg, i)
					require.Nil(t, err)
					if sut.MatchString(str) != r.MatchString(str) {
						f += 1
						fmt.Println(sut.MatchString(str), r.MatchString(str))
						fmt.Println(tt.args.reg)
						fmt.Println(str)
						fmt.Println(r.MatchString(str))
					}
				}
			}
			require.Equal(t, 0, f)
		})
	}
}
