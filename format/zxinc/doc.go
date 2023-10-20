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

package zxinc

/* zxipv6wry.dat Format (Little Endian)
+--------------------------------+--------------------------------+
|         "IPDB" (4byte)         |         Version (2byte)        |
+--------------------------------+--------------------------------+
|      Offset Length (1byte)     |        IP Length (1byte)       |
+--------------------------------+--------------------------------+
|         Count (8byte)          |        Start Index (8byte)     |
+--------------------------------+--------------------------------+
|                            Data Chunk                           |
+--------------------------------+--------------------------------+
|                            Index Chunk                          |
+--------------------------------+--------------------------------+

Data Chunk like qqwry.dat, but the End IP not included and use UTF-8 encoding.

*/
