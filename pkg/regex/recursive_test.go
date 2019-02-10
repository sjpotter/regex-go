package regex_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/sjpotter/regex-go/pkg/regex"
)

var _ = Describe("Recursive", func() {
	Context("Recursive Tests 1", func() {
		t := NewTokenizer("a(?R)?z")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Test 1", func() {
			ret, err := m.Match("aaazzz123")
			Expect(ret).Should(BeTrue())
			Expect(err).Should(BeNil())
			Expect(*m.GetGroup(0)).Should(Equal("aaazzz"))
		})

		It("Test 2", func() {
			ret, err := m.Match("aaabbzzz")
			Expect(ret).Should(BeFalse())
			Expect(err).Should(BeNil())
		})
	})

	Context("Recursive Tests 2", func() {
		t := NewTokenizer("(.)(?R)?(\\1)")
		m, err := NewMatcher(t)

		It("Verify Matcher Tokenization", func() {
			Expect(err).Should(BeNil())
			Expect(m).ShouldNot(BeNil())
		})

		It("Test", func() {
			ret, err := m.Match("abba")
			Expect(ret).Should(BeTrue())
			Expect(err).Should(BeNil())
			Expect(*m.GetGroup(0)).Should(Equal("abba"))
			Expect(*m.GetGroup(1)).Should(Equal("a"))
			Expect(*m.GetGroup(2)).Should(Equal("a"))
		})
	})
})
