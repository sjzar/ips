# IPS 转存命令说明

<!-- TOC -->
* [IPS 转存命令说明](#ips-转存命令说明)
  * [简介](#简介)
  * [使用方法](#使用方法)
  * [命令语法](#命令语法)
  * [示例](#示例)
    * [转存 IP 数据库内容到标准输出流](#转存-ip-数据库内容到标准输出流)
    * [转存 IP 数据库内容到文本文件](#转存-ip-数据库内容到文本文件)
    * [自定义导出字段](#自定义导出字段)
    * [设置输出语言和改写规则](#设置输出语言和改写规则)
  * [注意事项](#注意事项)
<!-- TOC -->>

## 简介

`ips dump` 命令允许用户从 IP 数据库文件中导出数据到文本文件中，用于数据分析或其他处理。此命令支持多种数据库格式，并允许自定义输出的数据字段。

## 使用方法

使用 `ips dump` 命令，可以指定输入文件、输入格式、导出字段等选项来执行数据导出操作。

## 命令语法

```shell
ips dump -i inputFile [--input-format] [-o outputFile] [flags]
```

- `-i, --input-file string`：指定输入 IP 数据库文件的路径。必填项。
- `--input-format string`：指定输入 IP 数据库文件的格式。默认为自动检测。
- `--input-option string`：数据库读取器指定选项。具体信息请查阅数据库文档。
- `-o, --output-file string`：指定转存文件的路径。不指定转存文件时，输出到标准输出流。
- `--lang string`：设置输出信息的语言。默认为 `zh-CN` (中文)。
- `-f, --fields string`：指定从输入文件中获取的字段。默认为所有字段。参数详细解释请参考 [IPS 配置说明](./config.md#fields)。
- `-r, --rewrite-files string`：指定需要载入的改写文件列表。参数详细解释请参考 [IPS 配置说明](./config.md#rewritefiles)。

## 示例

### 转存 IP 数据库内容到标准输出流

```shell
# 从 GeoLite2-City.mmdb 数据库文件导出数据到标准输出流
ips dump -i GeoLite2-City.mmdb
```

### 转存 IP 数据库内容到文本文件

```shell
# 从 GeoLite2-City.mmdb 数据库文件导出数据到 geoip.txt
ips dump -i GeoLite2-City.mmdb -o geoip.txt
```

### 自定义导出字段

```shell
# 从数据库文件中仅导出国家和城市字段
ips dump -i GeoLite2-City.mmdb -o geoip.txt --fields "country,city"
```

### 设置输出语言和改写规则

```shell
# 设置导出数据的语言为英文，并应用改写规则
ips dump -i GeoLite2-City.mmdb -o geoip.txt --lang en -r rewrite_rules.txt
```

## 注意事项

- 确保 `--input-file` 指向的数据库文件是存在且有效的。 
- 如果指定 `--input-format`，请确认格式与文件相符。 
- 使用 `--fields` 可以减少输出的数据量，仅导出需要的字段。 
- `--lang` 选项可以根据需要设置输出数据的语言，通常用于多语言数据库。 
- `--rewrite-files` 可以在导出数据前应用自定义的重写规则，以纠正数据库中的错误或进行数据定制。
