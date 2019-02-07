package regex

// An approach to abstract classes?
type quantifierToken struct {
	*baseToken
	min int
	max int
	t Token

	maxStrategy func(qt *quantifierToken, m *matcher) (bool, *RegexException)
}

func newQuantifierToken(t Token, min, max int, maxStrategy func(*quantifierToken, *matcher) (bool, *RegexException)) *quantifierToken {
	return &quantifierToken{
		baseToken:   newBaseToken(),
		min:         min,
		max:         max,
		t:           t,
		maxStrategy: maxStrategy,
	}
}

func newQuantifierGreedyToken(q *quantifier, t Token) *quantifierToken {
	return newQuantifierToken(t, q.min, q.max, greedyStrategy)
}

func newQuantifierNonGreedyToken(q *quantifier, t Token) *quantifierToken {
	return newQuantifierToken(t, q.min, q.max, nongreedyStrategy)
}

func newQuantifierPossessiveToken(q *quantifier, t Token) *quantifierToken {
	return newQuantifierToken(t, q.min, q.max, possessiveStrategy)
}

func greedyStrategy(tk *quantifierToken, m *matcher) (bool, *RegexException) {
	startPos := m.getTextPos()
	savedState := m.saveThenPushNextStack(tk.cloneDecrement())
	ret, err := tk.t.match(m)
	if err != nil {
		return false, err
	}
	if !ret {
		m.restoreNextStack(savedState)
		m.setTextPos(startPos)
		return tk.getNext().match(m)
	}

	return true, nil
}

func nongreedyStrategy(tk *quantifierToken, m *matcher) (bool, *RegexException) {
	startPos := m.getTextPos()
	savedState := m.saveNextStack()

	ret, err := tk.getNext().match(m)
	if err != nil {
		return false, err
	}
	if !ret {
		m.setTextPos(startPos)
		m.restoreNextStack(savedState)
		m.pushNextStack(tk.cloneDecrement())

		return tk.t.match(m)
	}

	return true, nil
}

func possessiveStrategy(tk *quantifierToken, m *matcher) (bool, *RegexException) {
	for i := 0; i < tk.max || tk.max == -1; i++ {
		savedState := m.saveAndResetNextStack()
		startPos := m.getTextPos()
		ret, err := tk.t.match(m)
		if err != nil {
			return false, err
		}
		if !ret {
			m.restoreNextStack(savedState)
            m.setTextPos(startPos)
            break
        }
    }

    return tk.getNext().match(m)
}


func (tk *quantifierToken) match(m *matcher) (bool, *RegexException) {
	if tk.min != 0 {
		m.pushNextStack(tk.cloneDecrement())
		return tk.t.match(m)
	}

	if tk.max != 0 {
		return tk.maxStrategy(tk, m)
	}

	return tk.getNext().match(m)
}

func (tk *quantifierToken) cloneDecrement() *quantifierToken {
	newQt := newQuantifierToken(tk.t, tk.decrement(tk.min), tk.decrement(tk.max), tk.maxStrategy)
	newQt.next = tk.next

	return newQt
}

func (tk *quantifierToken) decrement(v int) int {
	ret := v
	if v > 0 {
		ret = v - 1
	}

	return ret
}
