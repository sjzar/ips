# IPS Dump Command Documentation

<!-- TOC -->
* [IPS Dump Command Documentation](#ips-dump-command-documentation)
  * [Introduction](#introduction)
  * [Usage](#usage)
  * [Command Syntax](#command-syntax)
  * [Examples](#examples)
    * [Dump IP Database Contents to Standard Output](#dump-ip-database-contents-to-standard-output)
    * [Dump IP Database Contents to a Text File](#dump-ip-database-contents-to-a-text-file)
    * [Customize Export Fields](#customize-export-fields)
    * [Set Output Language and Rewrite Rules](#set-output-language-and-rewrite-rules)
  * [Notes](#notes)
<!-- TOC -->

## Introduction

The `ips dump` command allows users to export data from IP database files to plain text files for data analysis or other processing. This command supports multiple database formats and allows for customization of the output data fields.


## Usage

The `ips dump` command can be used to specify input files, input formats, export fields, and other options to perform data export operations.

## Command Syntax

```shell
ips dump -i inputFile [--input-format] [-o outputFile] [flags]
```

- `-i, --input-file string`：Specifies the path to the input IP database file. Required.
- `--input-format string`：Specifies the format of the input IP database file. Default is auto-detection.
- `--input-option string`：Specifies options for the database reader. For more information, refer to the database documentation.
- `-o, --output-file string`：Specifies the path to the dump file. When not specified, outputs to the standard output stream.
- `--lang string`：Sets the language for the output information. Default is `zh-CN` (Chinese).
- `-f, --fields string`：Specifies the fields to be extracted from the input file. Default is all fields. For a detailed explanation of the parameter, refer to  [IPS Configuration Documentation](./config_en.md#fields)。
- `-r, --rewrite-files string`：Specifies the list of rewrite files to load. For a detailed explanation of the parameter, refer to [IPS Configuration Documentation](./config_en.md#rewritefiles)。

## Examples

### Dump IP Database Contents to Standard Output

```shell
# Export data from the GeoLite2-City.mmdb database file to the standard output stream
ips dump -i GeoLite2-City.mmdb
```

### Dump IP Database Contents to a Text File

```shell
# Export data from the GeoLite2-City.mmdb database file to geoip.txt
ips dump -i GeoLite2-City.mmdb -o geoip.txt
```

### Customize Export Fields

```shell
# Export only the country and city fields from the database file
ips dump -i GeoLite2-City.mmdb -o geoip.txt --fields "country,city"
```

### Set Output Language and Rewrite Rules

```shell
# Set the output data language to English and apply rewrite rules
ips dump -i GeoLite2-City.mmdb -o geoip.txt --lang en -r rewrite_rules.txt
```

## Notes

- Ensure that the `--input-file` points to an existing and valid database file.
- If `--input-format` is specified, make sure it matches the file format.
- Using `--fields` can reduce the amount of data output, exporting only the necessary fields.
- The `--lang` option can set the language of the output data as needed, typically used for multilingual databases.
- `--rewrite-files` can be used to apply custom rewrite rules before exporting data, to correct errors in the database or for data customization.