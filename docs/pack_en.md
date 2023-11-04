# IPS Pack Command Documentation

<!-- TOC -->
* [IPS Pack Command Documentation](#ips-pack-command-documentation)
  * [Introduction](#introduction)
  * [Usage](#usage)
  * [Command Syntax](#command-syntax)
  * [Examples](#examples)
    * [Dump File Packaging IP Database](#dump-file-packaging-ip-database)
    * [Convert IP Database File Format](#convert-ip-database-file-format)
    * [Package IP Database and Specify Fields](#package-ip-database-and-specify-fields)
  * [Notes](#notes)
<!-- TOC -->

## Introduction

The `ips pack` command is used to convert an IP database file into another format. This is particularly useful for users who need to convert the database from one format to another to meet specific application requirements.

## Usage

With the `ips pack` command, users can specify the source database file and its format, and define the target database file and its format, as well as specify the fields to be included.

## Command Syntax

```shell
ips pack -i inputFile [--input-format format] -o outputFile [--output-format format] [flags]
```

- `-i, --input-file string`：Specifies the path to the input IP database file. required.
- `--input-format string`：Specifies the format of the input IP database file. The default is auto-detection.
- `--input-option string`：Specifies options for the database reader. For more information, please consult the relevant database format documentation or obtain professional support.
- `-o, --output-file string`：Specifies the path to the output IP database file. required.
- `--output-format string`：Specifies the format of the output IP database file. If not specified, the format is auto-detected based on the output file extension.
- `--output-option string`：Specifies options for the database writer. For more information, please consult the relevant database format documentation or obtain professional support.
- `--lang string`：Sets the language of the output information. The default is zh-CN (Chinese).
- `-f, --fields string`：Specifies the fields to be extracted from the input file. The default is all fields. For a detailed explanation of the parameters, please refer to [IPS Configuration Documentation](./config_en.md#fields)。
- `-r, --rewrite-files string`：Specifies a list of rewrite files to be loaded. For a detailed explanation of the parameters, please refer to [IPS Configuration Documentation](./config_en.md#rewritefiles)。

## Examples

### Dump File Packaging IP Database

```shell
# Convert dump.txt to ipdb format
ips pack -i dump.txt -o geoip.ipdb
```

### Convert IP Database File Format

```shell
# Convert GeoLite2-City.mmdb database file to ipdb format
ips pack -i GeoLite2-City.mmdb -o geoip.ipdb
```

### Package IP Database and Specify Fields

```shell
# Only export the country and city fields and convert the database to ipdb format
ips pack -i GeoLite2-City.mmdb -o geoip.ipdb --fields "country,city"
```

## Notes

- When specifying `--input-file`, ensure the path to the input file is correct and that the file exists.
- When specifying `--output-file`, ensure the path to the output file is accessible and that you have sufficient permissions to write to it.
- Using `--fields` allows you to customize the data fields included in the output file, reducing unnecessary data storage.
- The `--lang` option allows users to set a specific language for the output data, suitable for databases with multilingual support.
- Custom data rewrite rules can be applied with `--rewrite-files`, which is very useful when adjusting the content of the output file.