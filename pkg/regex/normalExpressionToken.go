package regex

type normalExpresionToken struct {
	*expressionToken
}

func newNormalExpressionToken() *normalExpresionToken {
	return &normalExpresionToken{
		expressionToken: newExpressionToken(),
	}
}

// Backtracking and Advancing are basially the same code path, try next iterator if it exists, so only difference
// is how state is setup (maintained state, or new state)
func (tk *normalExpresionToken) match(m *matcher) bool {
	var state *expressionState

	if m.tokenState[tk] != nil {
		var ok bool
		if state, ok = m.tokenState[tk].(*expressionState); !ok {
			panic(newRegexException("normalExpresionToken state is not an *expressionState"))
		}
	} else {
		state = &expressionState{
			it:       tk.altIterator(),
			startPos: m.getTextPos(),
			myNext:   tk.getNext(),
		}
		m.tokenState[tk] = state
	}

	if state.it.hasNext() {
		m.setTextPos(state.startPos)
		tk.deleteUntil(tk, state.myNext, m)
		tmp := state.it.next()
		tk.insertAfter(tk, tmp)
		return true
	}

	delete(m.tokenState, tk)
	return false
}

func (tk *normalExpresionToken) normalExpression() bool {
	return true
}

func (tk *normalExpresionToken) copy() Token {
	net := &normalExpresionToken{
		expressionToken: newExpressionToken(),
	}
	net.alts = tk.alts

	return net
}
