// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path/filepath"

	"github.com/bilus/oya/cmd/internal"
	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render TEMPLATE",
	Short: "Render a template using values from an Oyafile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		oyafilePath, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		templatePath := args[0]
		outputPath, err := os.Getwd()
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
		scopeSelector, err := cmd.Flags().GetString("scope")
		if err != nil {
			return err
		}
		return internal.Render(fullOyafilePath, templatePath, outputPath, autoScope, scopeSelector, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringP("file", "f", "Oyafile", "Read FILE as Oyafile")
	renderCmd.Flags().StringP("scope", "s", "", "Render template within the specified value scope")
	renderCmd.Flags().BoolP("auto-scope", "a", false, "Automatically detect value scope when running in imported tasks")
}
