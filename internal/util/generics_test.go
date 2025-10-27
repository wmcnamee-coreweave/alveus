package util

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CoalesceSlices", func() {
	var (
		slices [][]string

		got []string
	)

	BeforeEach(func() {
		slices = nil
	})

	JustBeforeEach(func() {
		got = CoalesceSlices(slices...)
	})

	When("slices is nil", func() {
		BeforeEach(func() {
			slices = nil
		})

		It("should return a nil slice", func() {
			Expect(got).To(BeNil())
		})
	})

	When("slices is empty", func() {
		BeforeEach(func() {
			slices = [][]string{}
		})

		It("should return a nil slice", func() {
			Expect(got).To(BeNil())
		})
	})

	When("first slice has something in it", func() {
		BeforeEach(func() {
			slices = [][]string{
				{"a"},
			}
		})

		It("should return the first slice", func() {
			Expect(got).To(HaveLen(1))
			Expect(got[0]).To(Equal("a"))
		})
	})

	When("second slice has something in it", func() {
		BeforeEach(func() {
			slices = [][]string{
				{},
				{"b"},
			}
		})

		It("should return the second slice", func() {
			Expect(got).To(HaveLen(1))
			Expect(got[0]).To(Equal("b"))
		})
	})

	When("both slices have something in them", func() {
		BeforeEach(func() {
			slices = [][]string{
				{"a"},
				{"b"},
			}
		})

		It("should still return the first slice", func() {
			Expect(got).To(HaveLen(1))
			Expect(got[0]).To(Equal("a"))
		})
	})
})
