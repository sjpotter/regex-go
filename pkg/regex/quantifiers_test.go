package regex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/yucs/spotter/regex/pkg/regex"
)

var _ = Describe("Quantifiers", func() {
	Context("* Quantifier", func() {
		t := NewTokenizer("^\\d*$")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Empty text", func() {
			ret, err := m.Match("")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("All Digits", func() {
			ret, err := m.Match("123123")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Mixed Characters", func() {
			ret, err := m.Match("23432vbwef23142")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("? Quantifier", func() {
		t := NewTokenizer("^\\d?a$")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("No Digit", func() {
			ret, err := m.Match("a")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("One Digit", func() {
			ret, err := m.Match("1a")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Two Digits", func() {
			ret, err := m.Match("12a")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("+ Quantifier", func() {
		t := NewTokenizer("\\d+")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("All Digit", func() {
			ret, err := m.Match("123213")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Mix Digits and Text", func() {
			ret, err := m.Match("23432vbwef23142")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})
	})

	Context("{1} Quantifier", func() {
		t := NewTokenizer("\\d{1}")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("All Digits", func() {
			ret, err := m.Match("123213")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Mix Digits and Text", func() {
			ret, err := m.Match("12sdaf3213")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("All Text", func() {
			ret, err := m.Match("SAdfds")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("{5,} Quantifier", func() {
		t := NewTokenizer("\\d{5,}")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("5 Digits", func() {
			ret, err := m.Match("asdfds12345dsaf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("9 Digits", func() {
			ret, err := m.Match("asdfds123456789dsaf")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Only 4 Digits in a row", func() {
			ret, err := m.Match("12sdaf1234")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("{10,11} Quantifier", func() {
		t := NewTokenizer("\\d{10,11}")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("10 Digits", func() {
			ret, err := m.Match("1234567890abc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("5 Digits", func() {
			ret, err := m.Match("12345")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})

	Context("Invalid Quantifier (parsed as character tokens) {10,abc}", func() {
		t := NewTokenizer("\\d{10,abc}")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Matching characters", func() {
			ret, err := m.Match("1{10,abc}")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeTrue())
		})

		It("Trying to match 10 numbers", func() {
			ret, err := m.Match("11234567890,abc")
			Expect(err).Should(BeNil())
			Expect(ret).Should(BeFalse())
		})
	})
})
