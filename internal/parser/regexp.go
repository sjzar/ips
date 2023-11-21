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

package parser

import (
	"regexp"
)

var (
	// IPv4Regexp IPv4 正则表达式
	IPv4Regexp = regexp.MustCompile(`(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(\.(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)){3}`)

	// IPv6Regexp IPv6 正则表达式
	// [fF][eE]80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?| # IPv6 Link-local (`net.ParseIP` does not support this format)
	// ([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}| # IPv6
	// ::([fF]{4}){1}:(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(\.(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)){3}| # IPv4-mapped IPv6 address
	// (([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})? # IPv6 with two colons
	IPv6Regexp = regexp.MustCompile(`([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|::([fF]{4}){1}:(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(\.(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)){3}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)

	// DomainRegexp 域名正则表达式
	// [a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,}
	// ([a-zA-Z0-9][-a-zA-Z0-9]{0,62}\.)+([a-zA-Z][-a-zA-Z]{0,62})
	// ^(xn--|_)?[a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?(\.(xn--|_)?[a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?)*(\.[a-zA-Z]{2,})$
	DomainRegexp = regexp.MustCompile(`(xn--|_)?[a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?(\.(xn--|_)?[a-zA-Z0-9]([a-zA-Z0-9-_]{0,61}[a-zA-Z0-9])?)*(\.[a-zA-Z]{2,})`)
)
