package regex

type quantifierPossessiveToken struct {
	*baseToken
	min     int
	max     int
	t       Token
	matched *bool
}

func newQuantifierPossessiveToken(q *quantifier, t Token) *quantifierPossessiveToken {
	var sharedMatch bool
	return &quantifierPossessiveToken{baseToken: newBaseToken(), min: q.min, max: q.max, t: t, matched: &sharedMatch}
}

// same exact logic as greedy, except we only backtrack to the furthest quantifier match, so all decremented
// quantifiers have to share a matched state to reset once we try to continue the regex past the quantifier
func (tk *quantifierPossessiveToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		if state, ok := m.tokenState[tk].(*nextState); ok {
			firstTime := tk.getNext() != state.myNext

			tk.deleteUntil(tk, state.myNext, m)

			if firstTime && tk.min == 0 && tk.max != 0 && *tk.matched {
				*tk.matched = false
				m.setTextPos(state.startPos)
				return true
			}

			*tk.matched = false
			delete(m.tokenState, tk)
			return false
		} else {
			panic(newRegexException("quantifierPossessiveToken state is not a *nextState"))
		}
	}

	*tk.matched = true

	m.tokenState[tk] = &nextState{myNext: tk.getNext(), startPos: m.getTextPos()}

	if tk.min != 0 || tk.max != 0 {
		tk.insertAfter(tk, tk.cloneDecrement())
		tk.insertAfter(tk, tk.t)
		return true
	}

	return true
}

func (tk *quantifierPossessiveToken) cloneDecrement() *quantifierPossessiveToken {
	q := &quantifier{min: tk.min, max: tk.max}

	newQt := newQuantifierPossessiveToken(q, tk.t)
	newQt.matched = tk.matched

	return newQt
}

func (tk *quantifierPossessiveToken) copy() Token {
	q := &quantifier{min: tk.min, max: tk.max}

	qt := newQuantifierPossessiveToken(q, tk.t)
	qt.matched = tk.matched

	return qt
}