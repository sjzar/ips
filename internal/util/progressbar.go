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

package util

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/sjzar/ips/ipnet"
)

// SetIPProgressBar sets the ip to the progress bar.
// FIXME Less effective, characteristics temporarily deactivated (´･ω･`)
func SetIPProgressBar(bar *progressbar.ProgressBar, ip net.IP) {
	_ = bar.Set64(int64(ipnet.IPToUint32(ip)))
}

// ProgressBar returns a new progress bar.
// copy from progressbar.DefaultBytes, add OptionUseANSICodes(true)
// ISSUE: https://github.com/schollz/progressbar/issues/102
func ProgressBar(maxBytes int64, description ...string) *progressbar.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	return progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionUseANSICodes(true),
	)
}
