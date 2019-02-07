package regex

type RegexException struct {
	s string
}

func newRegexException(s string) *RegexException {
	return &RegexException{s: s}
}

func (re *RegexException) Error() string {
	return re.s
}
