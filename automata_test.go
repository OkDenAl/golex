package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, regexp string) *FiniteState {
	scn := NewScanner(bufio.NewReader(strings.NewReader(regexp)), NewCompiler())
	scn.regularMode = true
	tokens := GetTokens(scn)
	parser := New(tokens)

	expr, ok := parser.regExpr()
	require.True(t, ok)

	return expr.Compile()
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
			count:  10,
			maxLen: 10,
		},
		{
			name:   "[(((a|b)|(abc)*)|p)[0-9]*]",
			args:   args{reg: "[(((a|b)|(abc)*)|p)[0-9]*]"},
			count:  10,
			maxLen: 10,
		},
		{
			name:   "[\\n\\t ]*",
			args:   args{reg: "[\\n\\t ]*"},
			count:  10,
			maxLen: 10,
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
			count:  10,
			maxLen: 10,
		},
		{
			name:   "(\\()(a|b*)\\)",
			args:   args{reg: "(\\()(a|b*)\\)"},
			count:  10,
			maxLen: 10,
		},
		{
			name:   "[\\^1 ]*",
			args:   args{reg: "[\\^1 ]*"},
			count:  10,
			maxLen: 10,
		},
		{
			name:   "(ab)*ac",
			args:   args{reg: "(ab)*ac"},
			count:  10,
			maxLen: 10,
		},
		{
			name:   "[A-Da-d]+a\\{",
			args:   args{reg: "[A-Da-d]+a\\{"},
			count:  100,
			maxLen: 100,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := 0
			sut := setup(t, tt.args.reg)
			r := regexp.MustCompile(tt.args.reg)

			for i := 1; i < tt.maxLen; i++ {
				for j := 0; j < tt.count; j++ {
					str, err := Generate(tt.args.reg, i)
					require.Nil(t, err)
					if sut.Execute(str) != r.MatchString(str) {
						f += 1
						fmt.Println(sut.Execute(str), r.MatchString(str))
						fmt.Println(tt.args.reg, str)
					}
				}
			}
			require.Equal(t, 0, f)
		})
	}
}
