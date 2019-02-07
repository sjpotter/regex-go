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

func (tk *backReferenceToken) match(m *matcher) (bool, *RegexException) {
	if m.getDirection() != -1 {
		return tk.forwardMatch(m)
    }

    return tk.reverseMatch(m)
}

func (tk *backReferenceToken) forwardMatch(m *matcher) (bool, *RegexException) {
	text := m.getText()
	textPos := m.getTextPos()
	stored, err := m.getGroup(tk.backReference)
	if err != nil {
		return false, err
	}

	if stored == nil { // TODO: unsure this is correct, maybe its an automatic and continue to next match
		return false, nil
	}

	for i :=0; i < len(stored); i++ {
		if textPos >= len(text) || stored[i] != text[textPos] {
			return false, nil
		}
        textPos++
	}

    m.setTextPos(textPos)

	return tk.getNext().match(m)
}

func (tk *backReferenceToken) reverseMatch(m *matcher) (bool, *RegexException) {
	text := m.getText()
	textPos := m.getTextPos() - 1
	stored, err := m.getGroup(tk.backReference)
	if err != nil {
		return false, err
	}

	if stored == nil { // TODO: unsure this is correct, maybe its an automatic and continue to next match
		return false, nil
	}

	for i := len(stored) -1; i >= 0; i-- {
		if textPos < 0 || stored[i] != text[textPos] {
			return false, nil
		}
        textPos--
	}

    m.setTextPos(textPos)

	return tk.getNext().match(m)
}

func (tk *backReferenceToken) quantifiable() bool {
	return true
}