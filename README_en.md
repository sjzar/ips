## IPS

[![Go Report Card](https://goreportcard.com/badge/github.com/sjzar/ips)](https://goreportcard.com/report/github.com/sjzar/ips)
[![GoDoc](https://godoc.org/github.com/sjzar/ips?status.svg)](https://godoc.org/github.com/sjzar/ips)
[![GitHub release](https://img.shields.io/github/release/sjzar/ips.svg)](https://github.com/sjzar/ips/releases)
[![GitHub license](https://img.shields.io/github/license/sjzar/ips.svg)](https://github.com/sjzar/ips/blob/main/LICENSE)

ips is a command-line tool and library that facilitates the querying, dumping, and packaging of IP geolocation databases.

[中文](./README.md) | English

### Download And Install

#### Installation from Source

```bash
go install github.com/sjzar/ips@latest
```

### Features

* One-click querying, dumping, and packaging of IP geolocation databases
* Compatibility with multiple database formats
* Querying through command-line arguments or piping
* Output in both text and JSON formats
* Customizable query fields with persistent configuration
* Flexible database field rewriting: add or remove fields and modify content as needed

### Supported Databases

| Database     | Query | Dump | Pack | Official Website                                              | Command   |
|:----------|:---|:---|:---|:--------------------------------------------------|:----------|
| txt       | -  | ✅  | ✅  | -                                                 | Used for project dumps  |
| ipdb      | ✅  | ✅  | ✅  | [Link](https://ipip.net)                          |           |
| mmdb      | ✅  | ✅  | ✅  | [Link](https://maxmind.com)                       |           |
| awdb      | ✅  | ✅  | -  | [Link](https://ipplus360.com)                     |           |
| qqwry     | ✅  | ✅  | -  | [Link](https://cz88.net)                          | IPv4 only |
| zxinc     | ✅  | ✅  | -  | [Link](https://ip.zxinc.org)                      | IPv6 only |
| ip2region | ✅  | ✅  | -  | [Link](https://github.com/lionsoul2014/ip2region) | IPv4 only |

### Usage

#### Query

```shell
# Basic query
ips <ip or text> [flags]

# Query IP
ips 61.144.235.160
# Output：61.144.235.160 [中国 广东 深圳 电信]

# Query IP using pipeline
echo "61.144.235.160" | ips
# Output：61.144.235.160 [中国 广东 深圳 电信]

# Query IP using a specific database file
ips -d ./GeoLite2-City.mmdb 61.144.235.160
# Output：61.144.235.160 [中国 广州]

# Query IP using a specific database file and fields
ips -d ./GeoLite2-City.mmdb --fields country 61.144.235.160
# Output：61.144.235.160 [中国]

# Query IP using a specific database file and output in JSON format
ips -d ./GeoLite2-City.mmdb --fields '*' -j 61.144.235.160
# Output：{"ip":"61.144.235.160","net":"61.144.192.0/18","data":{"city":"广州市","continent":"亚洲","country":"中国","latitude":"23.1181","longitude":"113.2539","utcOffset":"Asia/Shanghai"}}
```

#### Dump

```shell
# Basic dump command, output dump content
ips dump -i ./qqwry.dat
# Output：
#    # Dump Time: 2023-10-20 00:00:00
#    # Fields: country,area
#    ... <omitted part of the output> ...

# Specify fields for dumping
ips dump -i ./qqwry.dat -f country
# Output：
#    # Dump Time: 2023-10-20 00:00:00
#    # Fields: country
#    ... <omitted part of the output> ...

# Dump content and save to a file
ips dump -i ./qqwry.dat -o 1.txt
```

#### Pack

```shell
# Package from dump file
ips pack -i qqwry.txt -o qqwry.ipdb

# Package from database file
ips pack -i qqwry.dat -o qqwry.ipdb

# Package from database file specifying fields
ips pack -i qqwry.dat -f country -o country.ipdb
```

### License

`ips` is open-source software licensed under the Apache-2.0 License.

### Acknowledgments

* [IPIP.net](https://ipip.net) for the ipdb database format
* [MaxMind](https://maxmind.com) for the mmdb database format
* [埃文科技](https://ipplus360.com) for the awdb database format
* [纯真网络](https://cz88.net) for the qqwry database format
* [ip.zxinc.org](https://ip.zxinc.org) for the zxinc database format
* [@lionsoul2014](https://github.com/lionsoul2014) for the [ip2region](https://github.com/lionsoul2014/ip2region) database format
* [@zu1k](https://github.com/zu1k) for the [nali](https://github.com/zu1k/nali) project, from which this project's querying feature was inspired
* [@metowolf](https://github.com/metowolf) for the [qqwry.dat](https://github.com/metowolf/qqwry.dat) and ipdb project
* [GeoNames.org](https://geonames.org) for the geolocation data
* Contributors of various Go open-source libraries, such as [cobra](https://github.com/spf13/cobra), [viper](https://github.com/spf13/viper), [logrus](https://github.com/sirupsen/logrus), [progressbar](https://github.com/schollz/progressbar), etc.
