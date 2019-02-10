package regex

type quantifierGreedyToken struct {
	*baseToken
	min int
	max int
	t Token
	paired *quantifierGreedyToken
	matched bool
}

func newQuantifierGreedyToken(q *quantifier, t Token) *quantifierGreedyToken {
	return &quantifierGreedyToken{baseToken: newBaseToken(), min: q.min, max: q.max, t: t}
}

// 4 cases
// 1. going right - first time through
// 2. going left first time - matched quantifier
// 3. going left - didn't match quantifier
// 4. going left second time - matched quantifier

func (tk *quantifierGreedyToken) match(m *matcher) bool {
	if state, ok := m.tokenState[tk].(*nextState); ok {
		tk.deleteUntil(tk, state.myNext, m)

		if tk.min != 0 {
			delete(m.tokenState, tk)
			return false
		}

		if tk.matched || (tk.paired != nil && tk.paired.matched) {
			if tk.paired != nil {
				tk.paired.matched = false
			}
			tk.matched = false
			m.setTextPos(state.startPos)
			return true
		}

		delete(m.tokenState, tk)
		return false
	}

	tk.matched = false
	if tk.paired != nil {
		tk.paired.matched = true
	} else { // if 0 max matches, can still continue from here?
		tk.matched = true
	}

	m.tokenState[tk] = &nextState{myNext: tk.getNext(), startPos: m.getTextPos()}

	if tk.min != 0 || tk.max != 0 {
		nextQt := tk.cloneDecrement()
		nextQt.paired = tk

		tk.insertAfter(tk, nextQt)
		tk.insertAfter(tk, tk.t)
		return true
	}

	tk.matched = true
	return true
}

func (tk *quantifierGreedyToken) cloneDecrement() *quantifierGreedyToken {
	q := &quantifier{min: quantifierDecrement(tk.min), max: quantifierDecrement(tk.max)}

	newQt := newQuantifierGreedyToken(q, tk.t)

	return newQt
}

func (tk *quantifierGreedyToken) copy() Token {
	q := &quantifier{min: tk.min, max: tk.max}

	qt := newQuantifierGreedyToken(q, tk.t)
	qt.paired = tk.paired

	return qt
}