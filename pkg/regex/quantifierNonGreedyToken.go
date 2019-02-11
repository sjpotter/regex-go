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

/* 4 cases
 1) advancing
  a. if minimum is not zero, have to insert sub token and be like any other quantifier
  b. if mimimum is zero, just return true
 2) backtracking
  a. first time backtracking, if min == 0 (i.e. all minimum quantifiers matched) and max != 0 (i.e. can add), add
  b. otherwise fail
*/
func (tk *quantifierNonGreedyToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		if state, ok := m.tokenState[tk].(*nextState); ok {
			firstTime := tk.getNext() == state.myNext

			if firstTime && tk.min == 0 && tk.max != 0 {
				tk.insertAfter(tk, tk.cloneDecrement())
				tk.insertAfter(tk, tk.t)
				return true
			}

			tk.deleteUntil(tk, state.myNext, m)
			delete(m.tokenState, tk)
			return false
		} else {
			panic(newRegexException("quantifierNonGreedyToken state is not a *nextState"))
		}
	}
	// case 1a and 1b
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

func (tk *quantifierNonGreedyToken) copy() Token {
	q := &quantifier{min: tk.min, max: tk.max}

	qt := newQuantifierNonGreedyToken(q, tk.t)

	return qt
}