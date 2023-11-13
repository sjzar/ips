# IPS 打包命令说明

<!-- TOC -->
* [IPS 打包命令说明](#ips-打包命令说明)
  * [简介](#简介)
  * [使用方法](#使用方法)
  * [命令语法](#命令语法)
  * [示例](#示例)
    * [转存文件打包 IP 数据库](#转存文件打包-ip-数据库)
    * [转换 IP 数据库文件格式](#转换-ip-数据库文件格式)
    * [打包 IP 数据库并指定字段](#打包-ip-数据库并指定字段)
  * [注意事项](#注意事项)
<!-- TOC -->

## 简介

`ips pack` 命令用于将 IP 数据库文件转换成另一种格式的数据库文件。这对于需要将数据库从一种格式转换为另一种格式以满足特定应用需求的用户尤其有用。

## 使用方法

通过 `ips pack` 命令，用户可以指定源数据库文件及其格式，并定义目标数据库文件及其格式，还可以指定要包含在内的字段。

## 命令语法

```shell
ips pack -i inputFile [--input-format format] -o outputFile [--output-format format] [flags]
```

- `-i, --input-file string`：指定输入 IP 数据库文件的路径。必填项。
- `--input-format string`：指定输入 IP 数据库文件的格式。默认为自动检测。
- `--input-option string`：数据库读取器指定选项。具体信息请查阅相关的数据库格式文档或获取专业支持。
- `--hybrid-mode string`: 指定混合读取器的操作模式，可选值为 `comparison` 与 `aggregation`，参数详细解释请参考 [IPS 配置说明](./config.md#hybridmode)。
- `-o, --output-file string`：指定输出 IP 数据库文件的路径。必填项。
- `--output-format string`：指定输出 IP 数据库文件的格式。未指定时，使用输出文件的扩展名自动检测。
- `--output-option string`：数据库写入器指定选项。具体信息请查阅相关的数据库格式文档或获取专业支持。
- `--lang string`：设置输出信息的语言。默认为 `zh-CN` (中文)。
- `-f, --fields string`：指定从输入文件中获取的字段。默认为所有字段。参数详细解释请参考 [IPS 配置说明](./config.md#fields)。
- `-r, --rewrite-files string`：指定需要载入的改写文件列表。参数详细解释请参考 [IPS 配置说明](./config.md#rewritefiles)。

## 示例

### 转存文件打包 IP 数据库

```shell
# 将 dump.txt 转换为 ipdb 格式
ips pack -i dump.txt -o geoip.ipdb
```

### 转换 IP 数据库文件格式

```shell
# 将 GeoLite2-City.mmdb 数据库文件转换为 ipdb 格式
ips pack -i GeoLite2-City.mmdb -o geoip.ipdb
```

### 打包 IP 数据库并指定字段

```shell
# 仅导出国家和城市字段，并将数据库转换为 ipdb 格式
ips pack -i GeoLite2-City.mmdb -o geoip.ipdb --fields "country,city"
```

## 注意事项
- 在指定 `--input-file` 时，确保输入文件的路径正确，并且该文件存在。
- 在指定 `--output-file` 时，确保输出文件的路径可访问，并且有足够的权限进行写入操作。
- 使用 `--fields` 可以自定义输出文件中包含的数据字段，减少不必要的数据存储。
- `--lang` 选项允许用户为输出数据设置特定的语言，适用于多语言支持的数据库。
- 通过 `--rewrite-files` 可以应用自定义的数据改写规则，这在调整输出文件的数据内容时非常有用。