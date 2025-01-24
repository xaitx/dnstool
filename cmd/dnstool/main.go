package main

import (
	"context"
	"log"
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
		// flag.Usage()
		os.Exit(1)
	}

	resolver := dns.NewResolver(cfg.Server, cfg.Port)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	records, err := resolver.Lookup(ctx, cfg.Domain, cfg.QueryType)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	logs.Info("Domain:    ", cfg.Domain)
	logs.Info()
	for _, record := range records {
		logs.Info("Type:    ", record.Type)
		logs.Info("TTL:     ", record.TTL)
		logs.Info("Data:")
		// for _, data := range record.Data {
		// 	logs.Info("  ", data)
		// }
		logs.Info(record.Data)
		logs.Info()

	}
}
