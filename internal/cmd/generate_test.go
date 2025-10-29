package cmd

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("NewGenerateCommand", func() {
	var (
		cmd       *cobra.Command
		ctx       context.Context
		args      []string
		actualErr error
	)

	BeforeEach(func() {
		ctx = context.Background()
		args = []string{}
		actualErr = nil
	})

	JustBeforeEach(func() {
		cmd = NewGenerateCommand()
		Expect(cmd).NotTo(BeNil())
		cmd.SetArgs(args)
		actualErr = cmd.ExecuteContext(ctx)
	})

	When("no arguments are passed", func() {
		BeforeEach(func() {
			args = []string{}
		})

		It("should err", func() {
			Expect(actualErr).To(MatchError(`required flag(s) "repo-url" not set`))
		})
	})

	When("repo-url is set", func() {
		BeforeEach(func() {
			args = append(args, "--repo-url", "https://github.com/wmcnamee-coreweave/alveus.git")
		})

		It("should err as service definition passed in (filename or stdin)", func() {
			Expect(actualErr).To(MatchError(`stdin is from a terminal`))
		})

		When("service definition passed in via stdin", func() {
			// TODO need to not rely explicitly on os.Stdin
		})
	})
})
