package regex

type quantifierGreedyToken struct {
	*baseToken
	min     int
	max     int
	t       Token
}

func newQuantifierGreedyToken(q *quantifier, t Token) *quantifierGreedyToken {
	return &quantifierGreedyToken{baseToken: newBaseToken(), min: q.min, max: q.max, t: t}
}

/* 4 cases
 1) advancing
  a. if min and max are both not 0, insert a quantifier like normal
  b. if both are 0, nothing more to quantify, just advance normally
 2) backtracking
  a.if min is not 0, failure case
  b. if min is 0, are we the end of a

// 1. going right - first time through
// 2. going left first time - matched quantifier
// 3. going left - didn't match quantifier
// 4. going left second time - matched quantifier
*/
func (tk *quantifierGreedyToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		if state, ok := m.tokenState[tk].(*nextState); ok {
			firstTime := tk.getNext() != state.myNext

			tk.deleteUntil(tk, state.myNext, m)

			// 1. only try to continue with regex the first time backtracking
			// 2. only try to continue the regex if tk.min == 0 as that means all required quantifiers matched
			// 3. optimization: only continue regex if tk.max != 0 as otherwise already went down this path
			if firstTime && tk.min == 0 && tk.max != 0 {
				m.setTextPos(state.startPos)
				return true
			}

			delete(m.tokenState, tk)
			return false
		} else {
			panic(newRegexException("quantifierGreedyToken state is not a *nextState"))
		}
	}

	m.tokenState[tk] = &nextState{myNext: tk.getNext(), startPos: m.getTextPos()}

	if tk.min != 0 || tk.max != 0 {
		tk.insertAfter(tk, tk.cloneDecrement())
		tk.insertAfter(tk, tk.t)
		return true
	}

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

	return qt
}
