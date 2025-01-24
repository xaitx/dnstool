# DNSTool

一个简单而功能强大的 DNS 查询工具,支持传统 DNS 和 DoH(DNS over HTTPS)查询。

## 功能特点

- 支持多种 DNS 记录类型查询 (A, AAAA, CNAME, MX, NS, PTR, TXT, ANY)
- 支持传统 DNS 服务器查询
- 支持 DoH(DNS over HTTPS)查询
- 可配置查询超时时间
- 友好的命令行界面

## 安装

```bash
go install github.com/xaitx/dnstool@latest
```

## 使用方法

```bash
dnstool [选项] 域名
```

## 选项

* `-server`: DNS 服务器地址 (默认: "8.8.8.8")
* `-port`: DNS 服务器端口 (默认: 53)
* `-type`: 查询类型 (默认: "A")
* `-timeout`: 查询超时时间(秒) (默认: 5)

## 示例

查询A记录

```bash
dnstool example.com
```

使用 DoH 查询 NS 记录:

```bash
dnstool -type NS -server https://doh.360.cn/dns-query example.com
```

使用指定 DNS 服务器查询:

```bash
dnstool -type NS -server 8.8.8.8 example.com
```