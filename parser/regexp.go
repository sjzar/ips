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

package parser

import (
	"regexp"
)

var (
	// IPv4Regexp IPv4正则表达式
	IPv4Regexp = regexp.MustCompile(`(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(\.(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)){3}`)

	// IPv6Regexp IPv6正则表达式
	// [fF][eE]80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?| # IPv6 Link-local
	// ([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}| # IPv6
	// ::([fF]{4}){1}:(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(\.(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)){3}| # IPv4-mapped IPv6 address
	// (([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})? # IPv6 with two colons
	IPv6Regexp = regexp.MustCompile(`[fF][eE]80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|::([fF]{4}){1}:(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(\.(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)){3}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)
)
