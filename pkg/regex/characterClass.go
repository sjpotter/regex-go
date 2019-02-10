package regex

import "fmt"

var (
	digitSet = make(map[rune]bool)
	lowerSet = make(map[rune]bool)
	upperSet = make(map[rune]bool)
	wordSet  = make(map[rune]bool)
	whiteSet = make(map[rune]bool)
)

func init() {
	for c := '0'; c <= '9'; c++ {
		digitSet[c] = true
		wordSet[c] = true
	}
	for c := 'a'; c <= 'z'; c++ {
		lowerSet[c] = true
		wordSet[c] = true
	}

	for c := 'A'; c <= 'Z'; c++ {
		upperSet[c] = true
		wordSet[c] = true
	}

	whiteSet[' '] = true
	whiteSet['\t'] = true
	whiteSet['\r'] = true
	whiteSet['\n'] = true
	whiteSet['\f'] = true
}

type characterClass struct {
	characters map[rune]bool
	negated    map[rune]bool
	all        bool
}

func newCharacterClass(regex []rune, beg, end int) *characterClass {
	i := beg
	negate := false

	cc := &characterClass{
		characters: make(map[rune]bool),
		negated:    make(map[rune]bool),
	}

	if end > beg && regex[beg] == '^' {
		i++
		negate = true
	}

	for ; i <= end; i++ {
		if regex[i] == '\\' {
			cc.parseSlash(negate, regex, i)
			i++
		} else if i+2 <= end && regex[i+1] == '-' {
			cc.parseRange(negate, regex, i)
			i += 2
		} else {
			if !negate {
				cc.characters[regex[i]] = true
			} else {
				cc.negated[regex[i]] = true
			}
		}
	}

	return cc
}

func allCharacters() *characterClass {
	return &characterClass{all: true}
}

func (cc *characterClass) parseRange(negate bool, regex []rune, pos int) {
	if regex[pos] < regex[pos+2] {
		for c := regex[pos]; c <= regex[pos+2]; c++ {
			if !negate {
				cc.characters[c] = true
			} else {
				cc.negated[c] = true
			}
		}
	} else {
		panic(newRegexException("Character class ranged have to be in ascending order: " + string(regex[pos:pos+3])))
	}
}

func (cc *characterClass) parseSlash(negate bool, regex []rune, regexPos int) {
	if len(regex) <= regexPos+1 {
		panic(newRegexException(`Cannot end a regex with a single \`))
	}

	character := []rune(regex)[regexPos+1]
	switch character {
	case '.', '-', '\\', '+', '*', '?', '^', '$', '|', '(', ')':
		cc.characters[character] = true
	case 'd':
		cc.addSet(negate, digitSet)
	case 'D':
		cc.addSet(!negate, digitSet)
	case 'w':
		cc.addSet(negate, wordSet)
	case 'W':
		cc.addSet(!negate, wordSet)
	case 's':
		cc.addSet(negate, whiteSet)
	case 'S':
		cc.addSet(!negate, whiteSet)
	default:
		panic(newRegexException(fmt.Sprintf("parseSlash: unknown slash case: %v at index: %v", regex[regexPos+1:regexPos+2], regexPos+1)))
	}
}

func (cc *characterClass) addSet(negate bool, set map[rune]bool) {
	for k, v := range set {
		if negate {
			cc.negated[k] = v
		} else {
			cc.characters[k] = v
		}
	}
}

func (cc *characterClass) match(r rune) bool {
	return cc.characters[r] || len(cc.negated) > 0 && !cc.negated[r] || cc.all
}
