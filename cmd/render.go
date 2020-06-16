package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
	"github.com/tooploox/oya/pkg/errors"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:          "render TEMPLATE",
	Short:        "Render a template FILE or DIRECTORY using values from an Oyafile",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		oyafilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		templatePath := args[0]
		outputPath, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			return err
		}
		// This will turn "." or empty output path into full path relative to pwd.
		fullOutputPath, err := filepath.Abs(outputPath)
		if err != nil {
			return err
		}
		fullOyafilePath, err := filepath.Abs(oyafilePath)
		if err != nil {
			return err
		}
		autoScope, err := cmd.Flags().GetBool("auto-scope")
		if err != nil {
			return err
		}
		scopePath, err := cmd.Flags().GetString("scope")
		if err != nil {
			return err
		}
		exclude, err := cmd.Flags().GetStringArray("exclude")
		if err != nil {
			return err
		}

		overrides, err := parseValueOverrides(cmd, "set")
		if err != nil {
			return err
		}

		delimiters, err := cmd.Flags().GetString("delimiters")
		if err != nil {
			return err
		}

		return internal.Render(fullOyafilePath, templatePath, exclude, fullOutputPath, autoScope, scopePath, overrides, delimiters, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringP("file", "f", "./Oyafile", "Path to Oyafile to read")
	renderCmd.Flags().StringP("output-dir", "o", ".", "Specify the output DIRECTORY")
	renderCmd.Flags().StringP("scope", "s", "", "Render template within the specified value scope")
	renderCmd.Flags().StringP("delimiters", "d", "<%...%>", "Render template using delimiters for logic block")
	renderCmd.Flags().BoolP("auto-scope", "a", true, "When running in an imported pack's task, use the pack's scope, unless --")
	renderCmd.Flags().StringArrayP("exclude", "e", []string{}, "Relative paths to files or directories to exclude")
	renderCmd.Flags().StringArrayP("set", "", []string{}, "Value overrides, e.g. foo.bar=value")
}

func parseValueOverrides(cmd *cobra.Command, flag string) (map[string]interface{}, error) {
	assigns, err := cmd.Flags().GetStringArray("set")
	if err != nil {
		return nil, err
	}
	overrides := make(map[string]interface{})
	for _, a := range assigns {
		if a == "" {
			// This will not happen in real-life but it does happen in godog tests
			// because cmd.ResetFlags function used in oya_test.go doesn't reset
			// string arrays properly, making them into [""].
			continue
		}
		pv := strings.Split(a, "=")
		if len(pv) != 2 {
			return nil, errors.Errorf("unexpected --%s %s flag; must be in key=value or some.path=value format", flag, a)
		}
		path := pv[0]
		value := pv[1]
		overrides[path] = value
	}
	return overrides, nil
}
