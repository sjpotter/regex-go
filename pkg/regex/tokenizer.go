package regex

import (
	"fmt"
	"strconv"
	"unicode"
)

type tokenizer struct {
	regex        []rune
	captureCount int
	captureMap   map[int]*normalExpresionToken

	t Token
}

func NewTokenizer(r string) *tokenizer {
	return &tokenizer{
		regex:        []rune(r),
		captureCount: 0,
		captureMap:   make(map[int]*normalExpresionToken),
	}
}

func (t *tokenizer) Tokenize() (Token, *RegexException) {
	var err *RegexException

	if t.t == nil {
		capturePos := t.captureCount
		t.captureCount++

		t.t, err = t.createCapturedExpressionToken(capturePos, 0, len(t.regex))
		if err != nil {
			return nil, err
		}
	}

	return t.t, nil
}

func (t *tokenizer) tokenizeRange(regexPos, end int) (Token, *RegexException) {
	var token Token

	if regexPos >= end {
		return nullToken, nil
	}

	switch t.regex[regexPos] {
	case '^', '$': // start of line anchor token, end of line anchor token
		token = newAnchorToken(t.regex[regexPos])
		regexPos++
	case '\\':
		if regexPos+1 < len(t.regex) {
			// word boundary anchor token
			switch t.regex[regexPos+1] {
			case 'b', 'B':
				token = newAnchorToken(t.regex[regexPos+1])
				regexPos += 2
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				regexPos++
				val := int(t.regex[regexPos] - '0')
				for unicode.IsDigit(t.regex[regexPos+1]) {
					regexPos++
					val *= 10
					val += int(t.regex[regexPos] - '0')
				}

				token = newBackReferenceToken(val)
				regexPos++
			default:
				// other cases handled by characterToken below
			}
		}
	case '(': // There are many types of clauses that are within parens
		endParen, err := t.findMatchingParen(regexPos)
		if err != nil {
			return nil, err
		}

		if t.regex[regexPos+1] == '?' { // There are also many types of clauses that are within (? )
			switch t.regex[regexPos+2] {
			case '>':
				token, err = t.createAtomicExpressionToken(regexPos+3, endParen)

			// Look Ahead does not make sense to be quantified, position resets after they are done
			case '=': // Positive Look Ahead
				token, err = t.createLookAheadExpressionToken(regexPos+3, endParen, true)
			case '!': // Negative Look Ahead
				token, err = t.createLookAheadExpressionToken(regexPos+3, endParen, false)
			case '<':
				switch t.regex[regexPos+3] {
				case '=': // Positive Look Behind
					token, err = t.createLookBehindExpressionToken(regexPos+4, endParen, true)
				case '!': // Negative Look Behind
					token, err = t.createLookBehindExpressionToken(regexPos+4, endParen, false)
				default:
					return nil, newRegexException("Unknown lookbehind grouping")
				}
			case '(':
				token, err = t.createIfThenElseToken(regexPos+2, endParen)
			case 'R', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': //regex recursion
				token, err = t.createRecursiveToken(regexPos+2, endParen)
			default:
				return nil, newRegexException("Unknown grouping type")
			}
		} else { //normal capture
			capture := t.captureCount
			t.captureCount++
			token, err = t.createCapturedExpressionToken(capture, regexPos+1, endParen)
			token = newStartCaptureToken(capture, token)
		}
		regexPos = endParen+1
	}

	if token == nil {
		var err *RegexException

		if token, regexPos, err = newCharacterToken(t.regex, regexPos); err != nil {
			return nil, err
		}
	}

	if token.quantifiable() {
		var qt Token
		var err *RegexException

		qt, regexPos, err = quantifierParse(token, t.regex, regexPos)
		if err != nil {
			return nil, err
		}
		if qt != nil {
			token = qt
		}
	}

	next, err := t.tokenizeRange(regexPos, end)
	token.setNext(next)

	return token, err
}

func (t *tokenizer) createRecursiveToken(regexPos, endParen int) (Token, *RegexException) {
	var capture int64

	if t.regex[regexPos] == 'R' && regexPos+1 == endParen {
		capture = 0
	} else {
		var err error

		capture, err = strconv.ParseInt(string(t.regex[regexPos:endParen]), 10, 64)
		if err != nil {
			return nil, newRegexException(fmt.Sprintf("createRecursiveToken: couldn't parse %v as int", t.regex[regexPos:endParen]))
		}
	}

	return newRecursiveToken(int(capture)), nil
}

