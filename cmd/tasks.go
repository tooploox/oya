package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
)

// tasksCmd represents the init command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List available tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		recurse, err := cmd.Flags().GetBool("recurse")
		if err != nil {
			return err
		}
		changeset, err := cmd.Flags().GetBool("changeset")
		if err != nil {
			return err
		}
		return internal.Tasks(cwd, recurse, changeset, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.Flags().BoolP("recurse", "r", false, "Recursively process Oyafiles")
	tasksCmd.Flags().BoolP("changeset", "c", false, "Use the Changeset: directives")
}
