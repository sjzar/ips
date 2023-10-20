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

package qqwry

/* qqwry.dat Format (Little Endian + GBK Encoding)
+--------------------------------+--------------------------------+
|       Start Index (4byte)      |       End Index (4byte)        |
+--------------------------------+--------------------------------+
|                            Data Chunk                           |
+--------------------------------+--------------------------------+
|                            Index Chunk                          |
+--------------------------------+--------------------------------+

Data Chunk
+--------------------------------+--------------------------------+
|                           End IP (4byte)                        |
+--------------------------------+--------------------------------+
|         Country (n byte)       |     End Flag 0x00 (1byte)      |
+--------------------------------+--------------------------------+
|          Area (n byte)         |     End Flag 0x00 (1byte)      |
+--------------------------------+--------------------------------+

Redirect Mode1: redirect country AND area
+--------------------------------+--------------------------------+
|                           End IP (4byte)                        |
+--------------------------------+--------------------------------+
|   Redirect Mode1 0x01 (1byte)  |      Data Offset (3byte)       |
+--------------------------------+--------------------------------+

Redirect Mode2: redirect country OR area
+--------------------------------+--------------------------------+
|                           End IP (4byte)                        |
+--------------------------------+--------------------------------+
|   Redirect Mode2 0x02 (1byte)  |      Data Offset (3byte)       |
+--------------------------------+--------------------------------+
|          Area (n byte)         |     End Flag 0x00 (1byte)      |
+--------------------------------+--------------------------------+

Index Chunk
+--------------------------------+--------------------------------+
|        Start IP (4byte)        |       Data Offset (3byte)      |
+--------------------------------+--------------------------------+

*/
