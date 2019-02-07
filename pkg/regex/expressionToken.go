package regex

type expressionToken struct {
	*baseToken

	alts []Token
}

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

func (e *expressionToken) internalReverse() *RegexException {
	newAlts := []Token{}

	for i := 0; i < len(e.alts); i-- {
		reversed, err := e.alts[i].reverse()
		if err != nil {
			return err
		}

		newAlts = append(newAlts, reversed)
	}

	e.alts = newAlts

	return nil
}

func (e *expressionToken) reverse() (Token, *RegexException) {
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