func (t *tokenizer) createIfThenElseToken(regexPos, endParen int) (Token, *RegexException) {
	ifEndParen, err := t.findMatchingParen(regexPos)
	if err != nil {
		return nil, err
	}

	// Support testing if capture group exists.
	ifToken, err := t.tokenizeRange(regexPos, ifEndParen)
	if err != nil {
		return nil, err
	}
	var thenToken Token
	var elseToken Token

	if ifToken.normalExpression() {
		t.captureCount-- // TODO: HACK as the tokenize on the () string above would have incremented
		ifToken, err = newCaptureGroupTesterToken(t.regex[regexPos+1:ifEndParen])
		if err != nil {
			return nil, err
		}
	}

	if !ifToken.testable() {
		return nil, newRegexException(fmt.Sprintf("Didn't parse a testable token from: %v", t.regex[regexPos:ifEndParen]))
	}

	pipes, err := t.findPipes(ifEndParen+1, endParen)
	if err != nil {
		return nil, err
	}
	switch len(pipes) {
	case 0:
		thenToken, err = t.createNormalExpressionToken(ifEndParen+1, endParen)
		if err != nil {
			return nil, err
		}
		elseToken = nullToken
	case 1:
		thenToken, err = t.createNormalExpressionToken(ifEndParen+1, pipes[0])
		if err != nil {
			return nil, err
		}
		elseToken, err = t.createNormalExpressionToken(pipes[0]+1, endParen)
		if err != nil {
			return nil, err
		}
	default:
		return nil, newRegexException("Expected at most one pipe in if/then/else token parsing")
	}

	return newIfThenElseToken(ifToken, thenToken, elseToken)
}

func (t *tokenizer) createLookAheadExpressionToken(regexPos, endParen int, positive bool) (Token, *RegexException) {
	net, err := t.createNormalExpressionToken(regexPos, endParen)
	if err != nil {
		return nil, err
	}

	return newLookAheadExpressionToken(net, positive), nil
}

func (t *tokenizer) createLookBehindExpressionToken(regexPos, endParen int, positive bool) (Token, *RegexException) {
	net, err := t.createNormalExpressionToken(regexPos, endParen)
	if err != nil {
		return nil, err
	}

	err = net.internalReverse()
	if err != nil {
		return nil, err
	}

	return newLookBehindExpressionToken(net, positive), nil
}

func (t *tokenizer) createCapturedExpressionToken(capturePos int, regexPos int, endParen int) (*normalExpresionToken, *RegexException) {
	net, err := t.createNormalExpressionToken(regexPos, endParen)
	if err != nil {
		return nil, err
	}

	if capturePos != -1 {
		t.captureMap[capturePos] = net
	}

	return net, nil
}

func (t *tokenizer) createAtomicExpressionToken(regexPos, endParen int) (*atomicExpressionToken, *RegexException) {
	aet := newAtomicExpressionToken()
	if err := t.parseExpression(aet.expressionToken, regexPos, endParen); err != nil {
		return nil, err
	}

	return aet, nil
}

func (t *tokenizer) createNormalExpressionToken(regexPos, endParen int) (*normalExpresionToken, *RegexException) {
	net := newNormalExpressionToken()

	if err := t.parseExpression(net.expressionToken, regexPos, endParen); err != nil {
		return nil, err
	}

	return net, nil
}

func (t *tokenizer) parseExpression(et *expressionToken, regexPos, endParen int) *RegexException {
	pipes, err := t.findPipes(regexPos, endParen)

	for _, pipe := range pipes {
		alt, err := t.tokenizeRange(regexPos, pipe)
		if err != nil {
			return err
		}
		et.addAlt(alt)
		regexPos = pipe + 1
	}

	alt, err := t.tokenizeRange(regexPos, endParen)
	if err != nil {
		return err
	}
	et.addAlt(alt)

	return nil
}

// find the pipes that separate expressions at the same level within the start/end indices.
// if a section is enclosed in parens, its not at the same level, and hence not a pipe we care about
func (t *tokenizer) findPipes(start, end int) ([]int, *RegexException) {
	var parens []int
	var pipes []int

	var slashIsEscape = false

	switch t.regex[start] {
	case '\\':
		slashIsEscape = true
	case '|':
		pipes = append(pipes, start)
	case '(':
		parens = append(parens, start)
	case ')':
		return nil, newRegexException("unbalanced parens")
	}

	// As we are searching for alternates, only find pipes that are not in sub expressions (i.e. surrounded by ()
	for i := start + 1; i < end; i++ { //last element should be a paren
		switch t.regex[i] {
		case '\\':
			slashIsEscape = !slashIsEscape
		case '(':
			if !slashIsEscape || t.regex[i-1] != '\\' {
				parens = append(parens, i)
			}
		case ')':
			if !slashIsEscape || t.regex[i-1] != '\\' {
				if len(parens) == 0 {
					return nil, newRegexException("unbalanced parens")
				}
				parens = parens[:len(parens)-1]
			}
		case '|':
			if (!slashIsEscape || t.regex[i-1] != '\\') && len(parens) == 0 {
				pipes = append(pipes, i)
			}
		}
	}

	return pipes, nil
}

func (t *tokenizer) findMatchingParen(start int) (int, *RegexException) {
	if t.regex[start] != '(' {
		return -1, newRegexException("findMatchingParen: didn't start with an open paren")
	}

	var parens []int

	// by reading the first character first, can make the switch in loop simpler as don't have to check if index 0
	parens = append(parens, start)

	slashIsEscape := false
	for i := start + 1; i < len(t.regex); i++ {
		switch t.regex[i] {
		case '\\':
			slashIsEscape = !slashIsEscape
		case '(':
			if !slashIsEscape || t.regex[i-1] != '\\' {
				parens = append(parens, i)
			}
		case ')':
			if !slashIsEscape || t.regex[i-1] != '\\' {
				parens = parens[:len(parens)-1]
			}
			if len(parens) == 0 {
				return i, nil
			}
		}
	}

	return -1, newRegexException("unbalanced parens")
}
