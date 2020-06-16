package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init PROJECT_NAME",
	Short: "Initialize an Oya project",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return internal.Init(cwd, args[0], cmd.OutOrStdout(), cmd.OutOrStderr())
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(initCmd)
}
