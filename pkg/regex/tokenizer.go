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

func (t *tokenizer) Tokenize() Token {
	if t.t == nil {
		capturePos := t.captureCount
		t.captureCount++

		t.t = t.createCapturedExpressionToken(capturePos, 0, len(t.regex))
	}

	return t.t
}

func (t *tokenizer) tokenizeRange(regexPos, end int) Token {
	var token Token

	if regexPos >= end {
		return nil
	}

	switch t.regex[regexPos] {
	case '^', '$': // start of line anchor token, end of line anchor token
		token = newAnchorToken(t.regex[regexPos])
		regexPos++
	case '\\':
		if regexPos+1 < len(t.regex) {
			switch t.regex[regexPos+1] {
			case 'b', 'B':
				// word boundary anchor token
				token = newAnchorToken(t.regex[regexPos+1])
				regexPos += 2
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// back reference (match a previous group that the regex has stored)
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
				// other cases handled by characterToken below (i.e. things like /d /w)
			}
		}
	case '(': // There are many types of clauses that are within parens
		endParen := t.findMatchingParen(regexPos)

		if t.regex[regexPos+1] == '?' { // There are also many types of clauses that are within (? )
			switch t.regex[regexPos+2] {
			case '>':
				token = t.createAtomicExpressionToken(regexPos+3, endParen)

			// Look Ahead does not make sense to be quantified, position resets after they are done
			case '=': // Positive Look Ahead
				token = t.createLookAheadExpressionToken(regexPos+3, endParen, true)
			case '!': // Negative Look Ahead
				token = t.createLookAheadExpressionToken(regexPos+3, endParen, false)
			case '<':
				switch t.regex[regexPos+3] {
				case '=': // Positive Look Behind
					token = t.createLookBehindExpressionToken(regexPos+4, endParen, true)
				case '!': // Negative Look Behind
					token = t.createLookBehindExpressionToken(regexPos+4, endParen, false)
				default:
					panic(newRegexException("Unknown lookbehind grouping"))
				}
			case '(':
				token = t.createIfThenElseToken(regexPos+2, endParen)
			case 'R', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': //regex recursion
				token = t.createRecursiveToken(regexPos+2, endParen)
			default:
				panic(newRegexException("Unknown grouping type"))
			}
		} else { //normal capture
			capture := t.captureCount
			t.captureCount++

			token = t.createCapturedExpressionToken(capture, regexPos+1, endParen)
			token = newStartCaptureToken(capture, token)
		}
		regexPos = endParen + 1
	}

	if token == nil {
		token, regexPos = newCharacterToken(t.regex, regexPos)
	}

	if token.quantifiable() {
		var qt Token

		qt, regexPos = quantifierParse(token, t.regex, regexPos)
		if qt != nil {
			token = qt
		}
	}

	next := t.tokenizeRange(regexPos, end)
	token.setNext(next)
	if next != nil {
		next.setPrev(token)
	}

	return token
}

func (t *tokenizer) createRecursiveToken(regexPos, endParen int) Token {
	var capture int64

	if t.regex[regexPos] == 'R' && regexPos+1 == endParen {
		capture = 0
	} else {
		var err error

		capture, err = strconv.ParseInt(string(t.regex[regexPos:endParen]), 10, 64)
		if err != nil {
			panic(newRegexException(fmt.Sprintf("createRecursiveToken: couldn't parse %v as int", t.regex[regexPos:endParen])))
		}
	}

	return newRecursiveToken(int(capture))
}

func (t *tokenizer) createIfThenElseToken(regexPos, endParen int) Token {
	ifEndParen := t.findMatchingParen(regexPos)

	// Support testing if capture group exists.
	ifToken := t.tokenizeRange(regexPos, ifEndParen)
	var thenToken Token
	var elseToken Token

	if ifToken.normalExpression() {
		t.captureCount-- // TODO: HACK as the tokenize on the () string above would have incremented
		ifToken = newCaptureGroupTesterToken(t.regex[regexPos+1 : ifEndParen])
	}

	if !ifToken.testable() {
		panic(newRegexException(fmt.Sprintf("Didn't parse a testable token from: %v", t.regex[regexPos:ifEndParen])))
	}

	pipes := t.findPipes(ifEndParen+1, endParen)
	switch len(pipes) {
	case 0:
		thenToken = t.createNormalExpressionToken(ifEndParen+1, endParen)
		elseToken = nil
	case 1:
		thenToken = t.createNormalExpressionToken(ifEndParen+1, pipes[0])
		elseToken = t.createNormalExpressionToken(pipes[0]+1, endParen)
	default:
		panic(newRegexException("Expected at most one pipe in if/then/else token parsing"))
	}

	return newIfThenElseToken(ifToken, thenToken, elseToken)
}

func (t *tokenizer) createLookAheadExpressionToken(regexPos, endParen int, positive bool) Token {
	net := t.createNormalExpressionToken(regexPos, endParen)

	return newLookAheadExpressionToken(net, positive)
}

func (t *tokenizer) createLookBehindExpressionToken(regexPos, endParen int, positive bool) Token {
	net := t.createNormalExpressionToken(regexPos, endParen)

	return newLookBehindExpressionToken(net, positive)
}

func (t *tokenizer) createCapturedExpressionToken(capturePos int, regexPos int, endParen int) *normalExpresionToken {
	net := t.createNormalExpressionToken(regexPos, endParen)

	if capturePos != -1 {
		t.captureMap[capturePos] = net.copy().(*normalExpresionToken)
	}

	return net
}

func (t *tokenizer) createAtomicExpressionToken(regexPos, endParen int) *atomicExpressionToken {
	aet := newAtomicExpressionToken()
	t.parseExpression(aet.expressionToken, regexPos, endParen)

	return aet
}

func (t *tokenizer) createNormalExpressionToken(regexPos, endParen int) *normalExpresionToken {
	net := newNormalExpressionToken()

	t.parseExpression(net.expressionToken, regexPos, endParen)

	return net
}

func (t *tokenizer) parseExpression(et *expressionToken, regexPos, endParen int) {
	pipes := t.findPipes(regexPos, endParen)

	for _, pipe := range pipes {
		alt := t.tokenizeRange(regexPos, pipe)
		et.addAlt(alt)
		regexPos = pipe + 1
	}

	alt := t.tokenizeRange(regexPos, endParen)
	et.addAlt(alt)
}

// find the pipes that separate expressions at the same level within the start/end indices.
// if a section is enclosed in parens, its not at the same level, and hence not a pipe we care about
func (t *tokenizer) findPipes(start, end int) []int {
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
		panic(newRegexException("unbalanced parens"))
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
					panic(newRegexException("unbalanced parens"))
				}
				parens = parens[:len(parens)-1]
			}
		case '|':
			if (!slashIsEscape || t.regex[i-1] != '\\') && len(parens) == 0 {
				pipes = append(pipes, i)
			}
		}
	}

	return pipes
}

func (t *tokenizer) findMatchingParen(start int) int {
	if t.regex[start] != '(' {
		panic(newRegexException("findMatchingParen: didn't start with an open paren"))
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
				return i
			}
		}
	}

	panic(newRegexException("unbalanced parens"))
}