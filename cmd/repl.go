package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start REPL",
	Long: `Start a REPL session, allowing you
to build Oya tasks interactively, evaluating shell commands
in an identical environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.REPL(cmd.InOrStdin(), cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
