package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets in Oyafile.secrets",
}

var secretsViewCmd = &cobra.Command{
	Use:          "view",
	Short:        "View secrets",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		return internal.SecretsView(path,
			cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

var secretsEditCmd = &cobra.Command{
	Use:          "edit",
	Short:        "Edit secrets file",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		return internal.SecretsEdit(path,
			cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

var secretsEncryptCmd = &cobra.Command{
	Use:          "encrypt",
	Short:        "Encrypt secrets file",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		return internal.SecretsEncrypt(path,
			cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsEncryptCmd)
	rootCmd.AddCommand(secretsCmd)
}
