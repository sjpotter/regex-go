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

func (tk *atomicExpressionToken) match(m *matcher) (bool, *RegexException) {
	it := tk.altIterator()
	start := m.getTextPos()

	for it.hasNext() {
		savedStack := m.saveAndResetNextStack()

		ret, err := it.next().match(m)
		if err != nil {
			return false, err
		}
		if ret {
			m.restoreNextStack(savedStack)
			return tk.getNext().match(m)
		}

		m.restoreNextStack(savedStack)
		m.setTextPos(start)
	}

	return false, nil
}