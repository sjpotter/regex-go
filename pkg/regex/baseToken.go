package regex

type Token interface {
	match(m *matcher) bool
	getNext() Token
	setNext(Token)
	getPrev() Token
	setPrev(Token)
	copy() Token

	// "is a" interface type methods
	quantifiable() bool
	testable() bool
	normalExpression() bool
}

type baseToken struct {
	next Token
	prev Token
}

// super() like function
func newBaseToken() *baseToken {
	return &baseToken{}
}

func (tk *baseToken) match(m *matcher) bool {
	panic(newRegexException("Unimplemented: always needs to be overriden"))
}

func (tk *baseToken) getNext() Token {
	return tk.next
}

func (tk *baseToken) setNext(n Token) {
	tk.next = n
}

func (tk *baseToken) getPrev() Token {
	return tk.prev
}

func (tk *baseToken) setPrev(p Token) {
	tk.prev = p
}

func (tk *baseToken) copy() Token {
	panic(newRegexException("Unimplemented: always needs to be overriden"))
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

func (tk *baseToken) insertAfter(self Token, n Token) {
	if n == nil {
		return
	}

	head, last := copyList(n)
	if tk.getNext() != nil {
		last.setNext(tk.getNext())
		last.getNext().setPrev(last)
	}
	tk.setNext(head)
	head.setPrev(self)
}

func copyList(n Token) (Token, Token) {
	head := n.copy()
	prev := head

	next := n.getNext()
	for next != nil {
		tmp := next.copy()
		prev.setNext(tmp)
		tmp.setPrev(prev)
		prev = tmp
		next = next.getNext()
	}

	return head, prev
}

func (tk *baseToken) deleteUntil(self Token, n Token, m *matcher) {
	cur := tk.getNext()

	if n != nil {
		n.getPrev().setNext(nil)
	}

	self.setNext(n)
	if n != nil {
		n.setPrev(self)
	}

	for cur != nil && cur != n {
		delete(m.tokenState, cur)
		cur = cur.getNext()
	}
}

type nextState struct {
	myNext   Token
	startPos int
}
