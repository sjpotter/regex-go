package regex

type baseStack struct {
	s []interface{}
}

func (s *baseStack) Len() int {
	return len(s.s)
}

func (s *baseStack) Peek() interface{} {
	if len(s.s) == 0 {
		return nil
	}

	return s.s[len(s.s)-1]
}

func (s *baseStack) Push(e interface{}) {
	s.s = append(s.s, e)
}

func (s *baseStack) Pop() interface{} {
	if len(s.s) == 0 {
		return nil
	}

	ret := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]

	return ret
}

type TokenStack struct {
	s baseStack
}

func (ts *TokenStack) Len() int {
	return ts.s.Len()
}

func (ts *TokenStack) Peek() Token {
	return ts.s.Peek().(Token)
}

func (ts *TokenStack) Pop() Token {
	return ts.s.Pop().(Token)
}

func (ts *TokenStack) Push(t Token) {
	ts.s.Push(t)
}

func (ts *TokenStack) Copy(f *TokenStack) {
	for _, v := range f.s.s {
		ts.s.s = append(ts.s.s, v)
	}
}

type StringStack struct {
	s baseStack
}

func (ss *StringStack) Len() int {
	return ss.s.Len()
}

func (ss *StringStack) Peek() string {
	ret := ss.s.Peek()
	if ret == nil {
		return ""
	}

	return ret.(string)
}

func (ss *StringStack) Pop() string {
	ret := ss.s.Pop()
	if ret == nil {
		return ""
	}

	return ret.(string)
}

func (ss *StringStack) Push(s string) {
	ss.s.Push(s)
}

func (ss *StringStack) Copy(f *StringStack) {
	for _, v := range f.s.s {
		ss.s.s = append(ss.s.s, v)
	}
}
