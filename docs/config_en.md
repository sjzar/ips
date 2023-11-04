# IPS Configuration Documentation

<!-- TOC -->
* [IPS Configuration Documentation](#ips-configuration-documentation)
  * [Introduction](#introduction)
  * [Configuration Types](#configuration-types)
  * [Working Directory](#working-directory)
  * [Configuration Loading Logic](#configuration-loading-logic)
  * [Configuration Command](#configuration-command)
    * [Viewing Current Configuration](#viewing-current-configuration)
    * [Modifying Configuration Parameters](#modifying-configuration-parameters)
    * [Removing Configuration Parameters](#removing-configuration-parameters)
    * [Resetting Configuration Parameters](#resetting-configuration-parameters)
  * [Configuration Parameters](#configuration-parameters)
    * [lang](#lang)
    * [ipv4_file](#ipv4file)
    * [ipv4_format](#ipv4format)
    * [ipv6_file](#ipv6file)
    * [ipv6_format](#ipv6format)
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

## Introduction

In order to meet the individual needs of different users and simplify the process of querying, transferring, and packaging IP geolocation databases, IPS provides a flexible configuration system.

Through this system, users can customize the behavior of IPS, including but not limited to specifying the IP database, output format, query fields, etc.

Configuration parameters affect the efficiency of using IPS and the output of results. Reasonable configurations can improve query speed, optimize content output, and make transfer and packaging operations more in line with users' actual needs.

## Configuration Types

The configuration of IPS mainly falls into two categories:

- **Command Line Arguments**: Specified directly in the command line, suitable for one-time tasks or to override default settings.
- **Configuration Files**: Persistently saved in the `ips.json` file in the IPS working directory, suitable for routine operations and personalized settings.

## Working Directory

By default, IPS treats the `.ips` folder in the user's HOME directory as the working directory, with the path `${HOME}/.ips`. This directory is responsible for storing configuration files and IP database files downloaded by IPS.

If the HOME directory is unavailable, IPS will use the system's temporary folder `${TEMP}/.ips` as a fallback.

Users can specify a custom working directory path using the `IPS_DIR` environment variable.

## Configuration Loading Logic

IPS automatically loads the configuration file from the working directory at startup. If the configuration file does not exist or some parameters are missing, IPS will use the built-in default configuration and create or update the configuration file.

It should be noted that command line arguments will take precedence over settings in the configuration file, providing users with a more flexible configuration method.

## Configuration Command

IPS provides the ability to view and modify configuration parameters through the `ips config` command. It aims to offer users a safer and more convenient way to manage the configuration of IPS than directly editing the configuration file.

### Viewing Current Configuration

Running `ips config` without any additional parameters will output the current tool's configuration parameters. This is a quick method to check the current settings and can also be used to confirm whether previous modifications have taken effect.


```shell
# Output current configuration parameters
$ ips config
IPS CONFIG ====================================================================

ips dir:		[/Users/sarv/.ips]
ipv4_file(ipv4):	[ipv4.awdb]
ipv6_file(ipv6):	[GeoLite2-City.mmdb]
fields:			[country,province,city,isp]
text_format:		[%origin [%values] ]
text_values_sep:	[ ]

===============================================================================
... <Omitted output> ...
```

### Modifying Configuration Parameters

The `ips config set` subcommand can be used to modify individual configuration parameters. After a successful modification, IPS will display a confirmation message.

```shell
# Set the IPv4 database file path to the specified file
$ ips config set ipv4 ~/path/to/ipv4.db
INFO[2023-11-02T13:08:01+08:00] set ipv4_file: [/Users/sarv/path/to/ipv4.db] success

# Confirm that the parameter modification was successful
$ ips config
IPS CONFIG ====================================================================

ips dir:		[/Users/sarv/.ips]
ipv4_file(ipv4):	[/Users/sarv/path/to/ipv4.db]  # Modification successful
ipv6_file(ipv6):	[GeoLite2-City.mmdb]
... <Omitted part of the output> ...
```

### Removing Configuration Parameters

The `ips config unset` allows users to remove a configuration parameter, setting its value to empty. Please note that this will not revert to the default value but will require explicit specification or allow IPS to use the default setting in subsequent operations.

```shell
# Remove the configuration of the IPv4 database file path
$ ips config unset ipv4
INFO[2023-11-02T13:09:15+08:00] unset ipv4_file success
```

### Resetting Configuration Parameters

To restore all configuration parameters to their default values, the `ips config reset` subcommand can be used. This will discard all user-defined configurations and restore IPS to its initial state.

```shell
# Reset all configuration parameters to default values
$ ips config reset
INFO[2023-11-02T13:12:21+08:00] reset config success
```

## Configuration Parameters

IPS offers a wide range of configuration parameters to suit different usage scenarios. Below is a detailed introduction to all available configuration parameters, including their functions and default values.

### lang

Sets the language of the output information, a string parameter. The default is `zh-CN` (Chinese), with the following list of supported languages:

| Number | Language   | Code    |
|--------|:-----------|:--------|
| 1      | English    | `en`    |
| 2      | Chinese    | `zh-CN` |
| 3      | Russian    | `ru`    |
| 4      | Japanese   | `ja`    |
| 5      | German     | `de`    |
| 6      | French     | `fr`    |
| 7      | Spanish    | `es`    |
| 8      | Portuguese | `pt-BR` |
| 9      | Persian    | `fa`    |
| 10     | Korean     | `ko`    |

It is important to note that the translation data comes from [GeoNames.org](https://geonames.org). For multilingual IP databases, such as `mmdb`, the translation data in the database will be used directly.

### ipv4_file

Specifies the database file used for IPv4 address queries, a string parameter, with the default value `qqwry.dat`.

### ipv4_format

Specifies the format of the IPv4 database file, a string parameter.

When the database file suffix is insufficient to determine the file format, it is used to specify the format of the IPv4 database. Usually, it is not necessary to set this.

### ipv6_file

Specifies the database file used for IPv6 address queries, a string parameter, with the default value `zxipv6wry.db`.

### ipv6_format

Specifies the format of the IPv6 database file, a string parameter.

When the database file suffix is insufficient to determine the file format, it is used to specify the format of the IPv6 database. Usually, it is not necessary to set this.

### fields

Defines which fields to display in the query results, a string parameter. The default value is `country,province,city,isp`. This parameter supports various advanced usages:

**Generic Field Mapping**

In addition to using database fields directly, `fields` allows you to use generic fields, which means you can refer to the same type of data in different databases that may have different names using one field name.

For example, the generic field `country` represents country information and may be represented as `country`, `country_name`, `countryName`, etc., in different databases.

The current generic fields are:

| Number | Field            | Description                        |
|--------|:-----------------|:-----------------------------------|
| 1      | `country`        | Country                            |
| 2      | `province`       | Province                           |
| 3      | `city`           | City                               |
| 4      | `isp`            | ISP                                |
| 5      | `asn`            | Autonomous System Number           |
| 6      | `continent`      | Continent                          |
| 7      | `utcOffset`      | UTC Offset                         |
| 8      | `latitude`       | Latitude                           |
| 9      | `longitude`      | Longitude                          |
| 10     | `chinaAdminCode` | China Administrative Division Code |

Database fields refer to the unique fields in different databases, such as the `subdivisions` field in `mmdb`, the `region` field in `ip2region`, etc.

**Conditional Field Selection**

`fields` also supports conditional field selection. This allows you to change the output fields based on the values of specific fields in the query results.

Conditions use URL query string format and support basic logical operations, such as `country=中国`, `country=!中国` (not China), `country=中国/美国` (China or USA).

If no condition matches, you can specify a default list of fields. This is done by adding a list of fields without conditions at the end of the conditional selection statement.

```shell
<fields>[|<rule1>:<fields>|<rule2>:<fields>|<default fields>]
 @ <fields> - List of output fields
 @ <rule> - Matching condition
 @ <default fields> - Default list of output fields
 
Example:
  # When the queried IP address is located in China, output the country,province,city,isp fields
  # For IP addresses outside of China, only output the country field
  "country,province,city,isp|country=!中国:country"
```

**Magic Variables**

fields supports the use of specific magic variables to quickly select a set of fields. The currently supported magic variables are:

| Number | Magic Variable   | Actual Represented Fields                        | Description                                     |
|--------|:-----------------|:-------------------------------------------------|:------------------------------------------------|
| 1      | `*`              | -                                                | Wildcard, represents all fields in the database |
| 2      | `find`           | `country,province,city,isp`                      | Query operation                                 |
| 3      | `chinaCity`      | `country,province,city,isp\|country=!中国:country` | Chinese city                                    |
| 4      | `provinceAndISP` | `province,isp\|country=!中国:`                     | Chinese province and ISP                        |
| 5      | `cn`             | `country\|country=中国:country='CN'\|country='OV'` | Chinese country                                 |


### use_db_fields

Decide whether to use the field names that come with the database, a boolean parameter. It is generally used in combination with JSON format output.

When set to `true`, the output will not use common field mapping but will directly display the original field names in the database.

### rewrite_files

This parameter allows you to specify one or multiple files containing rules for rewriting database entries, separated by `,`(commas).

This feature can help correct errors or inaccuracies in the database or format the data into the desired form.

By default, the built-in rewrite files will be automatically loaded from the project's `internal/data` directory.

**Rewrite File Format**

Rewrite files follow a specific format, each line containing a matching condition and a replacement action, separated by a tab (`\t`):

```shell
<condition>\t<replace>\n
 @ <condition> - Uses URL query string format, defining the database fields and values to match.
 @ <replace> - Also uses URL query string format, specifying the new values for the fields.

Example:
  # Rewrite the province field of records with the province information "内蒙古" to "内蒙"
  province=内蒙古	province=内蒙

  # If ASN is 4134, rewrite the ISP field to "电信"
  asn=4134	isp=电信

  # For records with the country field as "辽宁省抚顺市", split it into country "中国", province "辽宁" and city "抚顺" (from qqwry.dat (´･_･`))
  country=辽宁省抚顺市	country=中国&province=辽宁&city=抚顺
```

In these examples, we can see how to use rewrite rules to modify data. This can be used to correct erroneous data, merge or split fields, or convert field values to meet specific output format requirements.

**Using Rewrite Files**

Rewrite files can be specified in the following way:

```shell
$ ips config set rewrite_files "/path/to/rewrite1.txt,/path/to/rewrite2.txt"
```

This command will specify two rewrite files, and IPS will apply the rules in these files when processing database queries.

By using rewrite files, you can ensure that even if the source database contains inaccurate information, the output data will meet the requirements.

### output_type

Specify the format of the command output information, a string field. It can be set to different types according to your needs, so that it can be used in different environments. The default value is `text`. The options are:

- `text`: Outputs the query results in plain text format.
- `json`: Outputs the query results in JSON format.
- `alfred`: Outputs the query results in the format required by Alfred Workflow.

### text_format

When using text output, this parameter defines the specific format of the output, a string parameter.

You can customize this format to display useful information when outputting results. The default value is `%origin [%values]`. The current supported format variables are:

- `%origin`: Represents the original query data.
- `%values`: Represents the field values of the query results.

### text_values_sep

When using text output, this parameter defines the separator between multiple field values, a string parameter. The default value is a ` `(space).

### json_indent

Controls whether to indent JSON format output to improve readability, a boolean parameter. The default value is `false`.

### dp_fields

Similar to the `fields` parameter, this parameter allows you to choose which fields will be included in the final data during the database storage or packaging operation. The default value is empty, which means all fields are included.

### dp_rewriter_files

Similar to the `rewrite_files` parameter, this parameter allows you to specify a list of rewrite files for storage or packaging operations. The default value is empty.

### reader_option

Some database formats provide additional reading options, which can be set during the initialization of the database reader through this parameter to affect the behavior of the reading operation.

For example, `mmdb` database's `disable_extra_data` and so on, please refer to the database documentation for specific functions.

### writer_option

Some database formats provide additional writing options, which can be set during the initialization of the database writer through this parameter to affect the behavior of the writing operation.

For example, `mmdb` database's `select_languages` and so on, please refer to the database documentation for specific functions.

### myip_count

When querying the local IP address, this parameter defines the minimum number of detectors that return the same IP address. The default value is `3`.

### myip_timeout_s

This parameter sets the timeout for each detector when querying the local IP address, in seconds. The default value is `10` seconds.

### addr

When starting the IPS service, this parameter defines the address where the service listens. The default value is `0.0.0.0:6860`, indicating that it listens on port `6860` on all network interfaces.
