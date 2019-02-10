package regex

import (
	"strings"
	"unicode"
)

type quantifier struct {
	min int
	max int
}

func quantifierParse(token Token, regex []rune, regexPos int) (Token, int) {
	var q *quantifier
	var qt Token

	if len(regex) == regexPos {
		return nil, regexPos
	}

	switch regex[regexPos] {
	case '{':
		endPos := strings.IndexRune(string(regex[regexPos:]), '}')
		if endPos != -1 {
			q = handleVariable(regex, regexPos+1)
			if q != nil {
				regexPos += endPos + 1 // endPos is relative to the start of regexPos because of slice method
			}
		}
	case '*':
		q = &quantifier{0, -1}
		regexPos++
	case '+':
		q = &quantifier{1, -1}
		regexPos++
	case '?':
		q = &quantifier{0, 1}
		regexPos++

	}

	if q != nil {
		qt = newQuantifierGreedyToken(q, token)

		if regexPos < len(regex) {
			switch regex[regexPos] {
			case '?':
				qt = newQuantifierNonGreedyToken(q, token)
				regexPos++
			case '+':
				qt = newQuantifierPossessiveToken(q, token)
				regexPos++
			}
		}
	}

	return qt, regexPos
}

func handleVariable(regex []rune, regexPos int) *quantifier {
	if len(regex) > regexPos && unicode.IsDigit(regex[regexPos]) { //is there a number after the {
		val := int(regex[regexPos] - '0')
		regexPos++
		for len(regex) > regexPos && unicode.IsDigit(regex[regexPos]) {
			val *= 10
			val += int(regex[regexPos] - '0')
			regexPos++
		}

		if len(regex) > regexPos && regex[regexPos] == '}' { //if it's {#} we need to match exact
			return &quantifier{val, val}
		} else if len(regex) > regexPos && regex[regexPos] == ',' { // can be {#,} or {#,#}
			min := val
			regexPos++
			if len(regex) > regexPos && regex[regexPos] == '}' { // {#,}
				return &quantifier{min, -1}
			} else { // determine if it's a valid // {#,#}
				if len(regex) > regexPos && unicode.IsDigit(regex[regexPos]) {
					val = int(regex[regexPos] - '0')
					regexPos++
					for len(regex) > regexPos && unicode.IsDigit(regex[regexPos]) {
						val *= 10
						val += int(regex[regexPos] - '0')
						regexPos++
					}
				}

				if len(regex) > regexPos && regex[regexPos] == '}' { //maybe valid {#,#}
					if min <= val { // {min,val} only valid if min <= val (max)
						return &quantifier{min, val}
					}
				}
			}
		}
	}

	//regex following { invalid as a quantifier
	return nil
}

func quantifierDecrement(v int) int {
	ret := v
	if v > 0 {
		ret = v - 1
	}

	return ret
}
