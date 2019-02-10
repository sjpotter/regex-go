package regex

type quantifierPossessiveToken struct {
	*baseToken
	min     int
	max     int
	t       Token
	paired  *quantifierPossessiveToken
	matched bool
}

func newQuantifierPossessiveToken(q *quantifier, t Token) *quantifierPossessiveToken {
	return &quantifierPossessiveToken{baseToken: newBaseToken(), min: q.min, max: q.max, t: t}
}

func (tk *quantifierPossessiveToken) match(m *matcher) bool {
	if state, ok := m.tokenState[tk].(*nextState); ok {
		tk.deleteUntil(tk, state.myNext, m)

		if tk.min != 0 {
			delete(m.tokenState, tk)
			return false
		}

		if tk.matched {
			tmp := tk.paired
			// if I matched, don't backtrack to previous matches
			for tmp != nil {
				tmp.matched = false
				tmp = tmp.paired
			}

			tk.matched = false
			return true
		}

		delete(m.tokenState, tk)
		return false
	}

	tk.matched = false

	if tk.paired != nil {
		tk.paired.matched = true
	}

	if tk.min != 0 || tk.max != 0 {
		m.tokenState[tk] = &nextState{myNext: tk.getNext(), startPos: m.getTextPos()}

		nextQt := tk.cloneDecrement()
		nextQt.paired = tk

		tk.insertAfter(tk, nextQt)
		tk.insertAfter(tk, tk.t)
		return true
	}

	return false
}

func (tk *quantifierPossessiveToken) cloneDecrement() *quantifierPossessiveToken {
	q := &quantifier{min: tk.min, max: tk.max}

	newQt := newQuantifierPossessiveToken(q, tk.t)
	newQt.next = tk.next

	return newQt
}
