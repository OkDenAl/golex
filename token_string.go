// Code generated by "stringer -type=DomainTag"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TagCode-0]
	_ = x[TagRulesMarker-1]
	_ = x[TagStateMarker-2]
	_ = x[TagName-3]
	_ = x[Tag-4]
	_ = x[TagRegexp-5]
	_ = x[TagNL-6]
	_ = x[TagCharacter-7]
	_ = x[TagOpenParen-8]
	_ = x[TagCloseParen-9]
	_ = x[TagOpenBracket-10]
	_ = x[TagCloseBracket-11]
	_ = x[TagOpenBrace-12]
	_ = x[TagCloseBrace-13]
	_ = x[TagStar-14]
	_ = x[TagPlus-15]
	_ = x[TagQuestion-16]
	_ = x[TagCaret-17]
	_ = x[TagEscape-18]
	_ = x[TagPipe-19]
	_ = x[TagDash-20]
	_ = x[TagComma-21]
	_ = x[TagAnyCharacter-22]
	_ = x[TagErr-23]
	_ = x[TagEOP-24]
}

const _DomainTag_name = "CodeRulesMarkerStateMarkerRuleNameTagRegexpNLCharacterOpenParenCloseParenOpenBracketCloseBracketOpenBraceCloseBraceStarPlusQuestionCaretEscapePipeDashCommaAnyCharacterErrEOP"

var _DomainTag_index = [...]uint8{0, 4, 15, 26, 34, 37, 43, 45, 54, 63, 73, 84, 96, 105, 115, 119, 123, 131, 136, 142, 146, 150, 155, 167, 170, 173}

func (i DomainTag) String() string {
	if i < 0 || i >= DomainTag(len(_DomainTag_index)-1) {
		return "DomainTag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DomainTag_name[_DomainTag_index[i]:_DomainTag_index[i+1]]
}
