package main

import "os"

type Regexp struct {
	RegexpVal           *FiniteState
	ActionName          string
	SwitchConditionName string
	Continued           bool
	Edit                bool
}

type Condition struct {
	Name        string
	Regexps     []Regexp
	RegexpsLen  int
	UnionRegexp []Regexp
}

func NewCondition(name string, regxps []Regexp, regxpsLen int) Condition {
	return Condition{
		Name:       name,
		Regexps:    regxps,
		RegexpsLen: regxpsLen,
	}
}

type GeneratorInfo struct {
	Conditions   map[string]Condition
	AllRegexps   []Regexp
	UnionRegexps []Regexp
}

const InitialCond = "INIT"

func (r *Program) ProcessOneAutomata() *GeneratorInfo {
	gi := GeneratorInfo{Conditions: make(map[string]Condition), AllRegexps: make([]Regexp, 0)}
	conds := make(map[string][]Regexp)
	res := NewAutomata()
	for _, rule := range r.rules.ruleArr {
		automata := rule.expr.Compile()
		automata.setLexemName(rule.name.val)

		res.Union(automata.copy())
		res.ToGraph(os.Stdout)
		startCondName := InitialCond
		if rule.startCondition != nil {
			startCondName = rule.startCondition.condition.val
		}
		if _, ok := conds[startCondName]; !ok {
			conds[startCondName] = make([]Regexp, 0)
		}

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
	}
	//res.ToGraph(os.Stdout)
	for key, val := range conds {
		gi.UnionRegexps = append(gi.UnionRegexps, Regexp{
			RegexpVal:  res,
			ActionName: "UnionRegexps",
		})
		gi.Conditions[key] = NewCondition(key, val, len(val))
	}

	return &gi
}

func (r *Program) Process() *GeneratorInfo {
	gi := GeneratorInfo{Conditions: make(map[string]Condition), AllRegexps: make([]Regexp, 0)}
	conds := make(map[string][]Regexp)
	for _, rule := range r.rules.ruleArr {
		automata := rule.expr.Compile()
		startCondName := InitialCond
		if rule.startCondition != nil {
			startCondName = rule.startCondition.condition.val
		}
		if _, ok := conds[startCondName]; !ok {
			conds[startCondName] = make([]Regexp, 0)
		}

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
	}

	for key, val := range conds {
		gi.Conditions[key] = NewCondition(key, val, len(val))
	}

	return &gi
}
