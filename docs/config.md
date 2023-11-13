# IPS 配置说明

<!-- TOC -->
* [IPS 配置说明](#ips-配置说明)
  * [简介](#简介)
  * [配置类型](#配置类型)
  * [工作目录](#工作目录)
  * [配置加载逻辑](#配置加载逻辑)
  * [配置命令](#配置命令)
    * [查看当前配置](#查看当前配置)
    * [修改配置参数](#修改配置参数)
    * [移除配置参数](#移除配置参数)
    * [重置配置参数](#重置配置参数)
  * [配置参数](#配置参数)
    * [lang](#lang)
    * [ipv4_file](#ipv4file)
    * [ipv4_format](#ipv4format)
    * [ipv6_file](#ipv6file)
    * [ipv6_format](#ipv6format)
    * [hybrid_mode](#hybridmode)
    * [fields](#fields)
    * [use_db_fields](#usedbfields)
    * [rewrite_files](#rewritefiles)
    * [output_type](#outputtype)
    * [text_format](#textformat)
    * [text_values_sep](#textvaluessep)
    * [json_indent](#jsonindent)
    * [dp_fields](#dpfields)
    * [dp_rewriter_files](#dprewriterfiles)
    * [reader_option](#readeroption)
    * [writer_option](#writeroption)
    * [myip_count](#myipcount)
    * [myip_timeout_s](#myiptimeouts)
    * [addr](#addr)
<!-- TOC -->

## 简介

IPS 为了满足不同用户的个性化需求，简化IP地理位置数据库的查询、转存和打包过程，提供了一套灵活的配置系统。

通过这套系统，用户可以定制IPS的行为，包括但不限于指定 IP 数据库、输出格式、查询字段等。

配置参数影响着 IPS 的使用效率和结果输出。合理的配置可以提升查询速度，优化输出内容，并且使得转存与打包操作更加符合用户的实际需求。

## 配置类型

IPS 的配置主要分为两大类：

- **命令行参数**：直接在命令行中指定，适用于一次性任务或覆盖默认配置。
- **配置文件**：在 IPS 工作目录下的 `ips.json` 文件中持久化保存，适用于常规操作和个性化设置。

## 工作目录

IPS 默认将用户的 HOME 目录下的 `.ips` 文件夹作为工作目录，路径为 `${HOME}/.ips`。 该目录负责存储配置文件和 IPS 下载的 IP 数据库文件。

如果 HOME 目录不可用，IPS 将使用系统的临时文件夹 `${TEMP}/.ips` 作为后备。

用户可以通过环境变量 `IPS_DIR` 来指定一个自定义的工作目录路径。

## 配置加载逻辑

IPS 在启动时会自动加载工作目录中的配置文件。如果配置文件不存在或部分参数缺失，IPS将采用内置的默认配置，并创建或更新配置文件。

需要注意的是，命令行参数会优先于配置文件中的设置，为用户提供更灵活的配置方式。

## 配置命令

IPS 通过 `ips config` 命令提供了查看和修改配置参数的能力。希望为用户提供一种比直接编辑配置文件更安全、更便捷的方式来管理IPS的配置。

### 查看当前配置

运行 `ips config` 不带任何额外参数将输出当前工具的配置参数。这是检查当前设置的快速方法，也可以用来确认之前的修改是否已经生效。

```shell
# 输出当前的配置参数
$ ips config
IPS CONFIG ====================================================================

ips dir:		[/Users/sarv/.ips]
ipv4_file(ipv4):	[ipv4.awdb]
ipv6_file(ipv6):	[GeoLite2-City.mmdb]
fields:			[country,province,city,isp]
text_format:		[%origin [%values] ]
text_values_sep:	[ ]

===============================================================================
... <省略部分输出> ...
```

### 修改配置参数

使用 `ips config set` 子命令可以修改单个配置参数。修改成功后，IPS将显示一个确认信息。

```shell
# 将IPv4数据库文件路径设置为指定的文件
$ ips config set ipv4 ~/path/to/ipv4.db
INFO[2023-11-02T13:08:01+08:00] set ipv4_file: [/Users/sarv/path/to/ipv4.db] success

# 确认参数修改成功
$ ips config
IPS CONFIG ====================================================================

ips dir:		[/Users/sarv/.ips]
ipv4_file(ipv4):	[/Users/sarv/path/to/ipv4.db]  # 修改成功
ipv6_file(ipv6):	[GeoLite2-City.mmdb]
... <省略部分输出> ...
```

### 移除配置参数

使用 `ips config unset` 允许用户移除一个配置参数，将其值置空。请注意，这并不会恢复到默认值，而是在后续操作中需要显式指定或让IPS使用默认设置。

```shell
# 将IPv4数据库文件路径的配置移除
$ ips config unset ipv4
INFO[2023-11-02T13:09:15+08:00] unset ipv4_file success
```

### 重置配置参数

如果需要将所有的配置参数恢复到默认值，可以使用 `ips config reset` 子命令。这将丢弃所有用户自定义的配置，将IPS恢复到初始状态。

```shell
# 重置所有配置参数到默认值
$ ips config reset
INFO[2023-11-02T13:12:21+08:00] reset config success
```

## 配置参数

IPS 提供了丰富的配置参数以适应不同的使用场景。以下是所有可用配置参数的详细介绍，包括它们的功能与默认值等信息。

### lang

设置输出信息的语言，字符串参数。默认为 `zh-CN` (中文)，支持的语言列表如下：

| 序号 | 语言   | 代码      |
|----|:-----|:--------|
| 1  | 英语   | `en`    |
| 2  | 中文   | `zh-CN` |
| 3  | 俄语   | `ru`    |
| 4  | 日语   | `ja`    |
| 5  | 德语   | `de`    |
| 6  | 法语   | `fr`    |
| 7  | 西班牙语 | `es`    |
| 8  | 葡萄牙语 | `pt-BR` |
| 9  | 波斯语  | `fa`    |
| 10 | 韩语   | `ko`    |

需要注意的是，翻译数据源自 [GeoNames.org](https://geonames.org)。对于支持多语言的 IP 数据库，比如 `mmdb`，将直接使用数据库中的翻译数据。

### ipv4_file

指定 IPv4 地址查询时使用的数据库文件，字符串参数，默认值为 `qqwry.dat`。

### ipv4_format

指定 IPv4 数据库文件的格式，字符串参数。

当数据库文件后缀不足以确定文件格式时，用来指定 IPv4 数据库的格式。通常不需要设置。

### ipv6_file

指定 IPv6 地址查询时使用的数据库文件，字符串参数，默认值为 `zxipv6wry.db`。

### ipv6_format

指定 IPv6 数据库文件的格式，字符串参数。

当数据库文件后缀不足以确定文件格式时，用来指定 IPv6 数据库的格式。通常不需要设置。

### hybrid_mode

指定了混合读取器（Hybrid Reader）的操作模式，字符串参数。操作模式决定了如何处理和组合来自多个 IP 数据库的数据。

可选参数为 `comparison` 与 `aggregation`，默认值为 `aggregation`。

- `comparison`：比较模式，适用于需要跨不同 IP 数据库比较数据的场景，输出所有集成数据库的数据，便于识别每个源之间的差异和变化。
- `aggregation`：聚合模式，适用于需要统一、全面视图的 IP 信息的情况，从多个源聚合数据，用一个数据库中的信息补充另一个数据库中缺失的字段。

### fields

定义查询结果中显示哪些字段，字符串参数。默认值是 `country,province,city,isp`。该参数支持多种高级用法：

**通用字段映射**

除了直接使用数据库字段，`fields` 允许您使用通用字段，这意味着您可以使用一个字段名来引用不同数据库中可能有不同名称的相同类型数据。 

例如通用字段 `country` 表示国家信息，可能在不同的数据库中表示为 `country`、`country_name`、`countryName` 等。

当前的通用字段有：

| 序号 | 字段     | 说明     |
|----|:-------|:-------|
| 1  | `country` | 国家     |
| 2  | `province` | 省份     |
| 3  | `city`    | 城市     |
| 4  | `isp`     | 运营商    |
| 5  | `asn`     | 自治域号   |
| 6  | `continent` | 大洲     |
| 7  | `utcOffset` | UTC 偏移值 |
| 8  | `latitude` | 纬度     |
| 9  | `longitude` | 经度     |
| 10 | `chinaAdminCode` | 中国行政区划代码 |

数据库字段指的是不同数据库中特有的字段，例如 `mmdb` 中的 `subdivisions` 字段，`ip2region` 中的 `region` 字段等。

**条件字段选择**

`fields` 还支持基于条件的字段选择。这允许您根据查询结果中的特定字段值来改变输出字段。

条件使用 URL 查询字符串格式，支持基本的逻辑运算，如 `country=中国`、`country=!中国`（非中国）、`country=中国/美国`（中国或美国）。

如果没有条件匹配，可以指定一个默认字段列表。这通过在条件选择语句的最后添加一个没有条件的字段列表来实现。

```shell
<fields>[|<rule1>:<fields>|<rule2>:<fields>|<default fields>]
 @ <fields> - 输出字段列表
 @ <rule> - 匹配条件
 @ <default fields> - 默认字段列表
 
举例：
  # 当查询的 IP 地址位于中国时，输出 country,province,city,isp 字段
  # 对于非中国的 IP 地址，只输出 country 字段
  "country,province,city,isp|country=!中国:country"
```

**魔法变量**

`fields` 支持使用特定的魔法变量快速选择一组字段。当前支持的魔法变量有：

| 序号 | 魔法变量             | 实际表示字段                                           | 说明               |
|----|:-----------------|:-------------------------------------------------|:-----------------|
| 1  | `*`              | -                                                | 通配符，表示输出数据库中所有字段 |
| 2  | `find`           | `country,province,city,isp`                      | 查询操作             |
| 3  | `chinaCity`      | `country,province,city,isp\|country=!中国:country` | 中国城市             |
| 4  | `provinceAndISP` | `province,isp\|country=!中国:`                     | 中国省份和运营商         |
| 5  | `cn`             | `country\|country=中国:country='CN'\|country='OV'` | 中国国家             |

### use_db_fields

确定是否使用数据库自带的字段名称，布尔值参数。一般与 JSON 格式输出配合使用。

当设置为 `true` 时，输出将不使用通用字段映射，而是直接显示数据库中的原始字段名。

### rewrite_files

此参数允许您指定一个或多个文件，这些文件包含了用于改写数据库条目的规则，文件之间使用 `,` 分隔。

此功能可以帮助修正数据库中的错误或不精确的数据，或者将数据格式化为所需的形式。

默认情况下，将从项目的 `internal/data` 目录自动载入内置的改写文件。

**改写文件格式**

改写文件遵循特定的格式，每一行包含一个匹配条件和一个替换动作，由制表符 (`\t`) 分隔：

```shell
<condition>\t<replace>\n
 @ <condition> - 使用 URL 查询字符串格式，定义需要匹配的数据库字段和值。
 @ <replace> - 同样使用 URL 查询字符串格式，指定字段的新值。

举例：
  # 将省份信息为 "内蒙古" 的记录的省份字段改写为 "内蒙"
  province=内蒙古	province=内蒙

  # 如果 ASN 为 4134，将 ISP 字段改写为 "电信"
  asn=4134	isp=电信

  # 对于国家字段是 "辽宁省抚顺市" 的记录，将其拆分为国家 "中国"、省份 "辽宁" 和城市 "抚顺"（来自 qqwry.dat (´･_･`)）
  country=辽宁省抚顺市	country=中国&province=辽宁&city=抚顺
```

在这些例子中，我们可以看到如何使用改写规则来修改数据。这可以用于纠正错误的数据，合并或拆分字段，或者转换字段值以适应特定的输出格式需求。

**使用改写文件**

可以通过以下方式指定改写文件：

```shell
$ ips config set rewrite_files "/path/to/rewrite1.txt,/path/to/rewrite2.txt"
```

此命令将指定两个改写文件，ips 将在处理数据库查询时应用这些文件中的规则。

通过使用改写文件，您可以确保即使源数据库包含了不精确的信息，输出数据也能符合要求。

### output_type

指定命令输出信息的格式，字符串字段。它可以根据您的需要设置为不同的类型，以便在不同的环境下使用。默认值为 `text`。 可选值有：

- `text`: 以纯文本格式输出查询结果。
- `json`: 以 JSON 格式输出查询结果。
- `alfred`: 以 Alfred Workflow 所需的格式输出查询结果。

### text_format

当您使用文本输出时，此参数定义了输出的具体格式，字符串参数。

您可以自定义这个格式，以便在输出结果时展示有用的信息。默认值为 `%origin [%values]`。 当前支持的格式变量有：

- `%origin`: 表示原始的查询数据。
- `%values`: 表示查询结果的字段值。

### text_values_sep

当您使用文本输出时，此参数定义了多个字段值之间的分隔符，字符串参数。默认值为 ` `(空格)。

### json_indent

控制 JSON 格式输出时是否进行缩进，以提高可读性，布尔值参数。默认值为 `false`。

### dp_fields

功能与 `fields` 字段类似，在进行数据库的转存或打包操作时，此参数允许您选择哪些字段将被包含在最终的数据中。默认值为空，代表包含所有字段。

### dp_rewriter_files

功能与 `rewrite_files` 字段类似，此参数允许您指定用于转存或打包操作的改写文件列表。默认值为空。

### reader_option

一些数据库格式提供了额外的读取选项，通过此参数可以在初始化数据库读取器时进行设置，用以影响读取操作的行为。

例如 `mmdb` 数据库的 `disable_extra_data` 等，具体功能请查阅数据库文档。

### writer_option

一些数据库格式提供了额外的写入选项，通过此参数可以在初始化数据库写入器时进行设置，用以影响写入操作的行为。

例如 `mmdb` 数据库的 `select_languages` 等，具体功能请查阅数据库文档。

### myip_count

在查询本机 IP 地址时，此参数定义了返回相同 IP 地址的最小探测器数量。默认值为 `3`。

### myip_timeout_s

此参数设置查询本机 IP 地址时，每个探测器的超时时间（以秒为单位）。默认值为 `10` 秒。

### addr

在启动 IPS 服务时，此参数定义了服务监听的地址。默认值为 `0.0.0.0:6860`，表示在所有网络接口的 `6860` 端口上监听。

