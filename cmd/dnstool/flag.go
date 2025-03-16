package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Server        string
	Port          int
	QueryType     string
	Domain        string
	Timeout       int
	Retries       int
	RetryInterval int
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Server, "server", "8.8.8.8", "DNS服务器地址")
	flag.IntVar(&cfg.Port, "port", 53, "DNS服务器端口")
	flag.StringVar(&cfg.QueryType, "type", "A", "查询类型 (A, AAAA, CNAME, MX, NS, PTR, TXT, ALL)")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "查询超时时间(秒)")
	flag.IntVar(&cfg.Retries, "retries", 3, "查询失败重试次数")
	flag.IntVar(&cfg.RetryInterval, "retry-interval", 1, "重试间隔时间(秒)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [Options] domain\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s example.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type NS -server https://doh.360.cn/dns-query xaitx.com\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -type NS -server 8.8.8.8 xaitx.com\n", os.Args[0])
	}

	flag.Parse()

	if flag.NArg() != 1 {
		return nil, fmt.Errorf("请提供要查询的域名")
	}
	cfg.Domain = flag.Arg(0)

	return cfg, nil
}
