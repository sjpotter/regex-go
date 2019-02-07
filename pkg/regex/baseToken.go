package regex

type Token interface {
	match(m *matcher) (bool, *RegexException)
	getNext() Token
	setNext(Token)
	reverse() (Token, *RegexException)
	reverseToken(cur Token) (Token, *RegexException)
	captureGroup() int
	quantifiable() bool
	testable() bool
	normalExpression() bool
}

type baseToken struct {
	next Token
}

// super() like function
func newBaseToken() *baseToken {
	return &baseToken{
		next: nullToken,
	}
}

func (tk *baseToken) match(m *matcher) (bool, *RegexException) {
	panic("Unimplemented: always needs to be overriden")
}

func (tk *baseToken) getNext() Token {
	if tk.next == nil {
		return nullToken
	}

	return tk.next
}

func (tk *baseToken) setNext(n Token) {
	tk.next = n
}

func (tk *baseToken) reverse() (Token, *RegexException) {
	return tk.reverseToken(tk)
}

func (tk *baseToken) reverseToken(cur Token) (Token, *RegexException) {
	if cur.getNext().(*null) != nil {
		return cur, nil
	}

	prev, err := cur.getNext().reverse()
	if err != nil {
		return nil, err
	}

	cur.setNext(nullToken)

	tmp := prev;
	for tmp.getNext() != nullToken {
		tmp = tmp.getNext()
	}
	tmp.setNext(tk)

	return prev, nil
}

func (tk *baseToken) captureGroup() int {
	return -1
}

func (tk *baseToken) quantifiable() bool {
	return false
}

func (tk *baseToken) testable() bool {
	return false
}

func (tk *baseToken) normalExpression() bool {
	return false
}