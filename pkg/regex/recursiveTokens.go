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

func (tk *recursiveToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil { //this isn't a valid path, so this isn't a valid capture
		delete(m.tokenState, tk)
		return false
	}

	t := m.getCaptureToken(tk.group).copy()

    m1 := m.copyMatcher()
    m1.t = t

    ret := m1.matchFrom(m.getTextPos())
    if ret {
    	m.copy(m1)
    	m.tokenState[tk] = 1
    	return true
	}

    return false
}

func (tk *recursiveToken) copy() Token {
	return &recursiveToken{baseToken: newBaseToken(), group: tk.group, recursive: tk.recursive}
}

func (tk *recursiveToken) reverse() Token {
	panic(newRegexException("Can't LookBehind with Recursive Tokens"))
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

func (tk *recursiveEndToken) match(m *matcher) bool {
	tk.mOld.copy(m)
	return tk.getNext().match(tk.mOld)
}
