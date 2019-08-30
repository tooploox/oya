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

var secretsInitCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize secret management",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		typ, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		desc, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}
		return internal.SecretsInit(typ, email, name, desc, format,
			cmd.OutOrStdout(), cmd.OutOrStderr())
	},
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
		return internal.SecretsEdit(args[0],
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
	secretsInitCmd.Flags().StringP("type", "t", "pgp", "Key st TODO")
	secretsInitCmd.Flags().StringP("email", "e", "", "Email address to use to generate the key pair")
	secretsInitCmd.Flags().StringP("name", "n", "", "Name to use to generate the key pair")
	secretsInitCmd.Flags().StringP("description", "d", "", "Key pair description")
	secretsInitCmd.Flags().StringP("format", "f", "text", "Output format (text/json)")
	secretsCmd.AddCommand(secretsInitCmd)
	secretsCmd.AddCommand(secretsViewCmd)
	secretsCmd.AddCommand(secretsEditCmd)
	secretsCmd.AddCommand(secretsEncryptCmd)
	rootCmd.AddCommand(secretsCmd)
}
