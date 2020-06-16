package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get URI",
	Short: "Get Oya pack from external repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		update, err := cmd.Flags().GetBool("update")
		if err != nil {
			return err
		}
		return internal.Get(cwd, args[0], update, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("update", "u", false,
		"Update the package to the latest available version")
	// NOTE: --update flag mentioned in one of the error messages
}
