package regex

type lookAheadExpressionToken struct {
	*baseToken
	t        *normalExpresionToken
	positive bool
	laet     bool
}

func newLookAheadExpressionToken(net *normalExpresionToken, positive bool) *lookAheadExpressionToken {
	return &lookAheadExpressionToken{
		baseToken: newBaseToken(),
		t:         net,
		positive:  positive,
	}
}

func (tk *lookAheadExpressionToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		delete(m.tokenState, tk)
		return false
	}

	m1 := m.copyMatcher()
	m1.t = tk.t

	ret := m1.matchFrom(m.getTextPos())

	if tk.positive {
		if !ret {
			return false
		}
	} else {
		if ret {
			return false
		}
	}

	m.tokenState[tk] = 1
	return true
}

func (tk *lookAheadExpressionToken) testable() bool {
	return true
}

func (tk *lookAheadExpressionToken) copy() Token {
	return newLookAheadExpressionToken(tk.t, tk.positive)
}
