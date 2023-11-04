# IPS Download Command Documentation

## Introduction

The `ips download` command is designed to help users simplify the process of acquiring and updating IP geolocation databases.

It is important to note that IPS does not own the copyrights of these databases. All database links are shared by community users or are official free versions.

IPS provides these links for user convenience but is not responsible for the content or copyright of the databases.

If copyright owners believe these links should not be provided within IPS, please contact the author of IPS for prompt removal.

## Usage

`ips download` supports direct downloading of databases through predefined URLs. Users can also provide custom URLs to download the required database files.

After downloading, the new database file path can be configured using the ips config command.

## Command Syntax

```shell
ips download [database_name] [custom_url]
```

- `database_name`: The predefined name of the database.
- `custom_url`: (Optional) A custom download link to download database files if not using a predefined one.

## Predefined Database List

IPS maintains a list of popular IP geolocation databases for download. Below is the list of available databases:

| Database Name       | Format | Download Link                                                                              | Description                |
|:--------------------|:-------|:-------------------------------------------------------------------------------------------|:---------------------------|
| GeoLite2-City.mmdb  | mmdb   | [Link](https://git.io/GeoLite2-City.mmdb)                                                  | MaxMind free edition       |
| city.free.ipdb      | ipdb   | [Link](https://raw.githubusercontent.com/ipipdotnet/ipdb-go/master/city.free.ipdb)         | IPIP.net free edition      |
| dbip-asn-lite.mmdb  | mmdb   | [Link](https://download.db-ip.com/free/dbip-asn-lite-2023-10.mmdb.gz)                      | db-ip free edition         |
| dbip-city-lite.mmdb | mmdb   | [Link](https://download.db-ip.com/free/dbip-city-lite-2023-10.mmdb.gz)                     | db-ip free edition         |
| ip2region.xdb       | xdb    | [Link](https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb) | ip2region free edition     |
| qqwry.dat           | dat    | [Link](https://github.com/metowolf/qqwry.dat/releases/download/20231011/qqwry.dat)         | CZ88.NET database (shared) |
| zxipv6wry.db        | db     | [Link](https://raw.githubusercontent.com/ZX-Inc/zxipdb-python/main/data/ipv6wry.db)        | ip.zxinc.org database      |

These databases are sourced from the Internet, and some are regularly updated. You can access and download the latest versions of the databases via the provided links.

## Examples

### Downloading a Predefined Database

```shell
# Download the free city database provided by IPIP.net
ips download city.free.ipdb
```

Downloading a Database Using a Custom URL and Setting as Default

```shell
# Download a database file using a custom URL
ips download city.ipdb https://foo.com/city.ipdb

# Set as the default database
ips config set ipv4 city.ipdb
```

## Notes

- The download directory is the IPS working directory. For the definition of the working directory, please refer to [IPS Configuration Documentation](./config_en.md#working-directory).
- After downloading a database, it is necessary to specify the database file path in the IPS configuration to use the new database for IP queries.