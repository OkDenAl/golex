package main

type Regexp struct {
	RegexpVal           *FiniteAutomata
	ActionName          string
	SwitchConditionName string
	Continued           bool
	Edit                bool
}

type Condition struct {
	Name        string
	Regexps     []Regexp
	RegexpsLen  int
	UnionRegexp Regexp
}

type GeneratorInfo struct {
	Conditions   map[string]Condition
	AllRegexps   []Regexp
	UnionRegexps []Regexp
}

func NewCondition(name string, regxps []Regexp, regxpsLen int, unionReg *FiniteAutomata) Condition {
	return Condition{
		Name:    name,
		Regexps: regxps,
		UnionRegexp: Regexp{
			RegexpVal:  unionReg,
			ActionName: name,
		},
		RegexpsLen: regxpsLen,
	}
}

type AutomataWithNaming struct {
	automata *FiniteAutomata
	naming   map[int]string
}

const InitialCond = "INIT"

func (r *Program) ProcessOneAutomata() *GeneratorInfo {
	gi := GeneratorInfo{Conditions: make(map[string]Condition), AllRegexps: make([]Regexp, 0)}
	conds := make(map[string][]Regexp)
	condsUnionAutomata := make(map[string]AutomataWithNaming)
	for _, rule := range r.rules.ruleArr {
		startCondName := InitialCond
		if rule.startCondition != nil {
			startCondName = rule.startCondition.condition.val
		}

		automata := rule.expr.Compile().CompileV2()
		automata.setLexemName(rule.name.val)

		var switchCond string
		if rule.switchCondition != nil {
			switchCond = rule.switchCondition.nextCondition.val
		}

		reg := Regexp{
			RegexpVal:           automata,
			ActionName:          rule.name.val,
			SwitchConditionName: switchCond,
			Continued:           rule.contin != nil,
			Edit:                rule.edit != nil,
		}

		conds[startCondName] = append(conds[startCondName], reg)
		gi.AllRegexps = append(gi.AllRegexps, reg)

		if _, ok := condsUnionAutomata[startCondName]; !ok {
			naming = make(map[int]string)
			for _, a := range automata.TerminalStates {
				naming[a.State+1] = a.LexemName
			}
			condsUnionAutomata[startCondName] = AutomataWithNaming{
				automata: automata,
				naming:   naming,
			}
		} else {
			naming = condsUnionAutomata[startCondName].naming
			condsUnionAutomata[startCondName] = AutomataWithNaming{
				automata: condsUnionAutomata[startCondName].automata.UnionNext(automata.copy()),
				naming:   naming,
			}
		}
	}

	for key, val := range conds {
		curAutomataInfo := condsUnionAutomata[key]
		flPos = curAutomataInfo.automata.flPos
		letters = curAutomataInfo.automata.letters
		naming = curAutomataInfo.naming
		curAutomata := curAutomataInfo.automata.CompileV2()

		gi.Conditions[key] = NewCondition(key, val, len(val), curAutomata)
		gi.UnionRegexps = append(gi.UnionRegexps, Regexp{
			RegexpVal:  curAutomata,
			ActionName: key,
		})
	}

	return &gi
}
