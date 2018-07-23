package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "secretly",
	Short: "Grab secrets from AWS SecretsManager and load into the env",
	Long: `Use secretly as a prefix for a command to load secrets from
AWS SecretsManager into the Env variable context for the specified command`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
