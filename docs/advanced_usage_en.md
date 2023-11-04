# IPS Advanced Usage Examples

<!-- TOC -->
* [IPS Advanced Usage Examples](#ips-advanced-usage-examples)
  * [Reducing Fields to Compress Database Size](#reducing-fields-to-compress-database-size)
  * [Creating Custom Databases](#creating-custom-databases)
<!-- TOC -->

## Reducing Fields to Compress Database Size

Commercial databases typically contain a wealth of detailed information such as latitude and longitude, time zones, and postal codes, with high-precision latitude and longitude data particularly taking up a significant amount of space.

However, users do not always require all this data in practice. Trimming unnecessary fields according to business needs can effectively reduce the database size.

For example, the commercial IPv4 database file from [埃文科技](https://www.ipplus360.com) is over 600MB, which is inconvenient for distribution and use.

If you only need country, province, and city information for your use, the following command can be used to repack the database:

```shell
# Repack ipv4.awdb to ipv4_new.awdb, including only country, province, city, and ISP fields
ips pack -i ./ipv4.awdb -f country,province,city,isp -o ./ipv4_new.ipdb
```

The database file size after repacking is about 120MB; further data volume reduction is possible, for example, if city and ISP granularity data is needed in China, but only country-level data is needed overseas, the following command can repack the database accordingly:

```shell
# Include country, province, city, and ISP information for China, and only country information for overseas
ips pack -i ./ipv4.awdb -f 'country,province,city,isp|country=!中国:country' -o ./ipv4_new.ipdb
```

The database size is reduced to approximately 7.7MB after repacking with this command, which is about a 98% reduction from the original file size. This significantly optimizes the efficiency of database distribution and reduces memory usage during operation.

## Creating Custom Databases

IPS allows users to dump database files into a text format, facilitating custom modifications. Subsequently, users can repack the modified text file into a new database file to create a custom database.

For example, to determine if an IP belongs to your company, you can write your company's IP range into a text file in the specified format and then pack it into a database file for custom query purposes.

Since IP database file queries typically use a prefix tree search algorithm, custom database queries usually offer higher efficiency than text file queries and are more convenient to distribute.

```shell
# Dump the database into a text file
ips dump -i ./ipv4.awdb -o ./ipv4.txt

# Pack the custom text file into a database file
ips pack -i ./custom.txt -o ./custom.ipdb
```

Some traffic shaping tools use `mmdb` format database files for routing, and IPS also supports generating `mmdb` format database files.

If your traffic shaping tool only requires `geoname_id` to obtain country information, you can even use the `--output-option` parameter to remove multilingual translation data from the `mmdb` file to further reduce the file size.

```shell
# Pack the custom text file into an mmdb format database file
ips pack -i ./custom.txt -o ./custom.mmdb

# Remove multilingual translation data from the mmdb database file
ips pack -i ./custom.txt -o ./custom.mmdb --output-option "select_languages=-"
```