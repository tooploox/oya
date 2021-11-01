package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
)

var importCmd = &cobra.Command{
	Use:   "import URI",
	Short: "Imports Oya pack in Oyafile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		alias, err := cmd.Flags().GetString("alias")
		if err != nil {
			return err
		}
		expose, err := cmd.Flags().GetBool("expose")
		if err != nil {
			return err
		}
		return internal.Import(cwd, args[0], alias, expose, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringP("alias", "a", "", "Import pack under alias name")
	importCmd.Flags().BoolP("expose", "e", false, "Expose imported tasks")
}
