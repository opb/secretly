package main

import (
	"github.com/opb/secretly/kvp"
	"github.com/spf13/cobra"
)

var listEnvFiles bool

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listEnvFiles, "use-files", "f", false, "Use files instead of AWS Secrets Manager")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the secret keys and which Secret Manager config they came from",
	Long:  "List the secret keys and which Secret Manager config they came from",
	Run: func(cmd *cobra.Command, args []string) {
		var provider kvp.Provider
		if listEnvFiles{
			provider = kvp.FileProvider{}
		}else{
			provider = kvp.SMProvider{}
		}
		kvp.List(args, provider)
	},
}
