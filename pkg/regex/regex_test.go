package regex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/sjpotter/regex-go/pkg/regex"
)

var _ = Describe("Regex", func() {
	Context("Dot Match", func() {
		t := NewTokenizer(".")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})
		It("Single Digit", func() {
			ret, err := m.Match("1")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})
		It("Multiple Digits", func() {
			ret, err := m.Match("2135")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})
		It("Multiple Letters", func() {
			ret, err := m.Match("sdfdsfs")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})
		It("Empty String", func() {
			ret, err := m.Match("")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Start Anchor Test", func() {
		t := NewTokenizer("^\\d")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Leading Digit", func() {
			ret, err := m.Match("1asdfse")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("No Leading Digit", func() {
			ret, err := m.Match("a1sdfse")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("End Anchor Test", func() {
		t := NewTokenizer("\\d$")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Last Digit", func() {
			ret, err := m.Match("abcdsfsdfs223dsfsdaf32")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("No Last Digit", func() {
			ret, err := m.Match("abcdsfsdfs223dsfsdaf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Character Classes", func() {
		It("Individual characters - ^[123]*$", func() {
			t := NewTokenizer("^[123]*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("1232321")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("12324235")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		It("Individual characters negation - ^[^abc]*$", func() {
			t := NewTokenizer("^[^abc]*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("1234")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("1234a")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		It("Character Class Macros - ^[\\d]*$", func() {
			t := NewTokenizer("^[\\d]*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("8192389172")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("81923sdfsdf89172")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		It("Character Class Negation Macros - ^[^\\d]*$", func() {
			t := NewTokenizer("^[^\\d]*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("81923sdfasdf89172")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())

			ret, err = m.Match("asfsdkfjlkjl")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Character Class Ranges - ^[a-f]*$", func() {
			t := NewTokenizer("^[a-f]*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("abcdebdca")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("asdfewrsdf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		Context("Slash Characters", func() {
			It(`\\`, func() {
				t := NewTokenizer("^\\\\$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("\\")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("A")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\+`, func() {
				t := NewTokenizer("^\\+$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("+")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\*`, func() {
				t := NewTokenizer("^\\*$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("+")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\?`, func() {
				t := NewTokenizer("^\\?$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("?")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\^`, func() {
				t := NewTokenizer("^\\^$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("^")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\$`, func() {
				t := NewTokenizer("^\\$$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("$")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\.`, func() {
				t := NewTokenizer("^\\.$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match(".")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})

			It(`\-`, func() {
				t := NewTokenizer("^\\-$")
				m, err := NewMatcher(t)
				Expect(err).Should(BeNil())
				Expect(m).ShouldNot(BeNil())

				ret, err := m.Match("-")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeTrue())

				ret, err = m.Match("*")
				Expect(err).Should(BeNil())
				Expect(ret).Should(BeFalse())
			})
		})
	})

	Context("Invalid Regexes", func() {
		It("Invalid Character Class", func() {
			t := NewTokenizer("[asdf")
			_, err := NewMatcher(t)
			Expect(err).ShouldNot(BeNil())
		})

		It("Invalid Range Order", func() {
			t := NewTokenizer("[9-0]")
			_, err := NewMatcher(t)
			Expect(err).ShouldNot(BeNil())
		})

		It("Invalid Character", func() {
			t := NewTokenizer("123**")
			_, err := NewMatcher(t)
			Expect(err).ShouldNot(BeNil())
		})

		It("Invalid End Slash Character", func() {
			t := NewTokenizer("\\")
			_, err := NewMatcher(t)
			Expect(err).ShouldNot(BeNil())
		})

		It("Invalid Slash Character", func() {
			t := NewTokenizer("^\\P$")
			_, err := NewMatcher(t)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("Other Shortcuts", func() {
		It(`\D`, func() {
			t := NewTokenizer("^\\D*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("abcdfd")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("asdf324sdf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		It(`\w`, func() {
			t := NewTokenizer("^\\w*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("fdsaf08ewws34DWERdsf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("sdf324r;32fw`")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		It(`\W`, func() {
			t := NewTokenizer("^\\W*$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match(";';'`@!#$@")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("aefsa;213jkjsafs@!@$")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})

		It(`\S`, func() {
			t := NewTokenizer("^\\S+$")
			m, err := NewMatcher(t)
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())

			ret, err := m.Match("Sdfs324")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())

			ret, err = m.Match("sdfsd sdfsf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Different types of spaces", func() {
		t := NewTokenizer("\\s+")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Regular Space", func() {
			ret, err := m.Match("asdfd sdf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Tab Character", func() {
			ret, err := m.Match("sdf\tadfsd")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Newline Character", func() {
			ret, err := m.Match("sdf\n32rfsf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("No space Character", func() {
			ret, err := m.Match("Sdferfsd324")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Alternates", func() {
		t := NewTokenizer("^(abc|def|(hij*|kl*m)nop)qrs$")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Test1", func() {
			ret, err := m.Match("hijjnopqrs")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test2", func() {
			ret, err := m.Match("abcqrs")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test3", func() {
			ret, err := m.Match("kmnopqrs")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})
	})

	Context("Word Boundaries 1", func() {
		t := NewTokenizer("a(?>bc|b)c")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Test1", func() {
			ret, err := m.Match("abcc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test2", func() {
			ret, err := m.Match("abc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Word Boundaries 2", func() {
		t := NewTokenizer("a(?>c|b)c")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Test 1", func() {
			ret, err := m.Match("acc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test 2", func() {
			ret, err := m.Match("abc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test 3", func() {
			ret, err := m.Match("aac")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Look Ahead", func() {
		t := NewTokenizer("(?=regex)regex(abc)")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Verify Match", func() {
			ret, err := m.Match("regexabc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Verify Group 0", func() {
			ret := m.GetGroup(0)
			Expect(*ret).Should(Equal("regexabc"))
		})

		It("Verify Group 1", func() {
			ret := m.GetGroup(1)
			Expect(*ret).Should(Equal("abc"))
		})
	})

	Context("IfThenElse", func() {
		t := NewTokenizer("^(?(?=regex)(regex)|(abc))$")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Test 1", func() {
			ret, err := m.Match("regex")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test 1, Verify Group 1", func() {
			Expect(m.GetGroup(1)).Should(Equal("regex"))
		})

		It("Test 2", func() {
			ret, err := m.Match("abc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Test 2, Verify Group 2", func() {
			Expect(m.GetGroup(2)).Should(Equal("abc"))
		})

		It("Test 2, Verify All Groups", func() {
			ret := m.GetGroups()
			Expect(ret[0]).Should(Equal("abc"))
			Expect(ret[1]).Should(Equal("")) // Thoughts are this is incorrect as should be nil, not empty string
			Expect(ret[2]).Should(Equal("abc"))
		})
	})
})
