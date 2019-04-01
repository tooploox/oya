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
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tooploox/oya/cmd/internal"
	"github.com/tooploox/oya/pkg/flags"
	"github.com/tooploox/oya/pkg/project"
	"github.com/tooploox/oya/pkg/task"
)

func execTask(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	taskName := cmd.Use
	cobraFlags, taskArgs, err := parseArgs(args)
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
	return internal.Run(cwd, taskName, taskArgs, recurse, changeset, cmd.OutOrStdout(), cmd.OutOrStderr())
}

func createCmd(name task.Name, desc string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                string(name),
		Short:              desc,
		Args:               cobra.ArbitraryArgs,
		SilenceUsage:       true,
		DisableFlagParsing: true,
		RunE:               execTask,
	}
	cmd.Flags().BoolP("recurse", "r", false, "Recursively process Oyafiles")
	cmd.Flags().BoolP("changeset", "c", false, "Use the Changeset: directives")
	return cmd
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	recurse := flagRecurse()
	changeset := flagChangeset()
	err = addTasksCommands(cwd, recurse, changeset)
	if err != nil {
		fmt.Println(err)
	}
}

func addTasksCommands(workDir string, recurse, changeset bool) error {
	installDir, err := project.InstallDir()
	if err != nil {
		return err
	}
	p, err := project.Detect(workDir, installDir)
	if err != nil {
		return err
	}
	err = p.InstallPacks()
	if err != nil {
		return err
	}
	oyafiles, err := p.RunTargets(workDir, recurse, changeset)
	if err != nil {
		return err
	}
	dependencies, err := p.Deps()
	if err != nil {
		return err
	}

	for _, o := range oyafiles {
		err = o.Build(dependencies)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		err = o.Tasks.ForEach(func(taskName task.Name, task task.Task, meta task.Meta) error {
			if !taskName.IsBuiltIn() {
				rootCmd.AddCommand(createCmd(taskName, meta.Doc))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func parseArgs(args []string) ([]string, internal.Args, error) {
	cobraFlags, rest := detectFlags(args)
	posArgs, flags, err := flags.Parse(rest)
	taskArgs := internal.Args{
		All:        rest,
		Positional: posArgs,
		Flags:      flags,
	}
	return cobraFlags, taskArgs, err
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

func flagRecurse() bool {
	re := regexp.MustCompile(`^-r$|^--recurse$`)
	return foundInArgs(re)
}

func flagChangeset() bool {
	re := regexp.MustCompile(`^-c$|^--changeset$`)
	return foundInArgs(re)
}

func foundInArgs(re *regexp.Regexp) bool {
	for _, arg := range os.Args {
		if re.MatchString(arg) {
			return true
		}
	}
	return false
}
