# IPS Server Command Documentation

## Introduction

The `ips server` command is used to launch an IPS server that can provide IP address query services.

## Usage

Using the `ips server` command, you can quickly start an IP query server that will listen on a specified address and port. Users can make HTTP requests to query IP addresses and receive responses in JSON format.

## Command Syntax

```shell
ips server [--addr address] [flags]
```

- `-a, --addr string`：Server listening address. Default is `0.0.0.0:6860`, which means listening on port 6860 on all network interfaces.
- `-i, --file string`：Specifies the path to both IPv4 and IPv6 database files.
- `--format string`：Specifies the format for both IPv4 and IPv6 database files; used in conjunction with `--file`. The default is auto-detection.
- `--database-option string`：Specifies options for the database reader. For more information, consult the documentation for the relevant database format or seek professional support.
- `--ipv4-file string`：Specifies the path to the IPv4 database file.
- `--ipv4-format string`：Specifies the format for the IPv4 database file; used in conjunction with `--ipv4-file`. The default is auto-detection.
- `--ipv6-file string`：Specifies the path to the IPv6 database file.
- `--ipv6-format string`：Specifies the format for the IPv6 database file; used in conjunction with `--ipv6-file`. The default is auto-detection.
- `--lang string`：Sets the language for the output. The default is `zh-CN` (Chinese). For more details, refer to [IPS Configuration Documentation](./config_en.md#lang)。
- `-f, --fields string`：Specifies the fields to retrieve from the input file. The default is all fields. For more details, refer to [IPS Configuration Documentation](./config_en.md#fields)。
- `-r, --rewrite-files string`：Specifies a list of files to be rewritten based on the provided configurations. For more details, refer to [IPS Configuration Documentation](./config_en.md#rewritefiles)。

## Examples

### Start an IP Query Server

```shell
# Start the server on local port 8080
ips server -a 127.0.0.1:8080
```

### Use a Custom Database File

```shell
# Start the server using a custom database file
ips server -i GeoLite2-City.mmdb
```

### Set Output Fields and Language

```shell
# Start the server and set output fields and language
ips server -f "country,city" --lang en
```

## API Interface

### Query IP Address

```http request
GET /api/v1/ip?ip=<ip>
Host: <ips host>
Authorization: <none>

200 OK
{
    "ip": <string>,     // IP address
    "net": <string>,    // Subnet of the IP address, in CIDR format
    "data": {}          // Geolocation information
}

400 InvalidArgs
```

### Parse Text and Query Information

```http request
GET /api/v1/query?text=<text>
Host: <ips host>
Authorization: <none>

200 OK
{
    "items": [                  // List of data
        {
            "ip": <string>,     // IP address
            "net": <string>,    // Subnet of the IP address, in CIDR format
            "data": {}          // Geolocation information
        }
    ]
}

400 InvalidArgs
```

## Notes

- The IPS server provides a simple web page at the default entry point (e.g., `http://localhost:6860/` ) for text queries and result display, serving as a demo presentation.
- The IPS server does not currently provide an authentication mechanism, so please avoid exposing the service directly on the public internet.
