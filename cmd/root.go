package cmd

import (
	"fmt"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tooploox/oya/cmd/internal"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oya",
	Short: "Oya is a task manager and runner",
	Long: `Oya takes the pain out of bootstrapping new deployable
projects with packaged boilerplate & scripts.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if exitCode, err := ExecuteE(); err != nil {
		os.Exit(exitCode)
	}
}

func SetOyaVersion(ver string) {
	rootCmd.Version = ver
}

// ExecuteE executes a command same as Execute but returns error.
func ExecuteE() (int, error) {
	_, err := rootCmd.ExecuteC()
	if err != nil {
		return internal.HandleError(rootCmd.OutOrStderr(), err), err
	}
	return 0, nil
}

// SetInput overrides cobra input (for testing).
func SetInput(in io.Reader) {
	rootCmd.SetIn(in)
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
