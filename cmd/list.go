package cmd

import (
	"github.com/opb/secretly/kvp"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the secret keys and which Secret Manager config they came from",
	Long:  "List the secret keys and which Secret Manager config they came from",
	Run: func(cmd *cobra.Command, args []string) {
		kvp.PrintKeyListBySecret(args)
	},
}
