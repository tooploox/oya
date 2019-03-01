// Copyright © 2019 Marcin Bilski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
