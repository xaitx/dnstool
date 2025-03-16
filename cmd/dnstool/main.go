package main

import (
	"context"
	"os"
	"time"

	"github.com/xaitx/dnstool/dns"
	"github.com/xaitx/logs"
)

func main() {
	logs.SetPrefix(false)
	logs.SetFlags(0)
	cfg, err := ParseFlags()
	if err != nil {
		logs.Error("Error parsing flags:", err)
		os.Exit(1)
	}

	resolver := dns.NewResolver(cfg.Server, cfg.Port)

	// 配置解析器选项
	opts := &dns.ResolverOptions{
		MaxRetries:    cfg.Retries,
		RetryInterval: time.Duration(cfg.RetryInterval) * time.Second,
		QueryTimeout:  time.Duration(cfg.Timeout) * time.Second,
	}
	resolver = resolver.WithOptions(opts)

	ctx := context.Background()
	records, err := resolver.Lookup(ctx, cfg.Domain, cfg.QueryType)
	if err != nil {
		if _, ok := err.(interface{ IsDNSError() bool }); ok {
			logs.Error("DNS查询失败:", err)
		} else {
			logs.Error("发生未知错误:", err)
		}
		os.Exit(1)
	}

	logs.Info("Domain:    ", cfg.Domain)
	logs.Info()
	for _, record := range records {
		logs.Info("Type:    ", record.Type)
		logs.Info("TTL:     ", record.TTL)
		logs.Info("Data:")
		logs.Info(record.Data)
		logs.Info()
	}
}
