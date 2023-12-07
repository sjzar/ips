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

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mdnsCmd)

	// mdns
	mdnsCmd.Flags().StringVarP(&dnsClientNet, "net", "", "udp", "Specifies the network protocol to be used by the DNS client. tcp, udp, tcp-tls.")
	mdnsCmd.Flags().IntVarP(&dnsClientTimeoutMs, "client-timeout", "", 2000, "Defines the timeout in milliseconds for DNS client requests.")
	mdnsCmd.Flags().BoolVarP(&dnsClientSingleInflight, "single-inflight", "", false, "Indicates whether the DNS client should avoid making duplicate queries concurrently.")
	mdnsCmd.Flags().IntVarP(&mdnsTimeoutS, "timeout", "", 20, "Specifies the timeout in seconds for MDNS operations.")
	mdnsCmd.Flags().StringVarP(&mdnsExchangeAddress, "exchange-address", "", "119.29.29.29", "Defines the address of the DNS server to be used for MDNS queries.")
	mdnsCmd.Flags().IntVarP(&mdnsRetryTimes, "retry-times", "", 3, "Sets the number of times an MDNS query should be retried on failure.")

	// operate
	mdnsCmd.Flags().StringVarP(&fields, "fields", "f", "", UsageFields)
	mdnsCmd.Flags().BoolVarP(&useDBFields, "use-db-fields", "", false, UsageUseDBFields)
	mdnsCmd.Flags().StringVarP(&rewriteFiles, "rewrite-files", "r", "", UsageRewriteFiles)
	mdnsCmd.Flags().StringVarP(&lang, "lang", "", "", UsageLang)

	// database
	mdnsCmd.Flags().StringSliceVarP(&rootFile, "file", "i", nil, UsageQueryFile)
	mdnsCmd.Flags().StringSliceVarP(&rootFormat, "format", "", nil, UsageQueryFormat)
	mdnsCmd.Flags().StringSliceVarP(&rootIPv4File, "ipv4-file", "", nil, UsageQueryIPv4File)
	mdnsCmd.Flags().StringSliceVarP(&rootIPv4Format, "ipv4-format", "", nil, UsageQueryIPv4Format)
	mdnsCmd.Flags().StringSliceVarP(&rootIPv6File, "ipv6-file", "", nil, UsageQueryIPv6File)
	mdnsCmd.Flags().StringSliceVarP(&rootIPv6Format, "ipv6-format", "", nil, UsageQueryIPv6Format)
	mdnsCmd.Flags().StringVarP(&readerOption, "database-option", "", "", UsageReaderOption)
	mdnsCmd.Flags().StringVarP(&hybridMode, "hybrid-mode", "", "aggregation", UsageHybridMode)

	// output
	mdnsCmd.Flags().StringVarP(&rootTextValuesSep, "text-values-sep", "", "", UsageTextValuesSep)
}

var mdnsCmd = &cobra.Command{
	Use:   "mdns <domain>",
	Short: "Execute Multi-Geolocation DNS queries.",
	Long: `The 'ips mdns' command is designed to query domain name resolutions across multiple regions.

Utilizing EDNS capabilities, this command sends the client's IP address to the DNS server, which then returns the domain name resolution results for the corresponding region.

The command provides a global perspective on domain name resolutions, enabling users to quickly identify any anomalies in DNS resolutions.

For more detailed information and advanced configuration options, please refer to https://github.com/sjzar/ips/blob/main/docs/mdns.md`,
	PreRun: PreRunInit,
	Run:    MDNS,
}

func MDNS(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	ret, err := manager.MDNSResolve(args[0])
	if err != nil {
		return
	}
	fmt.Println(ret)
}
