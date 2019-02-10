package regex

type Token interface {
	match(m *matcher) bool
	getNext() Token
	setNext(Token)
	getPrev() Token
	setPrev(Token)
	reverse() Token
	reverseToken(cur Token) Token
	captureGroup() int
	copy() Token
	quantifiable() bool
	testable() bool
	normalExpression() bool
	delete(Token)
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

func (tk *baseToken) reverse() Token {
	return tk.reverseToken(tk)
}

func (tk *baseToken) reverseToken(cur Token) Token {
	if cur.getNext() != nil {
		return cur
	}

	prev := cur.getNext().reverse()
	cur.setNext(nil)

	tmp := prev
	for tmp.getNext() != nil {
		tmp = tmp.getNext()
	}
	tmp.setNext(tk)

	return prev
}

func (tk *baseToken) captureGroup() int {
	return -1
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
/* 	if cur == n {
		return
	}
*/
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

func (tk *baseToken) delete(self Token) {
	prev := self.getPrev()
	next := self.getNext()
	if prev != nil {
		prev.setNext(next)
	}
	if next != nil {
		next.setPrev(prev)
	}
}

type nextState struct {
	myNext   Token
	startPos int
}
