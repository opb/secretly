package main

import (
	"github.com/opb/secretly/kvp"
	"github.com/spf13/cobra"
)

var compareEnvFiles bool

func init() {
	RootCmd.AddCommand(compareCmd)
	compareCmd.Flags().BoolVarP(&compareEnvFiles, "use-files", "f", false, "Use files instead of AWS Secrets Manager")
}

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare the contents of multiple secrets",
	Long:  "Compare the contents of multiple secrets",
	Run: func(cmd *cobra.Command, args []string) {
		var provider kvp.Provider
		if compareEnvFiles {
			provider = kvp.FileProvider{}
		} else {
			provider = kvp.SMProvider{}
		}
		kvp.Compare(args, provider)
	},
}
