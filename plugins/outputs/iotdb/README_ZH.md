[English](./README.md) | [中文](./README_ZH.md)

# Apache IoTDB

[![Main Mac and Linux](https://github.com/apache/iotdb/actions/workflows/main-unix.yml/badge.svg)](https://github.com/apache/iotdb/actions/workflows/main-unix.yml)
[![Main Win](https://github.com/apache/iotdb/actions/workflows/main-win.yml/badge.svg)](https://github.com/apache/iotdb/actions/workflows/main-win.yml)
[![coveralls](https://coveralls.io/repos/github/apache/iotdb/badge.svg?branch=master)](https://coveralls.io/repos/github/apache/iotdb/badge.svg?branch=master)
[![GitHub release](https://img.shields.io/github/release/apache/iotdb.svg)](https://github.com/apache/iotdb/releases)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
![](https://github-size-badge.herokuapp.com/apache/iotdb.svg)
![](https://img.shields.io/github/downloads/apache/iotdb/total.svg)
![](https://img.shields.io/badge/platform-win10%20%7C%20macox%20%7C%20linux-yellow.svg)
![](https://img.shields.io/badge/java--language-1.8-blue.svg)
[![Language grade: Java](https://img.shields.io/lgtm/grade/java/g/apache/iotdb.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/apache/iotdb/context:java)
[![IoTDB Website](https://img.shields.io/website-up-down-green-red/https/shields.io.svg?label=iotdb-website)](https://iotdb.apache.org/)
[![Maven Version](https://maven-badges.herokuapp.com/maven-central/org.apache.iotdb/iotdb-parent/badge.svg)](http://search.maven.org/#search|gav|1|g:"org.apache.iotdb")
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://join.slack.com/t/apacheiotdb/shared_invite/zt-qvso1nj8-7715TpySZtZqmyG5qXQwpg)

Apache IoTDB（物联网数据库）是一个物联网原生数据库，在数据管理和分析方面表现良好，可部署在边缘设备和云上。
由于其轻量级架构、高性能和丰富的功能集，以及与Apache Hadoop、Spark和Flink的深度集成，
Apache IoTDB可以满足物联网工业领域的海量数据存储、高速数据摄取和复杂数据分析的要求。

Apache IoTDB website: https://iotdb.apache.org
Apache IoTDB Github: https://github.com/apache/iotdb

# Apache IoTDB Go语言客户端

[![E2E Tests](https://github.com/apache/iotdb-client-go/actions/workflows/e2e.yml/badge.svg)](https://github.com/apache/iotdb-client-go/actions/workflows/e2e.yml)
[![GitHub release](https://img.shields.io/github/release/apache/iotdb-client-go.svg)](https://github.com/apache/iotdb-client-go/releases)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
![](https://github-size-badge.herokuapp.com/apache/iotdb-client-go.svg)
![](https://img.shields.io/badge/platform-win10%20%7C%20macos%20%7C%20linux-yellow.svg)
[![IoTDB Website](https://img.shields.io/website-up-down-green-red/https/shields.io.svg?label=iotdb-website)](https://iotdb.apache.org/)

Apache IoTDB 有一个go语言客户端，能够使用go语言原生接口支持 IoTDB 的数据增删改查。

Apache IoTDB Golang Client Github: https://github.com/apache/iotdb

# IoTDB 输出插件

IoTDB 输出插件可以把 Telegraf 采集到的数据保存到IoTDB数据库。该插件使用了go语言客户端的接口，能够支持会话连接、数据插入。

## 快速上手

使用该插件前，需要配置数据库服务器的ip地址、所使用的端口号、用户名、密码等信息，以及一些数据类型转换、时间单位等配置。

英文的配置文件请参考：[English Configuration](./sample.conf)，中文配置文件请参考[中文配置样例](./sample_zh.conf). 或者，对应版本的配置内容也在后文中列出。

## 注意事项

1. IoTDB 0.13.x版本以及之前的版本，**不支持无符号整数**。所以本插件提供了三种可选的无符号整数处理方式，只需要指定参数`convertUint64To`的取值即可。该参数有三个取值，分别对应不同的处理放肆，分别是：
   - `ToInt64`，默认的处理方式。对于未超出`int64`表示范围的无符号整数，以`int64`类型存储；如果超出表示范围，则保存`math.MaxInt64`，也即`9223372036854775807`。
   - `ForceToInt64`，强制类型转换为`int64`。如果数字超过`int64`表示范围，可能会抛出异常。
   - `Text`，强制转换为字符串。无论无符号整型多大，都会被转换为字符串保存，不丢失精度。

2. IoTDB支持多种时间精度，但无论何种精度，都以`int64`类型存储，所以用户需要指定时间戳的语义。用户需要指定参数`convertUint64To`的取值，该参数默认取值为`nanosecond`。

3. IoTDB目前不能很好地支持标签（Tag）索引，目前采用的处理方式请参考[InfluxDB-Protocol适配器](https://iotdb.apache.org/zh/UserGuide/Master/API/InfluxDB-Protocol.html)。用户需要指定参数`treateTagsAs`的取值，来决定如何处理标签：

   - `Measurements`，Tag会被看做一个普通的物理量，等同于Field。只不过Tag的取值总是字符串。
   - `DeviceID_subtree`，Tag会被看做设备标识路径（device id）的一部分。Tags的顺序是有序的，该顺序由Telegraf决定，一般为字典序升序排列。

   举例：当一个metric的取值为，`Name="root.sg.device", Tags={tag1="private", tag2="working"}, Fields={s1=100, s2="hello"}`。此时不同参数对应的处理结果为：

   - `Measurements`，处理结果：`root.sg.device, s1=100, s2="hello", tag1="private", tag2="working"`
   - `DeviceID_subtree`，处理结果：`root.sg.device.private.working, s1=100, s2="hello"`

## 测试

本插件自带测试。
