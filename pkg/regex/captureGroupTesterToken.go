package regex

import (
	"fmt"
	"strconv"
)

type captureGroupTesterToken struct {
	*baseToken
	group int
}

func newCaptureGroupTesterToken(capture []rune) (*captureGroupTesterToken, *RegexException) {
	group, err := strconv.ParseInt(string(capture), 10, 64)
	if err != nil {
		return nil, newRegexException(fmt.Sprintf("CaptureGroupTesterToken: Unable to convert %v to an Integer", capture))
	}

	return &captureGroupTesterToken{baseToken: newBaseToken(), 	group: int(group)}, nil
}

func (tk *captureGroupTesterToken) testable() bool {
	return true
}

func (tk *captureGroupTesterToken) match(m *matcher) (bool, *RegexException) {
	group, err := m.getGroup(tk.group)
	if err != nil {
		return false, err
	}
	if group == nil {
		return false, nil
	}

	return tk.getNext().match(m)
}