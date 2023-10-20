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

package ip2region

/* IP2Region Format (Little Endian)
+--------------------------------+--------------------------------+
|                      Header Chunk (256byte)                     |
+--------------------------------+--------------------------------+
|                  Vector Index Chunk (256*256byte)               |
+--------------------------------+--------------------------------+
|                           Data Chunk                            |
+--------------------------------+--------------------------------+
|                           Index Chunk                           |
+--------------------------------+--------------------------------+

Header Chunk (256byte)
+--------------------------------+--------------------------------+
| Version (2byte) | Cache Policy (2byte) |   Build Time (4byte)   |
+--------------------------------+--------------------------------+
|      Start Index (4byte)       |       End Index (4byte)        |
+--------------------------------+--------------------------------+
|                           Empty Data                            |
+--------------------------------+--------------------------------+

Vector Index Chunk (256*256byte)
+--------------------------------+--------------------------------+
|   Index Start Offset (4byte)   |    Index End Offset (4byte)    |
+--------------------------------+--------------------------------+

Index Chunk
+--------------------------------+--------------------------------+
|        Start IP (4byte)        |         End IP (4byte)         |
+--------------------------------+--------------------------------+
|       Data Length (2byte)      |       Data Offset (4byte)      |
+--------------------------------+--------------------------------+

Document: https://mp.weixin.qq.com/s/ndjzu0BgaeBmDOCw5aqHUg
*/
