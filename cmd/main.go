package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "secretly",
	Short: "Load secrets from AWS SecretsManager (or local files) and load into the env",
	Long: `Use secretly as a prefix for a command to load secrets from
AWS SecretsManager or local files into the Env variable context for the specified command`,
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
