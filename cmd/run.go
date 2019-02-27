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
	"fmt"
	"os"
	"strings"

	"github.com/bilus/oya/cmd/internal"
	"github.com/bilus/oya/pkg/flags"
	"github.com/spf13/cobra"
)

type ErrMissingTaskName struct{}

func (e ErrMissingTaskName) Error() string {
	return fmt.Sprintf("missing TASK name")
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                "run TASK",
	Short:              "Runs an Oya task",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		cobraFlags, taskName, positionalArgs, flags, err := parseArgs(args)
		if err != nil {
			return err
		}
		// BUG(bilus): Yack. This is what has to be done to support arbitrary flags passed to tasks.
		cmd.DisableFlagParsing = false
		defer func() { cmd.DisableFlagParsing = true }()
		if err := cmd.ParseFlags(cobraFlags); err != nil {
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
		return internal.Run(cwd, taskName, recurse, changeset, positionalArgs, flags, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("recurse", "r", false, "Recursively process Oyafiles")
	runCmd.Flags().BoolP("changeset", "c", false, "Use the Changeset: directives")
}

func parseArgs(args []string) ([]string, string, []string, map[string]string, error) {
	cobraFlags, rest := detectFlags(args)
	if len(rest) == 0 {
		return nil, "", nil, nil, ErrMissingTaskName{}
	}
	taskName := rest[0]
	posArgs, flags, err := flags.Parse(rest[1:])
	return cobraFlags, taskName, posArgs, flags, err
}

// detectFlags processed args consisting of flags followed by positional arguments, splitting them.
// For example this: ["--foo", "-b", "xxx", "--foo"] becomes: ["--foo", "-b"], ["xxx", "--foo"].
func detectFlags(args []string) ([]string, []string) {
	flags := make([]string, 0, len(args))
	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			return flags, args[i:]
		}
	}
	return flags, nil
}
