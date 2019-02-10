package regex

type ifThenElseToken struct {
	*baseToken
	ifToken   Token
	thenToken Token
	elseToken Token
}

func newIfThenElseToken(ifToken, thenToken, elsetoken Token) Token {
	if !ifToken.testable() {
		panic(newRegexException("IfThenElseToken: ifToken not a TestableToken"))
	}

	return &ifThenElseToken{
		baseToken: newBaseToken(),
		ifToken:   ifToken,
		thenToken: thenToken,
		elseToken: elsetoken,
	}
}

func (tk *ifThenElseToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		if state, ok := m.tokenState[tk].(*nextState); !ok {
			panic(newRegexException("ifThenElseToken state is not an *nextState"))
		} else {
			tk.deleteUntil(tk, state.myNext, m)
		}
	}

	m1 := m.copyMatcher()
	m1.t = tk.ifToken

	ret := m1.matchFrom(m.getTextPos())

	exec := tk.thenToken
	if !ret {
		exec = tk.elseToken
	}

	tk.insertAfter(tk, exec)

	state := &nextState{
		myNext: tk.getNext(),
	}
	m.tokenState[tk] = state

	return true
}

func (tk *ifThenElseToken) reverse() Token {
	panic(newRegexException("Can't LookBehind with ifThenElse Token"))

}

func (tk *ifThenElseToken) quantifiable() bool {
	return true
}

func (tk *ifThenElseToken) copy() Token {
	return newIfThenElseToken(tk.ifToken, tk.thenToken, tk.elseToken)
}
