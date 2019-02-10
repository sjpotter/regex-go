package regex

import (
	"fmt"
	"strings"
)

type characterToken struct {
	*baseToken
	cc *characterClass
}

func newCharacterToken(regex []rune, regexPos int) (*characterToken, int) {
	cc, regexPos := getCharacterClass(regex, regexPos)

	return &characterToken{baseToken: newBaseToken(), cc: cc}, regexPos
}

func (tk *characterToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		delete(m.tokenState, tk)
		return false
	}

	text := m.getText()
	textPos := m.getTextPos()

	dir := m.getDirection()
	if dir == -1 {
		if textPos == 0 {
			return false
		}
		textPos--
	}

	if textPos < len(text) && textPos >= 0 {
		if tk.cc.match(text[textPos]) {
			if dir == 1 {
				m.setTextPos(textPos + 1)
			} else {
				m.setTextPos(textPos)
			}
			m.tokenState[tk] = 1
			return true
		}
	}

	return false
}

func (tk *characterToken) quantifiable() bool {
	return true
}

func (tk *characterToken) copy() Token {
	return &characterToken{baseToken: newBaseToken(), cc: tk.cc}
}

func getCharacterClass(regex []rune, regexPos int) (*characterClass, int) {
	//NOTE: always make sure that the regex string is advanced if new cases are added
	switch regex[regexPos] {
	case '[':
		end := strings.Index(string(regex[regexPos:]), "]")
		if end == -1 {
			panic(newRegexException(fmt.Sprintf("need to end characters class (started at index: %v) with a brace", regexPos)))
		}
		//cut out the [ and ]
		c := newCharacterClass(regex, regexPos+1, regexPos+end-1)
		return c, regexPos + end + 1
	case '\\':
		c := newCharacterClass(regex, regexPos, regexPos+1)
		return c, regexPos + 2
	case '.':
		return allCharacters(), regexPos + 1
	case '+', '*', '?', '^', '$', '|', '(', ')':
		panic(newRegexException(fmt.Sprintf("invalid characters in regex: %v at index: %v", regex[regexPos:regexPos+1], regexPos)))
	default: //plain characters
		c := newCharacterClass(regex, regexPos, regexPos)
		return c, regexPos + 1
	}
}
