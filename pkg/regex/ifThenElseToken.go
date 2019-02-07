package regex

type ifThenElseToken struct {
	*baseToken
	ifToken   Token
	thenToken Token
	elseToken Token
}

func newIfThenElseToken(ifToken, thenToken, elsetoken Token) (Token, *RegexException) {
	if !ifToken.testable() {
		return nil, newRegexException("IfThenElseToken: ifToken not a TestableToken")
	}

	return &ifThenElseToken{
		baseToken: newBaseToken(),
		ifToken:   ifToken,
		thenToken: thenToken,
		elseToken: elsetoken,
	}, nil
}

func (tk *ifThenElseToken) match(m *matcher) (bool, *RegexException) {
	// Empty stack for if clause, as only the tokens within it define true/false for the then/else clauses
    savedStack := m.saveAndResetNextStack()
    ret, err := tk.ifToken.match(m)
    if err != nil {
    	return false, err
	}

    // stack is returned for then/else clause as they continue matching next tokens.
    m.restoreNextStack(savedStack)

    exec := tk.thenToken
    if !ret {
	    exec = tk.elseToken
    }

    return exec.match(m);
}

func (tk *ifThenElseToken) reverse() (Token, *RegexException) {
	ifReversed, err := tk.ifToken.reverse()
	if err != nil {
		return nil, err
	}
	thenReversed, err := tk.thenToken.reverse()
	if err != nil {
		return nil, err
	}
	elseReversed, err := tk.elseToken.reverse()
	if err != nil {
		return nil, err
	}

	cur, err := newIfThenElseToken(ifReversed, thenReversed, elseReversed)
	if err != nil {
		return nil, err
	}

	return tk.baseToken.reverseToken(cur)
}

func (tk *ifThenElseToken) quantifiable() bool {
	return true
}