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

func (e *expressionToken) altIterator(dir int) *altIterator {
	pos := 0
	if dir == -1 {
		pos = len(e.alts)-1
	}
	return &altIterator{
		alt: e.alts,
		pos: pos,
		dir: dir,
	}
}

type altIterator struct {
	alt []Token
	pos int
	dir int
}

func (ai *altIterator) hasNext() bool {
	if ai.dir == 1 {
		return len(ai.alt) > ai.pos
	} else {
		return ai.pos >= 0
	}
}

func (ai *altIterator) next() Token {
	ret := ai.alt[ai.pos]
	ai.pos += ai.dir

	return ret
}
