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
