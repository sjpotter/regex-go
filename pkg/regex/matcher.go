package regex

import (
	"runtime/debug"
)

type matcher struct {
	groups     []*StringStack
	text       []rune
	textPos    int
	direction  int
	parenCount int
	captureMap map[int]*normalExpresionToken
	tokenState map[Token]interface{}
	t          Token

	//compound tokens have a list of their own tokens to match against and what follows the compound token.
	//what follows is what's pushed onto the nextStack
	nextStack *TokenStack
}

func NewMatcher(t *tokenizer) (m *matcher, re *RegexException) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}

		m = nil
		ok := false
		if re, ok = e.(*RegexException); !ok {
			re = newRegexException(string(debug.Stack()))
		}
	}()

	token := t.Tokenize()

	var groups []*StringStack
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
		tokenState: make(map[Token]interface{}),
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
		tokenState: make(map[Token]interface{}),
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

func (m *matcher) Match(text string) (ret bool, re *RegexException) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}

		ret = false
		ok := false
		if re, ok = e.(*RegexException); !ok {
			re = newRegexException(string(debug.Stack()))
		}
	}()

	runeText := []rune(text)
	m.text = runeText

	for i := 0; i < len(runeText) || i == 0; i++ {
		m.tokenState = make(map[Token]interface{})
		for j := 0; j < m.parenCount; j++ {
			m.groups[j] = &StringStack{}
		}
		if m.matchFrom(i) {
			return true, nil
		}

	}

	ret = false
	re = nil

	return
}

func (m *matcher) matchFrom(pos int) bool {
	m.textPos = pos
	curTok := m.t

	for {
		ret := curTok.match(m)
		if ret {
			curTok = curTok.getNext()
			if curTok == nil {
				str := string(m.text[pos:m.textPos])
				m.groups[0].Push(&str)
				return true
			}
		} else {
			curTok = curTok.getPrev()
			if curTok == nil {
				return false
			}
		}
	}
}

func (m *matcher) GetGroups() []*string {
	var ret []*string

	for _, v := range m.groups {
		if v.Len() > 0 {
			ret = append(ret, v.Peek())
		} else {
			ret = append(ret, nil)
		}
	}

	return ret
}

func (m *matcher) getGroup(pos int) []rune {
	s := m.groups[pos].Peek()
	if s != nil {
		return []rune(*s)
	}

	return nil
}

func (m *matcher) GetGroup(pos int) *string {
	return  m.groups[pos].Peek()
}

func (m *matcher) pushGroup(pos int, t *string) {
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

func (m *matcher) matchNextStack() bool {
	if m.nextStack.Len() == 0 {
		return true
	}

	return m.nextStack.Pop().match(m)
}

func (m *matcher) getCaptureToken(pos int) *normalExpresionToken {
	if pos >= len(m.captureMap) {
		panic(newRegexException("Trying to retrieve a token for a capture group that doesn't exist"))
	}

	return m.captureMap[pos]
}
