package ips

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.0.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version " + Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
