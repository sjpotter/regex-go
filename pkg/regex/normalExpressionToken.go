package regex

type normalExpresionToken struct {
	*expressionToken
}

func newNormalExpressionToken() *normalExpresionToken {
	return &normalExpresionToken{
		expressionToken: newExpressionToken(),
	}
}

func (ne *normalExpresionToken) match(m *matcher) (bool, *RegexException) {
	return ne.internalMatch(m, true)
}

func (ne *normalExpresionToken) matchNoFollow(m *matcher) (bool, *RegexException) {
	return ne.internalMatch(m, false)
}

func (ne *normalExpresionToken) internalMatch(m *matcher, goNext bool) (bool, *RegexException) {
	it := ne.altIterator()

	start := m.getTextPos()

	for it.hasNext() {
		savedStack := m.saveNextStack()

		if goNext {
			m.pushNextStack(ne.next)
		}

		t := it.next()
		ret, err := t.match(m)
		if err != nil {
			return false, nil
		}
		if ret {
			return ret, nil
		}

		m.restoreNextStack(savedStack)
		m.setTextPos(start)
	}

	return false, nil
}

func (ne *normalExpresionToken) normalExpression() bool {
	return true
}
