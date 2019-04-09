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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oya",
	Short: "Oya is a task manager and runner",
	Long:  "Oya takes the pain out of bootstrapping new deployable projects with packaged boilerplate & scripts.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oya.yaml)")
	rootCmd.PersistentFlags().BoolP("recurse", "r", false, "Recursively process Oyafiles")
	rootCmd.PersistentFlags().BoolP("changeset", "c", false, "Recursively process Oyafiles")
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
