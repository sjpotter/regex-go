package regex

type lookAheadExpressionToken struct {
	*baseToken
	t *normalExpresionToken
	positive bool
	laet bool
}

func newLookAheadExpressionToken(net *normalExpresionToken, positive bool) *lookAheadExpressionToken {
	return &lookAheadExpressionToken{
		baseToken: newBaseToken(),
		t:         net,
		positive:  positive,
	}
}

func (tk *lookAheadExpressionToken) match(m *matcher) (bool, *RegexException) {
	textPos := m.getTextPos()

	// Empty stack as only matters that its string of tokens match
    savedState := m.saveAndResetNextStack()

    ret, err := tk.t.match(m);
    if err != nil {
    	return false, err
	}

    m.restoreNextStack(savedState);
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

func (tk *lookAheadExpressionToken) testable() bool {
	return true
}