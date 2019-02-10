package regex

type startCaptureToken struct {
	*baseToken
	t Token
	capture int
}

func newStartCaptureToken(capture int, t Token) *startCaptureToken {
	return &startCaptureToken{
		baseToken: newBaseToken(),
		t:         t,
		capture:   capture,
	}
}

func (tk *startCaptureToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		if state, ok := m.tokenState[tk].(*nextState); !ok {
			panic(newRegexException("startCaptureToken state is not an *startCaptureState"))
		} else {
			tk.deleteUntil(tk, state.myNext, m)
		}
	} else {
		state := &nextState{
			myNext:   tk.getNext(),
		}
		m.tokenState[tk] = state
	}

	startPos := m.getTextPos()
	end := newEndCaptureToken(tk.capture, startPos)
	tk.insertAfter(tk, end)
	tk.insertAfter(tk, tk.t)

	return true
}

func (tk *startCaptureToken) quantifiable() bool {
	return true
}

func (tk *startCaptureToken) copy() Token {
	return newStartCaptureToken(tk.capture, tk.t)
}


type endCaptureToken struct {
	*baseToken
	capture int
	startPos int
}

func newEndCaptureToken(capture, startPos int) *endCaptureToken {
	return &endCaptureToken{
		baseToken: newBaseToken(),
		capture:   capture,
		startPos:  startPos,
	}
}

func (tk *endCaptureToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil { //this isn't a valid path, so this isn't a valid capture
		delete(m.tokenState, tk)
		m.popGroup(tk.capture)
		return false
	}

	m.tokenState[tk] = 1

	subRuneSlice := m.getText()[tk.startPos:m.getTextPos()]

	s := string(subRuneSlice)
	m.pushGroup(tk.capture, &s)

	return true
}
