package regex

type quantifierNonGreedyToken struct {
	*baseToken
	min int
	max int
	t   Token
}

func newQuantifierNonGreedyToken(q *quantifier, t Token) *quantifierNonGreedyToken {
	return &quantifierNonGreedyToken{baseToken: newBaseToken(), min: q.min, max: q.max, t: t}
}

func (tk *quantifierNonGreedyToken) match(m *matcher) bool {
	if state, ok := m.tokenState[tk].(*nextState); ok {
		delete(m.tokenState, tk)

		// If we are backtracking before reaching the minimum amount of repititions, no match this path
		if tk.min != 0 {
			tk.deleteUntil(tk, state.myNext, m)
			return false
		}

		// Non Greedy backtracking when reached the minimum amount of reptitions has 2 paths
		// 1. have we tried adding one iterations of our quantification, if so add it
		// 2. if already added it, that means no path this direction
		if tk.max != 0 {
			if tk.next == state.myNext { //i.e. we haven't inserted anything. backtracking as non greedy matching
				m.tokenState[tk] = &nextState{myNext: state.myNext, startPos: m.getTextPos()}
				tk.insertAfter(tk, tk.cloneDecrement())
				tk.insertAfter(tk, tk.t)
				return true
			} else { // after insert, we didn't match, so backtracking as no match
				tk.deleteUntil(tk, state.myNext, m)
				m.setTextPos(state.startPos)
				return false
			}
		}
	}

	m.tokenState[tk] = &nextState{myNext: tk.getNext(), startPos: m.getTextPos()}

	if tk.min != 0 {
		tk.insertAfter(tk, tk.cloneDecrement())
		tk.insertAfter(tk, tk.t)
	}

	return true
}

func (tk *quantifierNonGreedyToken) cloneDecrement() *quantifierNonGreedyToken {
	q := &quantifier{min: tk.min, max: tk.max}

	newQt := newQuantifierNonGreedyToken(q, tk.t)
	newQt.next = tk.next

	return newQt
}
