package regex

type atomicExpressionToken struct {
	*expressionToken
	aet bool
}

func newAtomicExpressionToken() *atomicExpressionToken {
	return &atomicExpressionToken{
		expressionToken: newExpressionToken(),
	}
}

func (tk *atomicExpressionToken) match(m *matcher) bool {
	var state *expressionState

	if m.tokenState[tk] != nil {
		var ok bool
		if state, ok = m.tokenState[tk].(*expressionState); !ok {
			panic(newRegexException("atomicExpressionToken state is not an *expressionState"))
		}
	} else {
		state = &expressionState{
			it:       tk.altIterator(),
			startPos: m.getTextPos(),
			myNext:   tk.getNext(),
		}
		m.tokenState[tk] = state
	}

	if !state.matched && state.it.hasNext() {
		m.setTextPos(state.startPos)
		tk.deleteUntil(tk, state.myNext, m)
		tk.insertAfter(tk, newAtomicEndToken(state))
		tk.insertAfter(tk, state.it.next())

		return true
	}

	delete(m.tokenState, tk)
	return false
}

func (tk *atomicExpressionToken) copy() Token {
	aet := &atomicExpressionToken{
		expressionToken: newExpressionToken(),
	}

	aet.alts = tk.alts

	return aet
}

type atomicEndToken struct {
	*baseToken
	state *expressionState
}

func newAtomicEndToken(state *expressionState) *atomicEndToken {
	return &atomicEndToken{baseToken: newBaseToken(), state: state}
}

func (tk *atomicEndToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		delete(m.tokenState, tk)
		return false
	}

	tk.state.matched = true

	m.tokenState[tk] = 1
	return true
}

func (tk *atomicEndToken) copy() Token {
	return &atomicEndToken{baseToken: newBaseToken(), state: tk.state}
}
