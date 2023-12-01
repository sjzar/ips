# IPS MDNS Command Documentation

## Introduction

The `ips mdns` command is designed to query domain name resolutions across multiple regions.

Utilizing EDNS capabilities, this command sends the client's IP address to the DNS server, which then returns the domain name resolution results for the corresponding region.

The command provides a global perspective on domain name resolutions, enabling users to quickly identify any anomalies in DNS resolutions.

## Usage

Using the `ips mdns` command, users can specify a domain and a DNS server address for querying. The results are displayed in a table format, clearly showing the resolution results from different geographical locations.

## Command Syntax

```shell
ips mdns <domain> [flags]
```

- `--net string`: Specifies the network type for DNS requests, options include `tcp`, `udp`, and `tcp-tls`, with the default being `udp`.
- `--client-timeout int`: Sets the client-side timeout for each DNS request in milliseconds, defaulting to 1000 milliseconds.
- `--single-inflight`: Specifies whether to merge concurrent DNS requests, defaulting to `false`.
- `--timeout int`: Sets the overall timeout for the MDNS command in seconds, with a default of 20 seconds.
- `--exchange-address string`: Specifies the DNS server address, defaulting to `119.29.29.29`.
- `--retry-times`: Determines the number of retry attempts for DNS requests, defaulting to 3 times.

## Examples

### Querying Multi-Region Domain Name Resolutions

```shell
ips mdns i0.hdslb.com
+--------------------+-------------------------------------------------+-------------------------------------------------+
|       GEOISP       |                      CNAME                      |                       IP                        |
+--------------------+-------------------------------------------------+-------------------------------------------------+
| 27.224.0.0         | i0.hdslb.com.04f6a54d.c.cdnhwc1.com [华为]      | 60.165.116.47 [中国 甘肃 兰州 电信]             |
| [中国 甘肃 电信]   | hcdnw.biliv6.c.cdnhwc2.com [华为]               | 60.165.116.48 [中国 甘肃 兰州 电信]             |
+--------------------+-------------------------------------------------+-------------------------------------------------+
| 36.133.72.0        | i0.hdslb.com.w.kunlunno.com [阿里]              | 221.181.64.184 [中国 上海 上海 移动]            |
| [中国 上海 移动]   |                                                 | 221.181.64.148 [中国 上海 上海 移动]            |
+--------------------+-------------------------------------------------+-------------------------------------------------+
| 36.133.108.0       | i0.hdslb.com.04f6a54d.c.cdnhwc1.com [华为]      | 39.136.138.59 [中国 重庆 重庆 移动]             |
| [中国 重庆 移动]   | hcdnw.biliv6.d.cdn.chinamobile.com [移动]       | 39.136.138.58 [中国 重庆 重庆 移动]             |
+--------------------+-------------------------------------------------+-------------------------------------------------+
| 1.56.0.0           | i0.hdslb.com.04f6a54d.c.cdnhwc1.com [华为]      | 218.10.185.43 [中国 黑龙江 鹤岗 联通]           |
| [中国 黑龙江 联通] | hcdnw.biliv6.c.cdnhwc2.com [华为]               | 218.60.101.84 [中国 辽宁 大连 联通]             |
+--------------------+-------------------------------------------------+-------------------------------------------------+
| 42.202.0.0         | i0.hdslb.com.download.ks-cdn.com [金山]         | 123.184.57.130 [中国 辽宁 沈阳 电信]            |
| [中国 辽宁 电信]   | k1-ipv6.gslb.ksyuncdn.com [金山]                | 123.184.57.129 [中国 辽宁 沈阳 电信]            |
+--------------------+-------------------------------------------------+-------------------------------------------------+
|       TOTAL        |                       11                        |                       730                       |
+--------------------+-------------------------------------------------+-------------------------------------------------+
<omitted part of the output>
```

### Using a Custom DNS Server

```shell
ips mdns i0.hdslb.com --exchange-address 8.8.8.8
```

## Important Notes

- If the resolution results are identical across different geographical locations, it is advised to check for local DNS hijacking issues, such as firewalls redirecting traffic on port 53.