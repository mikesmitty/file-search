package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("file-search %s (%s) built at %s\n", Version, Commit, Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
