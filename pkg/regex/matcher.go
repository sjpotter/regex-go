package regex

type matcher struct {
	groups     []*StringStack
	text       []rune
	textPos    int
	direction  int
	parenCount int
	captureMap map[int]*normalExpresionToken
	t          Token

	//compound tokens have a list of their own tokens to match against and what follows the compound token.
	//what follows is what's pushed onto the nextStack
	nextStack *TokenStack
}

func NewMatcher(t *tokenizer) (*matcher, *RegexException) {
	token, err := t.Tokenize()
	if err != nil {
		return nil, err
	}

	groups := []*StringStack{}
	for i := 0; i < t.captureCount; i++ {
		groups = append(groups, &StringStack{})
	}

	return &matcher{
		t:          token,
		direction:  1,
		parenCount: t.captureCount,
		captureMap: t.captureMap,
		nextStack:  &TokenStack{},
		groups:     groups,
	}, nil
}

func (m *matcher) copyMatcher() *matcher {
	m1 := &matcher{
		text:       m.text,
		textPos:    m.textPos,
		direction:  1,
		parenCount: m.parenCount,
		captureMap: m.captureMap,
		nextStack:  m.nextStack,
	}

	for i := 0; i < m1.parenCount; i++ {
		newStack := &StringStack{}
		newStack.Copy(m.groups[i])
		m1.groups = append(m1.groups, newStack)
	}

	return m1
}

func (m *matcher) copy(m1 *matcher) {
	m.textPos = m1.textPos
	m.nextStack = m1.nextStack
}

func (m *matcher) Match(text string) (bool, *RegexException) {
	runeText := []rune(text)

	for i := 0; i < len(runeText) || i == 0; i++ {
		for j := 0; j < m.parenCount; j++ {
			m.groups[j] = &StringStack{}
		}
		m.text = runeText
		m.nextStack = &TokenStack{}
		m.textPos = i
		ret, err := m.t.match(m)
		if err != nil {
			return false, err
		}
		if ret {
			m.groups[0].Push(string(m.text[i:m.textPos]))
			return true, nil
		}
	}

	return false, nil
}

func (m *matcher) GetGroups() []string {
	var ret []string

	for _, v := range m.groups {
		if v.Len() > 0 {
			ret = append(ret, v.Peek())
		} else {
			ret = append(ret, "")
		}
	}

	return ret
}

func (m *matcher) getGroup(pos int) ([]rune, *RegexException) {
	return []rune(m.groups[pos].Peek()), nil //FIXME
}

func (m *matcher) GetGroup(pos int) (string) {
	return  m.groups[pos].Peek()
}

func (m *matcher) pushGroup(pos int, t string) {
	m.groups[pos].Push(t)
}

func (m *matcher) popGroup(pos int) {
	m.groups[pos].Pop()
}

func (m *matcher) getTextPos() int {
	return m.textPos
}

func (m *matcher) setTextPos(pos int) {
	m.textPos = pos
}

func (m *matcher) getText() []rune {
	return m.text
}

func (m *matcher) getDirection() int {
	return m.direction
}

func (m *matcher) setDirection(dir int) {
	m.direction = dir
}

func (m *matcher) saveNextStack() *TokenStack {
	savedState := m.nextStack
	m.nextStack = &TokenStack{}
	m.nextStack.Copy(savedState)

	return savedState
}

func (m *matcher) saveAndResetNextStack() *TokenStack {
	savedState := m.saveNextStack()
	m.nextStack = &TokenStack{}

	return savedState
}

func (m *matcher) pushNextStack(t Token) {
	m.nextStack.Push(t)
}

func (m *matcher) saveThenPushNextStack(t Token) *TokenStack {
	saved := m.saveNextStack()
	m.pushNextStack(t)

	return saved
}

func (m *matcher) restoreNextStack(savedStack *TokenStack) {
	m.nextStack = savedStack
}

func (m *matcher) matchNextStack() (bool, *RegexException) {
	if m.nextStack.Len() == 0 {
		return true, nil
	}

	return m.nextStack.Pop().match(m)
}

func (m *matcher) getCaptureToken(pos int) (*normalExpresionToken, *RegexException) {
	if pos >= len(m.captureMap) {
		return nil, newRegexException("Trying to retrieve a token for a capture group that doesn't exist");
	}

	return m.captureMap[pos], nil
}
