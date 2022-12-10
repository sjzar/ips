/*
 * Copyright (c) 2022 shenjunzheng@gmail.com
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

package rewriter

import (
	"net"

	"github.com/sjzar/ips/ipx"
)

// EmptyRewriter 空改写
type EmptyRewriter struct {
}

// Rewrite 改写
func (r *EmptyRewriter) Rewrite(ip net.IP, ipRange *ipx.Range, data map[string]string) (net.IP, *ipx.Range, map[string]string, error) {
	return ip, ipRange, data, nil
}
