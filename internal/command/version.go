package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"unifai/internal/version"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print build version",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "unifai %s\n", version.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", version.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "built: %s\n", version.BuildDate)
		},
	}
}
