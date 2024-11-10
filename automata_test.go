package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, regexp string) *FiniteState {
	scn := NewScanner(bufio.NewReader(strings.NewReader(regexp)), NewCompiler())
	scn.regularMode = true
	tokens := scn.GetTokens()
	parser := New(tokens)

	expr, ok := parser.regExpr()
	require.True(t, ok)

	res := expr.Compile()
	res.ToGraph(os.Stdout)
	return res
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
			name:   "[\\^1 ]*",
			args:   args{reg: "[\\^1 ]*"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(ab)*ac",
			args:   args{reg: "(ab)*acge"},
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
			name:   "a{1}",
			args:   args{reg: "a{1}"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "a{0}",
			args:   args{reg: "a{0}bc"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(abc){3}",
			args:   args{reg: "(abc){3}"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(abc){2,}",
			args:   args{reg: "(abc){2,}"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "(abc){2,p}",
			args:   args{reg: "(abc){2,p}"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "a{2,4}",
			args:   args{reg: "a{2,4}"},
			count:  100,
			maxLen: 100,
		},
		{
			name:   "[A-Za-z]",
			args:   args{reg: "[A-Za-z]"},
			count:  100,
			maxLen: 100,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := 0
			sut := setup(t, tt.args.reg)

			for i := 1; i < tt.maxLen; i++ {
				for j := 0; j < tt.count; j++ {
					r := regexp.MustCompile("^" + tt.args.reg + "$")
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
