# IPDB Database Format

<!-- TOC -->
* [IPDB Database Format](#ipdb-database-format)
  * [Introduction](#introduction)
  * [Format Analysis](#format-analysis)
  * [MetaData Analysis](#metadata-analysis)
  * [NodeBlock Analysis](#nodeblock-analysis)
  * [DataBlock Analysis](#datablock-analysis)
  * [Query Operation](#query-operation)
  * [Packaging Process](#packaging-process)
<!-- TOC -->

## Introduction

The IPDB database is an IP database format designed and utilized by [IPIP.net](https://www.ipip.net/), renowned for its compact size, high query efficiency, and multi-language support.

## Format Analysis

```shell
    +--------------------------------+--------------------------------+
    |    MetaData Length (4byte)     |     MetaData (Json Format)     |
    +--------------------------------+--------------------------------+
    |                 Node Chunk (Prefix Tree / Trie)                 |
    +--------------------------------+--------------------------------+
    |                            Data Chunk                           |
    +--------------------------------+--------------------------------+
```

* The file is divided into three parts: metadata, node block, and data block.

## MetaData Analysis

Metadata is a JSON object containing the database build information and metadata required for querying. Below is an example:

```shell
{
    "build": 1632971142,    // Build time
    "ip_version": 1,        // IP database version (IPv4: 0x1, IPv6: 0x2)
    "languages": {
        "CN": 0             // Languages and field offsets
    },
    "node_count": 8705098,  // Number of nodes
    "total_size": 90028407, // Total size of the node block + data block
    "fields": [             // Fields for each data set in the data block
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

## NodeBlock Analysis

The node block consists of a prefix tree (trie). Each node is 8 bytes and stores the offset to the next node. If the offset exceeds the number of nodes, it jumps to the data block, indicating a result has been found.

## DataBlock Analysis

The data block stores IP database data. Identical data is stored only once to reduce redundancy.

```shell
    +--------------------------------+--------------------------------+--------------------------------+
    | Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
    +--------------------------------+--------------------------------+--------------------------------+
    | Data Length (2byte) | Data Fields(<country>\t<province>\t<city>\t<isp>\t<country>\t<province>)   |
    +--------------------------------+--------------------------------+--------------------------------+
```

* In the data block, data is split into two parts: length and data.
* The data part uses `\t` to separate fields. In multi-language versions of the database, different language data is returned using field offsets.

## Query Operation

* A CIDR address is a way of describing a network segment that combines an IP address with a subnet mask, e.g., 10.0.0.1/8 represents an 8-bit subnet mask (255.0.0.0). All IPs within a CIDR segment have identical subnet mask parts. 
* IP addresses are treated as 32-bit/128-bit binary strings. In the node block, the prefix tree is used to search from the beginning, and once a CIDR match is found, it jumps to the data block to return the corresponding data. 
* More CIDR groupings result in a larger number of nodes. 
* CIDRs should not overlap; nested CIDRs (like 10.0.0.1/8 and 10.0.0.1/16) may be precluded from further matching due to earlier matches. 
* For a detailed explanation of the query process, refer to the paper [IPv4 route lookup on Linux](https://vincent.bernat.ch/en/blog/2017-ipv4-route-lookup-linux).

## Packaging Process

* Build the prefix tree and place different datasets into the data block in load order. 
* Prefix trees are built according to IPv6 specifications; IPv4 data needs to fill in 96 bits of subnet data at the front, i.e., 80 bits of 0 and 16 bits of 1, corresponding to the IPv6 mapped address (::FFFF:<IPv4>). This allows fast offset to the 96-bit mask position for IPv4 querying. 
* In theory, the ipdb format database can store both IPv4 and IPv6 data in the same file, but ensure that no other CIDR records exist on the ::FFFF: path to keep the IPv4 query path clear (modifying the query SDK can better support simultaneous IPv4/IPv6 queries).