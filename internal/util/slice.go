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

// DeleteEmptyValue removes empty strings from a slice and returns a new slice.
func DeleteEmptyValue(s []string) []string {
	ret := make([]string, 0, len(s))
	for _, str := range s {
		if str != "" {
			ret = append(ret, str)
		}
	}
	return ret
}
