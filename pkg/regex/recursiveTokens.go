package regex

type recursiveToken struct {
	*baseToken

	group int

	recursive bool
}

func newRecursiveToken(capture int) *recursiveToken {
	return &recursiveToken{
		baseToken:    newBaseToken(),
		group: capture,
	}
}

func (tk *recursiveToken) match(m *matcher) (bool, *RegexException) {
	t, err := m.getCaptureToken(tk.group)
	if err != nil {
		return false, err
	}

    // As recursive consumes state, need to have nextStack work correctly, this token resets the matcher state
    // for future matches after it executes
    m.pushNextStack(newRecursiveEndToken(m, tk.getNext()))

    m1 := m.copyMatcher()

    return t.matchNoFollow(m1)
}

func (tk *recursiveToken) reverse() (Token, *RegexException) {
	return nil, newRegexException("Can't LookBehind with Recursive Tokens")
}

func (tk *recursiveToken) quantifiable() bool {
	return true
}

type recursiveEndToken struct {
	*baseToken

	mOld *matcher
}

func newRecursiveEndToken(m *matcher, next Token) *recursiveEndToken {
	ret := &recursiveEndToken{
		baseToken: newBaseToken(),
		mOld:      m,
	}

	ret.setNext(next)

	return ret
}

func (tk *recursiveEndToken) match(m *matcher) (bool, *RegexException) {
	tk.mOld.copy(m);
	return tk.getNext().match(tk.mOld);
}
