package main

type Regexp struct {
	RegexpVal  *FiniteState
	ActionName string
}

type GeneratorInfo struct {
	RegexpsLen int
	Regexps    []Regexp
}

func (r *Rules) Process() *GeneratorInfo {
	gi := GeneratorInfo{Regexps: make([]Regexp, 0)}
	for _, rule := range r.rules {
		automata := rule.expr.Compile()
		gi.Regexps = append(gi.Regexps, Regexp{
			RegexpVal:  automata,
			ActionName: rule.name.val,
		})
	}

	gi.RegexpsLen = len(gi.Regexps)

	return &gi
}
