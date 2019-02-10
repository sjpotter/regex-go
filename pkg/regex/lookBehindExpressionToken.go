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

func (tk *lookBehindExpressionToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		delete(m.tokenState, tk)
		return false
	}

	m1 := m.copyMatcher()
	m1.t = tk.t
	m1.setDirection(-1)


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

func (tk *lookBehindExpressionToken) testable() bool {
	return true
}

func (tk *lookBehindExpressionToken) copy() Token {
	return newLookBehindExpressionToken(tk.t, tk.positive)
}