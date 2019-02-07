package regex

import (
	"fmt"
	"unicode"
)

type anchorToken struct {
	*baseToken
	anchor rune
}

func newAnchorToken(r rune) *anchorToken {
	return &anchorToken{
		baseToken: newBaseToken(),
		anchor:    r,
	}
}

func (tk *anchorToken) match(m *matcher) (bool, *RegexException) {
	text := m.getText()
	textPos := m.getTextPos()

	switch tk.anchor {
	case '^':
		if textPos != 0 {
			return false, nil
		}
	case '$':
		if textPos != len(text) {
			return false, nil
		}
	case 'b', 'B': //word break anchor, negative word break anchor
		negative := tk.anchor == 'B'

		if textPos == 0 {
			if unicode.IsLetter(text[textPos]) {
				if negative { //don't want to match word breaks
					return false, nil
				}
			}
		} else if textPos == len(text) {
			if unicode.IsLetter(text[textPos - 1]) {
				if negative { //don't want to match words breaks
					return false, nil
				}
			}
		} else {
			if (unicode.IsSpace(text[textPos -1]) && unicode.IsLetter(text[textPos])) ||
				(unicode.IsSpace(text[textPos]) && unicode.IsLetter(text[textPos-1])) {
				if negative { // don't want to match words breaks
					return false, nil
				}
			}
		}

		// failed to match word breaks
		if !negative {
			return false, nil
		}
	default:
		return false, newRegexException(fmt.Sprintf("Unexpected ANCHOR token: %v", tk.anchor))
	}

	return tk.getNext().match(m)
}
