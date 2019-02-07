package regex

type lookBehindExpressionToken struct {
	*baseToken
	t        *normalExpresionToken
	positive bool
	lbet bool
}

func newLookBehindExpressionToken(net *normalExpresionToken, positive bool) *lookBehindExpressionToken {
	return &lookBehindExpressionToken{
		baseToken: newBaseToken(),
		t:         net,
		positive:  positive,
	}
}

func (tk *lookBehindExpressionToken) match(m *matcher) (bool, *RegexException) {
	textPos := m.getTextPos()
	m.setDirection(-1)

	// Empty stack as only matters that its string of tokens match
	savedState := m.saveAndResetNextStack()

	ret, err := tk.t.match(m)
	if err != nil {
		return false, err
	}

	m.restoreNextStack(savedState)

	m.setDirection(1)
	m.setTextPos(textPos)

	if tk.positive {
		if !ret {
			return false, nil
		}
	} else {
		if ret {
			return false, nil
		}
	}

	return tk.getNext().match(m)
}

func (tk *lookBehindExpressionToken) testable() bool {
	return true
}