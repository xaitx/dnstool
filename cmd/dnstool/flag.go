package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Server    string
	Port      int
	QueryType string
	Domain    string
	Timeout   int
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Server, "server", "8.8.8.8", "dns server,and support doh server")
	flag.IntVar(&cfg.Port, "port", 53, "dns server port,doh server port is not used")
	flag.StringVar(&cfg.QueryType, "type", "A", "can be A, AAAA, CNAME, MX, NS, PTR, SOA, SRV, TXT or ANY. case insensitive")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "dns query timeout")

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
		return nil, fmt.Errorf("please input domain")
	}
	cfg.Domain = flag.Arg(0)

	return cfg, nil
}
