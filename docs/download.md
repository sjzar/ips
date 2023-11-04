# IPS 下载命令说明

<!-- TOC -->
* [IPS 下载命令说明](#ips-下载命令说明)
  * [简介](#简介)
  * [使用方法](#使用方法)
  * [命令语法](#命令语法)
  * [预定义的数据库列表](#预定义的数据库列表)
  * [示例](#示例)
    * [下载预定义数据库](#下载预定义数据库)
    * [使用自定义 URL 下载数据库并设置为默认数据库](#使用自定义-url-下载数据库并设置为默认数据库)
  * [注意事项](#注意事项)
<!-- TOC -->

## 简介

`ips download` 命令用于帮助用户简化 IP 地理位置数据库的获取和更新过程。

需要注意的是，IPS 并不拥有这些数据库的版权，所有的数据库链接均来自于社区用户的分享或是官方提供的免费版本。

IPS 提供这些链接是为了便利用户，但不对数据库内容或版权负责。

如果版权所有者认为不应该在 IPS 中提供这些链接，请联系 IPS 的作者以便及时移除。

## 使用方法

`ips download` 支持直接通过预定义的 URL 下载数据库，用户也可以提供自定义的 URL 来下载所需的数据库文件。

下载完成后，可以通过 `ips config` 命令配置新下载的数据库文件路径。

## 命令语法

```shell
ips download [database_name] [custom_url]
```

- `database_name`：预定义的数据库名称。
- `custom_url`：（可选）自定义的下载链接，如果未使用预定义的文件，则可从中下载数据库文件。

## 预定义的数据库列表

IPS 维护一个包含流行 IP 地理位置数据库的列表，以下是可供下载的数据库：

| 数据库名称               | 格式        | 下载地址                                                                                       | 说明               |
|:--------------------|:----------|:-------------------------------------------------------------------------------------------|:-----------------|
| GeoLite2-City.mmdb  | mmdb      | [Link](https://git.io/GeoLite2-City.mmdb)                                                  | MaxMind 免费版数据库   |
| city.free.ipdb      | ipdb      | [Link](https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb)         | IPIP.net 免费版数据库  |
| dbip-asn-lite.mmdb  | mmdb      | [Link](https://download.db-ip.com/free/dbip-asn-lite-2023-10.mmdb.gz)                      | db-ip 免费版数据库     |
| dbip-city-lite.mmdb | mmdb      | [Link](https://download.db-ip.com/free/dbip-city-lite-2023-10.mmdb.gz)                     | db-ip 免费版数据库     |
| ip2region.xdb       | ip2region | [Link](https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb) | ip2region 免费版数据库 |
| qqwry.dat           | qqwry     | [Link](https://github.com/metowolf/qqwry.dat/releases/download/20231011/qqwry.dat)         | 纯真数据库(社区分享)      |
| zxipv6wry.db        | zxinc     | [Link](https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db)        | ip.zxinc.org 数据库 |

这些数据库均来源于互联网，部分数据库会定期进行更新。您可以通过提供的链接访问和下载最新版本的数据库。

## 示例

### 下载预定义数据库

```shell
# 下载 IPIP.net 提供的免费城市数据库
ips download city.free.ipdb
```

### 使用自定义 URL 下载数据库并设置为默认数据库

```shell
# 通过自定义 URL 下载数据库文件
ips download city.ipdb https://foo.com/city.ipdb

# 设置为默认数据库
ips config set ipv4 city.ipdb
```

## 注意事项

- 下载目录为 IPS 工作目录，关于工作目录的定义可翻阅 [IPS 配置说明](./config.md#工作目录)。
- 下载数据库后，需要在 IPS 的配置中指定数据库文件路径，以便使用新数据库进行 IP 查询。