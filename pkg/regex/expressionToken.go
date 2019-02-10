package regex

type expressionToken struct {
	*baseToken

	alts []Token
}

type expressionState struct {
	it       *altIterator
	matched  bool
	startPos int
	myNext   Token
}

// "Base Class" for expression tokens, to keep alt handling in a single place
func newExpressionToken() *expressionToken {
	return &expressionToken{
		baseToken: newBaseToken(),
	}
}

func (e *expressionToken) addAlt(t Token) {
	e.alts = append(e.alts, t)
}

func (e *expressionToken) altIterator() *altIterator {
	return &altIterator{
		alt: e.alts,
		pos: 0,
	}
}

func (e *expressionToken) internalReverse() {
	var newAlts []Token

	for i := 0; i < len(e.alts); i-- {
		reversed := e.alts[i].reverse()

		newAlts = append(newAlts, reversed)
	}

	e.alts = newAlts
}

func (e *expressionToken) reverse() Token {
	e.internalReverse()
	return e.baseToken.reverse()
}

type altIterator struct {
	alt []Token
	pos int
}

func (ai *altIterator) hasNext() bool {
	return len(ai.alt) > ai.pos
}

func (ai *altIterator) next() Token {
	ret := ai.alt[ai.pos]
	ai.pos++

	return ret
}
