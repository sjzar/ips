# IPS 服务命令说明

<!-- TOC -->
* [IPS 服务命令说明](#ips-服务命令说明)
  * [简介](#简介)
  * [使用方法](#使用方法)
  * [命令语法](#命令语法)
  * [示例](#示例)
    * [启动 IP 查询服务](#启动-ip-查询服务)
    * [使用自定义数据库文件](#使用自定义数据库文件)
    * [设置输出字段和语言](#设置输出字段和语言)
  * [API 接口](#api-接口)
    * [查询 IP 地址](#查询-ip-地址)
    * [解析文本并查询信息](#解析文本并查询信息)
  * [注意事项](#注意事项)
<!-- TOC -->

## 简介

`ips server` 命令用于启动一个 IPS 服务，该服务能够提供 IP 地址查询服务。

## 使用方法

使用 `ips server` 命令可以快速启动一个 IP 查询服务，它将监听指定的地址和端口。用户可以通过 HTTP 请求来查询 IP 地址，得到 JSON 格式的响应。

## 命令语法

```shell
ips server [--addr address] [flags]
```

- `-a, --addr string`：服务监听地址。默认值为 `0.0.0.0:6860`，表示在所有网络接口的 `6860` 端口上监听。
- `-i, --file string`：同时指定 IPv4 和 IPv6 数据库文件的路径。
- `--format string`：指定 IPv4 和 IPv6 数据库文件的格式，需要与 `--file` 配合使用。默认为自动检测。
- `--database-option string`：数据库读取器指定选项。具体信息请查阅相关的数据库格式文档或获取专业支持。
- `--ipv4-file string`：指定 IPv4 数据库文件的路径。
- `--ipv4-format string`：指定 IPv4 数据库文件的格式，需要与 `--ipv4-file` 配合使用。默认为自动检测。
- `--ipv6-file string`：指定 IPv6 数据库文件的路径。
- `--ipv6-format string`：指定 IPv6 数据库文件的格式，需要与 `--ipv6-file` 配合使用。默认为自动检测。
- `--hybrid-mode string`: 指定混合读取器的操作模式，可选值为 `comparison` 与 `aggregation`，参数详细解释请参考 [IPS 配置说明](./config.md#hybridmode)。
- `--lang string`：设置输出信息的语言。默认为 `zh-CN` (中文)。参数详细解释请参考 [IPS 配置说明](./config.md#lang)。
- `-f, --fields string`：指定从输入文件中获取的字段。默认为所有字段。参数详细解释请参考 [IPS 配置说明](./config.md#fields)。
- `-r, --rewrite-files string`：指定需要载入的改写文件列表。参数详细解释请参考 [IPS 配置说明](./config.md#rewritefiles)。

## 示例

### 启动 IP 查询服务

```shell
# 在本地 8080 端口启动服务
ips server -a 127.0.0.1:8080
```

### 使用自定义数据库文件

```shell
# 使用自定义数据库文件启动服务
ips server -i GeoLite2-City.mmdb
```

### 设置输出字段和语言

```shell
# 启动服务，并设置输出字段和语言
ips server -f "country,city" --lang en
```

## API 接口

### 查询 IP 地址

```http request
GET /api/v1/ip?ip=<ip>
Host: <ips host>
Authorization: <none>

200 OK
{
    "ip": <string>,     // IP 地址
    "net": <string>,    // IP 地址所在子网，CIDR 格式
    "data": {}          // 地理位置信息
}

400 InvalidArgs
```

### 解析文本并查询信息

```http request
GET /api/v1/query?text=<text>
Host: <ips host>
Authorization: <none>

200 OK
{
    "items": [                  // 数据列表
        {
            "ip": <string>,     // IP 地址
            "net": <string>,    // IP 地址所在子网，CIDR 格式
            "data": {}          // 地理位置信息
        }
    ]
}

400 InvalidArgs
```

## 注意事项

- IPS 服务在默认入口(例如 `http://localhost:6860/` )提供了一个简单的 Web 页面，提供文本查询和结果展示功能，用作 Demo 演示。
- IPS 服务暂未提供鉴权机制，请避免直接将服务暴露在公网环境下运行。