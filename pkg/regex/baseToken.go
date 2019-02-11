package regex

/* A regexp is a doubly linked list of tokens (enabling advancing and backtracking without a stack) that can be edited
   in place, so that compound tokens can insert their subtokens when advancing and remove them when backtracking.

   If match returns true, we will advance to the next token in the list (possibly modified by the token itself) and when
   it returns false, we will backtrack to the previous token.
 */
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

/* nextState is the standard state strut used by multiple tokens that enables it to record the current text position
   (for resetting if trying a different path available to the token) and myNext which is the token's current next token
   so it can know what to remove up to if it has to delete the tokens it has inserted
 */
type nextState struct {
	myNext   Token
	startPos int
}

// super() like function, doesn't do much now, but good to abstract it away to make it easy to modify
func newBaseToken() *baseToken {
	return &baseToken{}
}

// how we fake treat baseToken as an abstract class
func (tk *baseToken) match(m *matcher) bool {
	panic(newRegexException("match() unimplemented: always needs to be overriden"))
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

// more treating baseToken as an abstract class
func (tk *baseToken) copy() Token {
	panic(newRegexException("copy() unimplemented: always needs to be overriden"))
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

// clones a list of tokens by copying them
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

/* 
*  insertAfter() and deleteUntil() work together to enable compound tokens to maintain the linkedlist structure
   A compound token is a token that is composed of other "subtokens" (ex: expressions with alternates, quantifiers...)
   A compound token will insert its subtokens into the linked list when required but also remove them on failures.
   To ensure pristine tokens, we always copy the token nodes that are part of the subtoken list before inserting them.
   In many cases (notably quantifiers) the same token list can be inserted multiple times and hence needs to be copied
*/
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
