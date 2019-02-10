package regex

import (
	"fmt"
	"strconv"
)

type captureGroupTesterToken struct {
	*baseToken
	group int
}

func newCaptureGroupTesterToken(capture []rune) *captureGroupTesterToken {
	group, err := strconv.ParseInt(string(capture), 10, 64)
	if err != nil {
		panic(newRegexException(fmt.Sprintf("CaptureGroupTesterToken: Unable to convert %v to an Integer", capture)))
	}

	return &captureGroupTesterToken{baseToken: newBaseToken(), 	group: int(group)}
}

func (tk *captureGroupTesterToken) testable() bool {
	return true
}

func (tk *captureGroupTesterToken) match(m *matcher) bool {
	if m.tokenState[tk] != nil {
		delete(m.tokenState, tk)
		return false
	}

	group := m.getGroup(tk.group)
	if group == nil {
		return false
	}

	m.tokenState[tk] = 1
	return true
}

func (tk *captureGroupTesterToken) copy() Token {
	return &captureGroupTesterToken{baseToken: newBaseToken(), group: tk.group}
}