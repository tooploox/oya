// Copyright © 2018 Marcin Bilski
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
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tooploox/oya/pkg/oyafile"
	"github.com/tooploox/oya/pkg/task"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oya",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleError(err)
	}
}

func SetOyaVersion(ver string) {
	rootCmd.Version = ver
}

// ExecuteE executes a command same as Execute but returns error.
func ExecuteE() error {
	_, err := rootCmd.ExecuteC()
	return err
}

// SetOutput overrides cobra output (for testing).
func SetOutput(out io.Writer) {
	rootCmd.SetOutput(out)
}

// Reset all flags for all commands to their default values (testing).
func ResetFlags() {
	resetFlagsRecurse(rootCmd)
}

func resetFlagsRecurse(cmd *cobra.Command) {
	cmd.Flags().VisitAll(
		func(flag *pflag.Flag) {
			if flag.DefValue == "[]" {
				// BUG(bilus): I don't know how to set default value for StringArray flag type.
				flag.Value.Set("")

			} else {
				flag.Value.Set(flag.DefValue)
			}

		})
	for _, child := range cmd.Commands() {
		resetFlagsRecurse(child)
	}
}

func handleError(err error) {
	switch err := err.(type) {
	case oyafile.ErrTaskFail:
		handleTaskFail(err)
	default:
		logrus.Println(err)
		os.Stderr.WriteString("Error: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}

func handleTaskFail(err oyafile.ErrTaskFail) {
	fmt.Fprintf(os.Stderr, "--- RUN ERROR ---------------------- %v\n", err.OyafilePath)
	var showArgs string
	if len(err.Args) > 0 {
		showArgs = fmt.Sprintf(" invoked with arguments %q", strings.Join(err.Args, " "))
	}
	fmt.Fprintf(os.Stderr, "Error in task %q%v: %v\n", string(err.TaskName), showArgs, err.Cause.Error())
	if err.ImportPath != nil {
		fmt.Fprintf(os.Stderr, "  (task imported from %v)\n", *err.ImportPath)
	}
	switch cause := err.Cause.(type) {
	case task.ErrScriptFail:
		fmt.Fprintf(os.Stderr, "\n%s\n\n", pinpointScriptFail(cause))
		os.Exit(cause.ExitCode)
	default:
		os.Exit(1)
	}
}

func pinpointScriptFail(err task.ErrScriptFail) string {
	output := make([]string, 0, 3)
	lines := strings.Split(err.Script, "\n")
	var start uint = 1
	if err.Line > 1 {
		start = err.Line - 1
	}
	digits := int(math.Log10(float64(err.Line)) + 1)
	lineFmt := fmt.Sprintf("%%v %%%vv %%v", digits)
	for i := start; i <= err.Line; i++ {
		var marker string
		if i == err.Line {
			marker = ">"
		} else {
			marker = " "
		}
		output = append(output, fmt.Sprintf(lineFmt, marker, i, lines[i-1]))
	}
	output = append(output, fmt.Sprintf(lineFmt, " ", " ", colMarker(err.Column)))
	return strings.Join(output, "\n")
}

func colMarker(col uint) string {
	return strings.Repeat(" ", int(col)-1) + "^"
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oya.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".oya" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".oya")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
