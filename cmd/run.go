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

	"github.com/bilus/oya/cmd/internal"
	"github.com/bilus/oya/pkg/flags"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                "run TASK",
	Short:              "Runs an Oya task",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		taskName := args[0]
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		positionalArgs, flags, err := parseArgs(args)
		if err != nil {
			return err
		}
		return internal.Run(cwd, taskName, positionalArgs, flags, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func parseArgs(args []string) ([]string, map[string]string, error) {
	return flags.Parse(args[1:])
}
