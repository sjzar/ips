/*
 * Copyright (c) 2023 shenjunzheng@gmail.com
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ips

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	version   = "v0.0.1"
	buildInfo = debug.BuildInfo{}
	versionM  bool
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionM, "module", "m", false, "module version information")
}

var versionCmd = &cobra.Command{
	Use:   "version [-m]",
	Short: "Show the version of ips",
	PreRun: func(cmd *cobra.Command, args []string) {
		if bi, ok := debug.ReadBuildInfo(); ok {
			buildInfo = *bi
			version = bi.Main.Version
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if versionM {
			mod := buildInfo.String()
			if len(mod) > 0 {
				fmt.Printf("\t%s\n", strings.ReplaceAll(mod[:len(mod)-1], "\n", "\n\t"))
			}
		} else {
			fmt.Printf("ips version %s %s %s/%s\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		}
	},
}
