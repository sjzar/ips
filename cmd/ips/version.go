package ips

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var Version = "v0.0.1"
var BuildInfo = debug.BuildInfo{}
var (
	versionM bool
)

func init() {
	if bi, ok := debug.ReadBuildInfo(); ok {
		BuildInfo = *bi
		Version = bi.Main.Version
	}
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionM, "module", "m", false, "module version information")
}

var versionCmd = &cobra.Command{
	Use:   "version [-m]",
	Short: "Print the version of ips",
	Run: func(cmd *cobra.Command, args []string) {
		if versionM {
			mod := BuildInfo.String()
			if len(mod) > 0 {
				fmt.Printf("\t%s\n", strings.ReplaceAll(mod[:len(mod)-1], "\n", "\n\t"))
			}
		} else {
			fmt.Printf("ips version %s %s %s/%s\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		}
	},
}
