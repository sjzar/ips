# IPS Command Documentation

<!-- TOC -->
* [IPS Command Documentation](#ips-command-documentation)
  * [Introduction](#introduction)
  * [Usage](#usage)
  * [Command Syntax](#command-syntax)
  * [Examples](#examples)
    * [Command-Line Query for IP Address](#command-line-query-for-ip-address)
    * [Pipeline Query for IP Address](#pipeline-query-for-ip-address)
    * [Customize Query Fields and Output Format](#customize-query-fields-and-output-format)
  * [Notes](#notes)
<!-- TOC -->

## Introduction

The `ips` command serves not only as the main entry point to the IPS command-line tool, offering a variety of sub-commands for managing IP database operations but also functions as a query command, providing the capability to retrieve geolocation information for IP addresses.

The query command supports both command-line parameter and pipeline methods for queries, suitable for both IPv4 and IPv6 addresses. It can output custom field information based on user configuration, meeting the demand for personalized information display.

## Usage

As a command-line tool, `ips` offers intuitive command-line parameter queries and flexible pipeline queries, enabling users to swiftly obtain geolocation information for IP addresses.

Users can also customize the fields included in the query results, as well as the output format and language, to obtain the display that best suits their needs.

## Command Syntax

```shell
# Command-line parameter query
ips <ip or text> [flags]

# Pipeline query
echo <ip or text> | ips [flags]
```

- `-i, --file string`：Specifies the path to both IPv4 and IPv6 database files.
- `--format string`：Specifies the format for both IPv4 and IPv6 database files; used in conjunction with `--file`. The default is auto-detection.
- `--database-option string`：Specifies options for the database reader. For more information, consult the documentation for the relevant database format or seek professional support.
- `--ipv4-file string`：Specifies the path to the IPv4 database file.
- `--ipv4-format string`：Specifies the format for the IPv4 database file; used in conjunction with `--ipv4-file`. The default is auto-detection.
- `--ipv6-file string`：Specifies the path to the IPv6 database file.
- `--ipv6-format string`：Specifies the format for the IPv6 database file; used in conjunction with `--ipv6-file`. The default is auto-detection.
- `--text-format string`：Specifies the format for text output, supporting `%origin` and `%values` parameters.
- `--text-values-sep string`：Specifies the separator for values in text output, with the default being a space.
- `-j, --json bool`：Outputs results in JSON format.
- `--json-indent bool`：Outputs results in indented JSON format. For more details, refer to [IPS Configuration Documentation](./config_en.md#jsonindent)。
- `--use-db-fields bool`：Uses field names as they appear in the database, typically used with JSON output. For more details, refer to [IPS Configuration Documentation](./config_en.md#usedbfields)。
- `--lang string`：Sets the language for the output. The default is `zh-CN` (Chinese). For more details, refer to [IPS Configuration Documentation](./config_en.md#lang)。
- `-f, --fields string`：Specifies the fields to retrieve from the input file. The default is all fields. For more details, refer to [IPS Configuration Documentation](./config_en.md#fields)。
- `-r, --rewrite-files string`：Specifies a list of files to be rewritten based on the provided configurations. For more details, refer to [IPS Configuration Documentation](./config_en.md#rewritefiles)。
- `--loglevel string`：Sets the logging level, a global parameter with possible values of `trace`, `debug`, `info`, `warn`, `error`, `fatal`, and `panic`, with the default being `info`.

## Examples

### Command-Line Query for IP Address

```shell
# Query geolocation information for an IP address
ips 8.8.8.8

# Query geolocation information for multiple IP addresses
ips 8.8.8.8 119.29.29.29
```

### Pipeline Query for IP Address

```shell
# Use echo to pass the IP address to the ips command
echo 8.8.8.8 | ips

# Use in combination with commands like cat or dig
dig +short google.com | ips
```

### Customize Query Fields and Output Format

```shell
# Query an IP address, outputting only the country and city fields
ips 8.8.8.8 -f "country,city"

# Customize text output format
ips 8.8.8.8 --text-format "%values" --text-values-sep ":" --fields "country,city"
```

## Notes

- If used for the first time without specifying a database file path, the IP database file will be downloaded automatically.
- Regarding the default IP database selection:
  - The default IPv4 database is `qqwry.dat`. The QQWry database is chosen due to its continuous updates from the community (thanks to [@metowolf](https://github.com/metowolf)). For international users, it is recommended to use `GeoLite2-City.mmdb` or a commercial database.
  - The default IPv6 database is `zxipv6wry.db`. Its smaller size can offer a better initial experience when first using the IPS tool, but the database content is relatively outdated (last updated in July 2021). It is recommended to use `GeoLite2-City.mmdb` or a commercial database.
- The IP geolocation information for some countries is updated frequently. If this tool is used for commercial projects, it is essential to replace it with a more recent commercial database!

