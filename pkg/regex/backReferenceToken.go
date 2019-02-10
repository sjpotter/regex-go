package regex

type backReferenceToken struct {
	*baseToken
	backReference int
}

func newBackReferenceToken(backReference int) *backReferenceToken {
	return &backReferenceToken{
		baseToken:     newBaseToken(),
		backReference: backReference,
	}
}

func (tk *backReferenceToken) match(m *matcher) bool {
	var ret bool

	if m.tokenState[tk] != nil {
		delete(m.tokenState, tk)
		return false
	}

	if m.getDirection() != -1 {
		ret = tk.forwardMatch(m)
	} else {
		return tk.reverseMatch(m)
	}

	if ret {
		m.tokenState[tk] = 1
	}

	return ret
}

func (tk *backReferenceToken) forwardMatch(m *matcher) bool {
	text := m.getText()
	textPos := m.getTextPos()
	stored := m.getGroup(tk.backReference)

	if stored == nil { // TODO: unsure this is correct, maybe its an automatic and continue to next match
		return false
	}

	for i := 0; i < len(stored); i++ {
		if textPos >= len(text) || stored[i] != text[textPos] {
			return false
		}
		textPos++
	}

	m.setTextPos(textPos)

	return true
}

func (tk *backReferenceToken) reverseMatch(m *matcher) bool {
	text := m.getText()
	textPos := m.getTextPos() - 1
	stored := m.getGroup(tk.backReference)

	if stored == nil { // TODO: unsure this is correct, maybe its an automatic and continue to next match
		return false
	}

	for i := len(stored) - 1; i >= 0; i-- {
		if textPos < 0 || stored[i] != text[textPos] {
			return false
		}
		textPos--
	}

	m.setTextPos(textPos)

	return true
}

func (tk *backReferenceToken) quantifiable() bool {
	return true
}

func (tk *backReferenceToken) copy() Token {
	return newBackReferenceToken(tk.backReference)
}
