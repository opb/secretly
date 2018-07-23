package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "unknown"

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of secretly",
	Long: "Print the version number of secretly",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("secretly - version "+Version)
	},
}