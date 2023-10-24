## IPS

[![Go Report Card](https://goreportcard.com/badge/github.com/sjzar/ips)](https://goreportcard.com/report/github.com/sjzar/ips)
[![GoDoc](https://godoc.org/github.com/sjzar/ips?status.svg)](https://godoc.org/github.com/sjzar/ips)
[![GitHub release](https://img.shields.io/github/release/sjzar/ips.svg)](https://github.com/sjzar/ips/releases)
[![GitHub license](https://img.shields.io/github/license/sjzar/ips.svg)](https://github.com/sjzar/ips/blob/main/LICENSE)

ips 是一个命令行工具与库，可以轻松完成 IP 地理位置数据库的查询、转存与打包。

中文 | [English](./README_en.md)

### 下载与安装

#### 源码安装

```bash
go install github.com/sjzar/ips@latest
```

### 特性

* 一键查询、转存和打包 IP 地理位置数据库
* 兼容多种数据库格式
* 通过命令行参数或管道进行查询
* 输出支持文本和 JSON 格式
* 可自定义查询字段并持久化配置
* 灵活的数据库字段改写：按需增减字段和内容修改

### 数据库支持列表

| 数据库       | 查询 | 转存 | 打包 | 官方网站                                              | 说明        |
|:----------|:---|:---|:---|:--------------------------------------------------|:----------|
| txt       | -  | ✅  | ✅  | -                                                 | 本项目转存时使用  |
| ipdb      | ✅  | ✅  | ✅  | [Link](https://ipip.net)                          |           |
| mmdb      | ✅  | ✅  | ✅  | [Link](https://maxmind.com)                       |           |
| awdb      | ✅  | ✅  | -  | [Link](https://ipplus360.com)                     |           |
| qqwry     | ✅  | ✅  | -  | [Link](https://cz88.net)                          | IPv4 only |
| zxinc     | ✅  | ✅  | -  | [Link](https://ip.zxinc.org)                      | IPv6 only |
| ip2region | ✅  | ✅  | -  | [Link](https://github.com/lionsoul2014/ip2region) | IPv4 only |

### 使用方法

#### 查询

```shell
# 基础查询
ips <ip或文本> [选项]

# 查询 IP
ips 61.144.235.160
# 输出：61.144.235.160 [中国 广东 深圳 电信]

# 使用管道查询 IP
echo "61.144.235.160" | ips
# 输出：61.144.235.160 [中国 广东 深圳 电信]

# 使用指定的数据库文件查询 IP
ips -d ./GeoLite2-City.mmdb 61.144.235.160
# 输出：61.144.235.160 [中国 广州]

# 使用指定的数据库文件并设置查询字段
ips -d ./GeoLite2-City.mmdb --fields country 61.144.235.160
# 输出：61.144.235.160 [中国]

# 使用指定的数据库文件，以 JSON 格式输出结果
ips -d ./GeoLite2-City.mmdb --fields '*' -j 61.144.235.160
# 输出：{"ip":"61.144.235.160","net":"61.144.192.0/18","data":{"city":"广州市","continent":"亚洲","country":"中国","latitude":"23.1181","longitude":"113.2539","utcOffset":"Asia/Shanghai"}}
```

#### 转存

```shell
# 基础转存命令，输出转存内容
ips dump -i ./qqwry.dat
# 输出：
#    # Dump Time: 2023-10-20 00:00:00
#    # Fields: country,area
#    ... <省略部分输出> ...

# 指定字段进行转存
ips dump -i ./qqwry.dat -f country
# 输出：
#    # Dump Time: 2023-10-20 00:00:00
#    # Fields: country
#    ... <省略部分输出> ...

# 转存内容并保存到文件
ips dump -i ./qqwry.dat -o 1.txt
```

#### 打包

```shell
# 使用转存文件进行打包
ips pack -i qqwry.txt -o qqwry.ipdb

# 使用数据库文件进行打包
ips pack -i qqwry.dat -o qqwry.ipdb

# 使用数据库文件并指定字段进行打包
ips pack -i qqwry.dat -f country -o country.ipdb
```

### 许可

`ips` 是在 Apache-2.0 许可下的开源软件。

### 致谢

* [IPIP.net](https://ipip.net) 的 ipdb 数据库格式
* [MaxMind](https://maxmind.com) 的 mmdb 数据库格式
* [埃文科技](https://ipplus360.com) 的 awdb 数据库格式
* [纯真网络](https://cz88.net) 的 qqwry 数据库格式
* [ip.zxinc.org](https://ip.zxinc.org) 的 zxinc 数据库格式
* [@lionsoul2014](https://github.com/lionsoul2014) 的 [ip2region](https://github.com/lionsoul2014/ip2region) 数据库格式
* [@zu1k](https://github.com/zu1k) 的 [nali](https://github.com/zu1k/nali) 项目，本项目查询功能参考了 nali 的方案
* [@metowolf](https://github.com/metowolf) 的 [qqwry.dat](https://github.com/metowolf/qqwry.dat) 和 ipdb 项目
* [GeoNames.org](https://geonames.org) 的地理信息数据
* 各个 Go 开源库的贡献者们，例如 [cobra](https://github.com/spf13/cobra)、[viper](https://github.com/spf13/viper)、[logrus](https://github.com/sirupsen/logrus)、[progressbar](https://github.com/schollz/progressbar) 等
