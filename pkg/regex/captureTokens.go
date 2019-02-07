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

func (tk *startCaptureToken) match(m *matcher) (bool, *RegexException) {
	startPos := m.getTextPos()

	end := newEndCaptureToken(tk.capture, startPos)
	end.setNext(tk.getNext())
	m.pushNextStack(end)

	return tk.t.match(m)
}

func (tk *startCaptureToken) quantifiable() bool {
	return true
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

func (tk *endCaptureToken) match(m *matcher) (bool, *RegexException) {
	subRuneSlice := m.getText()[tk.startPos:m.getTextPos()]

	m.pushGroup(tk.capture,string(subRuneSlice))

	ret, err := tk.getNext().match(m)
	if err != nil {
		return false, err
	}
	if ret == true {
		return true, nil
	}

	//this isn't a valid path, so this isn't a valid capture
	m.popGroup(tk.capture)

	return false, nil
}
