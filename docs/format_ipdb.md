# IPDB 格式数据库

<!-- TOC -->
* [IPDB 格式数据库](#ipdb-格式数据库)
  * [简介](#简介)
  * [格式分析](#格式分析)
    * [MetaData 分析](#metadata-分析)
  * [NodeChunk 分析](#nodechunk-分析)
    * [DataChunk 分析](#datachunk-分析)
  * [查询操作](#查询操作)
  * [打包流程](#打包流程)
<!-- TOC -->

## 简介

IPDB 数据库是由 [IPIP.net](https://www.ipip.net/) 设计并使用的一种 IP 数据库格式，以体积小、查询效率高、支持多语言等特点而知名。（高老师牛逼～）

## 格式分析

```shell
    +--------------------------------+--------------------------------+
    |    MetaData Length (4byte)     |     MetaData (Json Format)     |
    +--------------------------------+--------------------------------+
    |                 Node Chunk (Prefix Tree / Trie)                 |
    +--------------------------------+--------------------------------+
    |                            Data Chunk                           |
    +--------------------------------+--------------------------------+
```

* 文件分为三个部分：MetaData、NodeChunk、DataChunk。

### MetaData 分析

MetaData 是一个 JSON 对象，包含数据库构建信息和查询所需的元数据。以下是一个示例：

```shell
{
    "build": 1632971142,    // 构建时间
    "ip_version": 1,        // IP库版本 IPv4:0x1 IPv6:0x2
    "languages": {
        "CN": 0             // 语言 & 字段偏移量
    },
    "node_count": 8705098,  // Node数量
    "total_size": 90028407, // NodeChunk+DataChunk 数据大小
    "fields": [             // DataChunk中每组数据的字段
        "country_name",
        "region_name",
        "city_name",
        "owner_domain",
        "isp_domain",
        "latitude",
        "longitude",
        "timezone",
        "utc_offset",
        "china_admin_code",
        "idd_code",
        "country_code",
        "continent_code"
    ]
}
```

## NodeChunk 分析

NodeChunk 由前缀树（字典树）构成。每个 Node 为 8 字节，存储到下一个节点的偏移量。如果偏移量超过节点数，则跳转到 DataChunk，表示找到了结果。

### DataChunk 分析

DataChunk 存储 IP 数据库数据。相同的数据仅存储一次，以减少冗余。

```shell
    +--------------------------------+--------------------------------+--------------------------------+
    | Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
    +--------------------------------+--------------------------------+--------------------------------+
    | Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
    +--------------------------------+--------------------------------+--------------------------------+
```

* DataChunk 在数据块中，数据分为长度和数据两部分。
* 数据部分使用 `\t` 分隔字段，在多语言版本的数据库中使用字段偏移返回不同语言的数据。

## 查询操作

* CIDR 地址是一种结合了 IP 地址和子网掩码的网段描述方式，如 10.0.0.1/8 表示一个 8 位子网掩码（255.0.0.0）。CIDR 网段内的所有 IP 在子网掩码部分完全相同。
* 将 IP 地址视为一个 32位/128位 的二进制字符串，在节点块中使用前缀树从前往后查询，一旦找到 CIDR 匹配，就跳转到数据块返回相应数据。
* 更多的 CIDR 分组会导致更大的节点数。
* CIDR 不得相互包含，嵌套的 CIDR（如 10.0.0.1/8 和 10.0.0.1/16）可能会因为先匹配到的原因而阻止进一步匹配。
* 关于查询过程的详细解释，请参考论文 [IPv4 route lookup on Linux](https://vincent.bernat.ch/en/blog/2017-ipv4-route-lookup-linux)。

## 打包流程

* 构建前缀树，并按照加载顺序将不同的数据集放入数据块。
* 根据 IPv6 的规范构建前缀树，IPv4 数据需要填充前 96 位的子网数据，即 80 位的 0 和 16 位的 1，对应于 IPv6 的映射地址（::FFFF:<IPv4>）。查询 IPv4 时，这允许快速偏移到 96 位掩码位置开始查询。
* 理论上，ipdb 格式数据库支持将 IPv4 和 IPv6 数据存储在同一个文件中，但确保 ::FFFF: 的路径上没有其他 CIDR 记录，以保持 IPv4 查询路径畅通（修改查询 SDK 可以更好地支持同时查询 IPv4/IPv6）。
