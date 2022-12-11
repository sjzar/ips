## IPS

ips is a command-line tool for querying, scanning, and packing IP geolocation database files.

ips is support for multiple database protocols, including popular ones such as MaxMind GeoLite2 and ipip.net, which allows for easy integration with existing systems. In addition, it has flexible query methods that can be chosen based on specific requirements, improving query efficiency.

Using ips, it is easy to query the geolocation information of an IP address, including country, city, and isp. In addition, it can scan the entire database, quickly extract records that meet specific conditions, and package them into a new database file for further processing.

ips is a powerful and easy-to-use tool for IP geolocation, suitable for a wide range of applications that require handling of IP geolocation information.

### Install

```bash
go install github.com/sjzar/ips@latest
```

### Feature
* IP Geolocation Databases Querying, Scanning, Packing
* Multiple Databases Support

### Databases Support Status

| Database | Query | Scan | Pack | Official              | Comments  |
|:---------|:------|:-----|:-----|:----------------------|:----------|
| ipdb     | ✅    | ✅   | ✅    | https://ipip.net      |           |
| awdb     | ✅    | ✅   | -    | https://ipplus360.com |           |
| mmdb     | ✅    | ✅   | -    | https://maxmind.com   |           |
| qqwry    | ✅    | ✅   | -    | https://cz88.net      | IPv4 only |
| zxinc    | ✅    | ✅   | -    | https://ip.zxinc.org  | IPv6 only |


### Usage

#### Query

```shell
ips [option] <ip or text>

# Query IP
ips 61.144.235.160
61.144.235.160 [广东省深圳市 电信]

# Qeury IP with pipeline
echo "61.144.235.160" | ips
61.144.235.160 [广东省深圳市 电信]

# Query IP with database
ips -d ./city.free.ipdb 61.144.235.160
61.144.235.160 [中国 广东 深圳]

# Query IP with database and fields
ips -d ./city.free.ipdb --fields country,province 61.144.235.160
61.144.235.160 [中国 广东]

# Query IP with database and set format
ips -d ./city.free.ipdb.rename --format ipdb 61.144.235.160
61.144.235.160 [中国 广东 深圳]
ips -d ipdb:./city.free.ipdb.rename 61.144.235.160
61.144.235.160 [中国 广东 深圳]
```

#### Scan

```shell
# Scan database
ips scan ./qqwry.dat
# ScanTime: 2006-01-02 15:04:05
# Fields: country,area
# IPVersion: 1
# Meta: {"IPVersion":1,"Fields":["country","area"]}
0.0.0.0/8       IANA,保留地址
1.0.0.0/32      美国,亚太互联网络信息中心(CloudFlare节点)
1.0.0.1/32      美国,APNIC&CloudFlare公共DNS服务器
1.0.0.2/31      美国,亚太互联网络信息中心(CloudFlare节点)
<ignore more content>

# Scan database with fields
ips/ips scan qqwry.dat --fields country
# ScanTime: 2006-01-02 15:04:05
# Fields: country
# IPVersion: 1
# Meta: {"IPVersion":1,"Fields":["country"]}
0.0.0.0/8       IANA
1.0.0.0/24      美国
1.0.1.0/24      福建省
1.0.2.0/23      福建省
1.0.4.0/22      澳大利亚
<ignore more content>

# Scan database with rewrite
# rewrite file format: <field>\t<match>\t<replace>\n
# make rewrite file like:
# country\t美国\t美利坚合众国
ips scan -r ./countrydemo.map qqwry.dat
# ScanTime: 2022-12-10 21:17:01
# Fields: country,area
# IPVersion: 1
# Meta: {"IPVersion":1,"Fields":["country","area"]}
0.0.0.0/8       IANA,保留地址
1.0.0.0/32      美利坚合众国,亚太互联网络信息中心(CloudFlare节点)
1.0.0.1/32      美利坚合众国,APNIC&CloudFlare公共DNS服务器
1.0.0.2/31      美利坚合众国,亚太互联网络信息中心(CloudFlare节点)
<ignore more content>

# Scan database and output to file
ips scan qqwry.dat -o qqwry.ipscan
ll
-rw-r--r--  1 sarv  staff    10M 12 10 20:33 qqwry.dat
-rw-r--r--  1 sarv  staff    48M 12 10 21:19 qqwry.ipscan
```

#### Pack

```shell
# Pack ipscan file
ips pack qqwry.ipscan

# Pack ipscan file and output to another file
ips pack qqwry.ipscan -o demo1.ipdb

# Pack database
ips pack qqwry.dat --format qqwry

# Pack database with fields
ips pack qqwry.dat --format qqwry -f country

# Pack database with rewrite
ips pack qqwry.dat --format qqwry -r ./countrydemo.map
```

### Examples

#### Scan and Rewrite qqwry.dat
```shell
# scan qqwry.dat with fields
ips scan qqwry.dat -f "country,province,city,isp|isp=:country,province,city,area" -o qqwry.ipscan

# scan qqwry.dat with rewrite
ips scan qqwry.dat -f "country,province,city,isp|isp=:country,province,city,area" -r ./data/qqwry_area.map,./data/qqwry_country.map -o qqwry_rewrite.ipscan

# diff qqwry.ipscan and qqwry_rewrite.ipscan
# qqwry.ipscan
1.24.24.0/21	内蒙古鄂尔多斯市,,,联通
108.162.243.0/24	美国华盛顿州西雅图,,,CloudFlare节点
219.138.4.0/26	长江大学,,, CZ88.NET
# qqwry_rewrite.ipscan
1.24.24.0/21	中国,内蒙古,鄂尔多斯,联通
108.162.243.0/24	美国,华盛顿,西雅图,CloudFlare
219.138.4.0/26	中国,湖北,荆州/长江大学, CZ88.NET
```