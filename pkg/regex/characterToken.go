package regex

import (
	"fmt"
	"strings"
)

type characterToken struct {
	*baseToken
	cc *characterClass
}

func newCharacterToken(regex []rune, regexPos int) (*characterToken, int, *RegexException) {
	cc, regexPos, err := getCharacterClass(regex, regexPos)

	return &characterToken{baseToken: newBaseToken(), cc: cc}, regexPos, err
}

func (tk *characterToken) match(m *matcher) (bool, *RegexException) {
	text := m.getText()
	textPos := m.getTextPos()

	dir := m.getDirection()
	if dir == -1 {
		if textPos == 0 {
			return false, nil
		}
		textPos--
	}


	if textPos < len(text) && textPos >= 0 {
		if tk.cc.match(text[textPos]) {
			if dir == 1 {
				m.setTextPos(textPos+1)
			} else {
				m.setTextPos(textPos)
			}
			return tk.getNext().match(m)
		}
	}

	return false, nil
}

func (tk *characterToken) quantifiable() bool {
	return true
}

func getCharacterClass(regex []rune, regexPos int) (*characterClass, int, *RegexException) {
	//NOTE: always make sure that the regex string is advanced if new cases are added
	switch regex[regexPos] {
	case '[':
		end := strings.Index(string(regex[regexPos:]), "]")
		if end == -1 {
			return nil, -1, newRegexException(fmt.Sprintf("need to end characters class (started at index: %v) with a brace", regexPos))
		}
		//cut out the [ and ]
		c, err := newCharacterClass(regex, regexPos+1, regexPos+end-1)
		if err != nil {
			return nil, -1, err
		}
		return c, regexPos + end + 1, nil
	case '\\':
		c, err := newCharacterClass(regex, regexPos, regexPos+1)
		if err != nil {
			return nil, -1, err
		}
		return c, regexPos + 2, nil
	case '.':
		return allCharacters(), regexPos + 1, nil
	case '+', '*', '?', '^', '$', '|', '(', ')':
		return nil, -1, newRegexException(fmt.Sprintf("invalid characters in regex: %v at index: %v", regex[regexPos:regexPos+1], regexPos))
	default: //plain characters
		c, _ := newCharacterClass(regex, regexPos, regexPos)
		return c, regexPos + 1, nil
	}
}

