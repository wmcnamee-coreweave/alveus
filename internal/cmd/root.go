package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "generate NOUN",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(
		NewInitCommand(),
		NewGenerateCommand(),
	)

	return cmd
}
