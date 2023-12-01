# Domain Information Query Library

<!-- TOC -->
* [Domain Information Query Library](#domain-information-query-library)
  * [Introduction](#introduction)
  * [Data Sources](#data-sources)
  * [China Entity List](#china-entity-list)
  * [International Entity List](#international-entity-list)
<!-- TOC -->

## Introduction

The Domain Information Query Library is a Golang code repository embedded with data, enabling the retrieval of related information based on domain names. Currently, it only provides principal information about domains.

Domains are categorized into two types according to the usage scenario: `platform` and `application`.

`platform` refers to domain entities that function as foundational services, providing support to other applications, such as CDN, cloud services, etc.

`application` denotes domain entities that offer services directly to end consumers, such as websites, apps, etc.

In practical use, the focus is not primarily on the type of domain, as a single entity may contain both CDN and website domain names.

## Data Sources

Within China, domain information is sourced from ICP registration records, which presently encompass major CDN providers and select active internet companies.

Internationally, due to the lack of a centralized registration data source, the library currently hosts domain information primarily from major CDN providers, which might not be exhaustive. Plans to extend this dataset further are forthcoming.

## China Entity List

| Type        | Entity      | ICP ID           |
|:------------|:------------|:-----------------|
| application | 4399        | 闽B2-20040099     |
| application | 4399        | 闽ICP备14017048号   |
| application | 爱奇艺         | 京ICP备11032965号   |
| application | 贝壳找房        | 津ICP备18000836号   |
| application | 贝壳找房        | 京ICP备16057509号   |
| application | 哔哩哔哩        | 沪ICP备13002172号   |
| application | 哔哩哔哩        | 沪ICP备13044862号   |
| application | 哔哩哔哩        | 沪ICP备2021021610号 |
| application | 波克城市        | 沪ICP备10215395号   |
| application | 博思软件        | 闽ICP备17021229号   |
| application | 得物          | 沪ICP备16019780号   |
| application | 点点互动        | 京ICP备13035911号   |
| application | 叮咚买菜        | 沪ICP备14046000号   |
| application | 东方财富        | 沪ICP备05006054号   |
| application | 抖音          | 京ICP备12025439号   |
| application | 抖音          | 京ICP备16016397号   |
| application | 抖音          | 京ICP备17044250号   |
| application | 斗鱼          | 鄂ICP备15011961号   |
| application | 斗鱼          | 鄂ICP备2023000041号 |
| application | 钢银          | 沪ICP备15007934号   |
| application | 湖南卫视        | 湘B2-20090004     |
| application | 虎牙          | 粤ICP备16120983号   |
| application | 虎牙          | 粤ICP备17027278号   |
| application | 货拉拉         | 粤ICP备14091972号   |
| application | 竞技世界        | 京ICP备09016390号   |
| application | 巨人网络        | 沪B2-20050107     |
| application | 恺英网络        | 沪ICP备10215773号   |
| application | 看准科技        | 京ICP备14013441号   |
| application | 快手          | 京ICP备15023266号   |
| application | 快手          | 京ICP备16009329号   |
| application | 马上消费        | 渝ICP备15005075号   |
| application | 蚂蚁          | 沪ICP备15027489号   |
| application | 美图          | 闽B2-20040192     |
| application | 美团          | 京ICP备10211739号   |
| application | 米哈游         | 沪ICP备19018275号   |
| application | 米哈游         | 沪ICP备2021018283号 |
| application | 拼多多         | 沪ICP备15010535号   |
| application | 汽车之家        | 京ICP备09113703号   |
| application | 三七互娱        | 沪ICP备11049082号   |
| application | 三七互娱        | 沪ICP备14000728号   |
| application | 盛趣游戏        | 沪ICP备11006268号   |
| application | 盛趣游戏        | 沪ICP备14052292号   |
| application | 世纪华通        | 浙ICP备10211015号   |
| application | 搜狐          | 京ICP证030367号     |
| application | 淘天          | 浙ICP备2023015617号 |
| application | 网龙          | B2-20050038      |
| application | 网易          | 京ICP备10005211号   |
| application | 网易          | 粤B2-20090191     |
| application | 唯品会         | 粤ICP备08114786号   |
| application | 喜马拉雅        | 沪ICP备13027243号   |
| application | 小米          | 京ICP备10046444号   |
| application | 携程          | 沪ICP备08023580号   |
| application | 新浪          | 京ICP备12002058号   |
| application | 新浪          | 京ICP证000007      |
| application | 央视          | 京ICP备06036302号   |
| application | 央视          | 京ICP备10003349号   |
| application | 用友          | 京ICP备05007539号   |
| application | 智联招聘        | 京ICP备17067871号   |
| application | 中原大易        | 豫ICP备16036132号   |
| application | TT语音        | 粤ICP备15000434号   |
| platform    | 阿卡迈         | 京ICP备08100795号   |
| platform    | 阿里          | 浙B2-20080101     |
| platform    | 阿里          | 浙B2-20080224     |
| platform    | 阿里          | 浙B2-20130133     |
| platform    | 阿里          | 浙ICP备09002987号   |
| platform    | 阿里          | 浙ICP备09109183号   |
| platform    | 阿里          | 浙ICP备12022327号   |
| platform    | 白山          | 京ICP备16037635号   |
| platform    | 百度          | 京ICP证030173号     |
| platform    | 北京和顺泰科技     | 京ICP备2021023791号 |
| platform    | 北京数据互通      | 京ICP备16046765号   |
| platform    | 北京兴羽网络      | 京ICP备2023011887号 |
| platform    | 必优          | 京ICP备15065975号   |
| platform    | 博纳云         | 闽ICP备17014081号   |
| platform    | 博睿          | 京ICP备08104257号   |
| platform    | 创世云         | 京ICP备09071955号   |
| platform    | 创世云         | 京ICP备17002507号   |
| platform    | 帝联          | 沪ICP备06052350号   |
| platform    | 甘肃乐天云       | 陇ICP备2022001142号 |
| platform    | 高升          | 沪ICP备14054681号   |
| platform    | 高升          | 吉ICP备11003921号   |
| platform    | 广州烽云        | 粤ICP备18120621号   |
| platform    | 浩云长盛        | 粤ICP备14076948号   |
| platform    | 湖南思极科技      | 湘ICP备2021019142号 |
| platform    | 华为          | 黔ICP备20004760号   |
| platform    | 华为          | 粤A2-20044005号    |
| platform    | 火山          | 京ICP备19059916号   |
| platform    | 火山          | 京ICP备20018813号   |
| platform    | 火山(飞书)      | 京ICP备16045432号   |
| platform    | 江苏云工场       | 苏ICP备16006509号   |
| platform    | 金山          | 京ICP备12032080号   |
| platform    | 京东          | 京ICP备11041704号   |
| platform    | 竞信          | 沪ICP备17043634号   |
| platform    | 浪潮          | 鲁ICP备05019369号   |
| platform    | 联通          | 京ICP备11000964号   |
| platform    | 联通          | 京ICP备16061603号   |
| platform    | 领雾          | 京ICP备2021034657号 |
| platform    | 龙云天下        | 京ICP备17037494号   |
| platform    | 美团          | 京ICP备15052537号   |
| platform    | 魔门云         | 京ICP备15045020号   |
| platform    | 牧云时代        | 京ICP备18028165号   |
| platform    | 南昌首页科技      | 赣ICP备12007401号   |
| platform    | 宁波本电        | 浙ICP备19015321号   |
| platform    | 七牛          | 沪ICP备11037377号   |
| platform    | 奇虎360       | 沪ICP备12043592号   |
| platform    | 奇虎360       | 京ICP备08010314号   |
| platform    | 奇虎360       | 京ICP备18020962号   |
| platform    | 青云          | 京ICP备13019086号   |
| platform    | 锐速云         | 粤ICP备16119720号   |
| platform    | 三全网络        | 苏B2-20070126     |
| platform    | 山迅网络        | 浙ICP备15047040号   |
| platform    | 上海彤旺        | 沪ICP备2021030271号 |
| platform    | 上海云盾        | 沪ICP备11032572号   |
| platform    | 石家庄沐川网络     | 冀ICP备2022002630号 |
| platform    | 世纪互联        | 沪ICP备13015306号   |
| platform    | 世纪互联        | 京ICP备09025197号   |
| platform    | 数掘          | 黔ICP备17001567号   |
| platform    | 思杰系统        | 沪ICP备2022015269号 |
| platform    | 腾讯          | 京ICP备11018762号   |
| platform    | 腾讯          | 鲁ICP备09090609号   |
| platform    | 腾讯          | 粤B2-20090059     |
| platform    | 腾讯          | 粤ICP备14029750号   |
| platform    | 天翼          | 京ICP备12007914号   |
| platform    | 天翼          | 京ICP备2021034386号 |
| platform    | 听云          | 京ICP备08104828号   |
| platform    | 万物云联        | 闽ICP备20008510号   |
| platform    | 网聚云联        | 粤ICP备19040330号   |
| platform    | 网宿          | 沪B2-20030144     |
| platform    | 网宿          | 京ICP备09037435号   |
| platform    | 网宿(同兴万点)    | 京ICP备2022012379号 |
| platform    | 网心          | 粤ICP备14008884号   |
| platform    | 网心          | 粤ICP备2023090978号 |
| platform    | 网易          | 浙ICP备17041593号   |
| platform    | 微软          | 京ICP备09042378号   |
| platform    | 唯一网络        | 粤B1.B2-20070259  |
| platform    | 西部数码        | 蜀ICP备12028237号   |
| platform    | 犀思云         | 京ICP备10009882号   |
| platform    | 新一云         | 沪ICP备19002892号   |
| platform    | 迅达云         | 京ICP备13015005号   |
| platform    | 一九零五(CCTV6) | 京ICP备2023021460号 |
| platform    | 移动          | 京ICP备05002571号   |
| platform    | 亿安天下        | 京ICP备12038793号   |
| platform    | 亿速云         | 粤ICP备17096448号   |
| platform    | 易电通和        | 京ICP备2022021142号 |
| platform    | 优云          | B1-20173072      |
| platform    | 又拍          | 浙ICP备14025602号   |
| platform    | 又拍          | 浙ICP备15031660号   |
| platform    | 云端智度        | 京ICP备16036225号   |
| platform    | 云帆          | 粤ICP备15015607号   |
| platform    | 云瑞智通        | 沪ICP备2022014501号 |
| platform    | 知道创宇        | 京ICP备10040895号   |
| platform    | NiuLink     | 沪ICP备2022005262号 |
| platform    | PPIO        | 沪ICP备18048206号   |
| platform    | UCloud      | 沪ICP备12020087号   |
| platform    | VeryCloud   | 苏ICP备11030873号   |

## International Entity List

| Type     | Entity                    |
|:---------|:--------------------------|
| platform | Adobe Ads                 |
| platform | Adobe EM                  |
| platform | Akamai                    |
| platform | Alibaba                   |
| platform | Amazon Web Services (AWS) |
| platform | CDNetworks                |
| platform | ChinaNetCenter            |
| platform | Cloudflare                |
| platform | Conviva                   |
| platform | DigitalOcean              |
| platform | Edgio                     |
| platform | Fastly                    |
| platform | Google                    |
| platform | Heroku                    |
| platform | IBM                       |
| platform | Imperva                   |
| platform | Leaseweb                  |
| platform | Linode                    |
| platform | Lumen CDN                 |
| platform | Microsoft Azure           |
| platform | Oracle                    |
| platform | OVHcloud                  |
| platform | Rackspace Technology      |
| platform | StackPath                 |
| platform | Statuspage                |
| platform | Verizon                   |
| platform | Zendesk                   |
