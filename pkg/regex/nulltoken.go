package regex

var (
	nullToken = &null{}
)

type null struct {}

func (n *null) match(m *matcher) (bool, *RegexException) { return m.matchNextStack() }

func (n *null) reverse() (Token, *RegexException) { return n, nil }

func (n *null) getNext() Token { return n }

func (n *null) setNext(Token) {}
func (n *null) reverseToken(cur Token) (Token, *RegexException) { return n, nil }

func (n *null) captureGroup() int { return -1 }
func (n *null) quantifiable() bool { return false }
func (n *null) testable() bool { return false }
func (n *null) normalExpression() bool { return false }