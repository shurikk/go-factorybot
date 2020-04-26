package factorybot_test

import (
	"fmt"

	. "factorybot"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sequence", func() {
	simple := NewSequence()

	Describe("#N", func() {
		It("increases sequence counter", func() {
			last := simple.N()
			for index := 1; index < 10; index++ {
				Expect(simple.N()).To(Equal(index + last))
			}
		})
	})

	Describe("#Rewind", func() {
		It("resets the sequence", func() {
			_ = simple.N()
			simple.Rewind()
			Expect(simple.N()).To(Equal(1))
		})
	})

	Describe("#One", func() {
		s := NewSequence(func(n int) interface{} {
			return fmt.Sprintf("address%d@example.com", n)
		})

		It("returns next value from the sequence", func() {
			Expect(s.One().(string)).To(Equal("address1@example.com"))
		})

		It("returns next value from the simple sequence", func() {
			Expect(simple.One()).To(Equal(simple.N() - 1))
		})
	})
})
