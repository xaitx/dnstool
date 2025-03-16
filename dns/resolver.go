package dns

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/xaitx/dnstool/internal"
	"golang.org/x/net/dns/dnsmessage"
)

// ResolverOptions 定义解析器的配置选项
type ResolverOptions struct {
	// 重试次数，默认为3次
	MaxRetries int
	// 重试间隔，默认为1秒
	RetryInterval time.Duration
	// 单次查询超时时间，默认为5秒
	QueryTimeout time.Duration
}

// DefaultResolverOptions 返回默认配置
func DefaultResolverOptions() *ResolverOptions {
	return &ResolverOptions{
		MaxRetries:    3,
		RetryInterval: time.Second,
		QueryTimeout:  5 * time.Second,
	}
}

type Record struct {
	Name  string
	Type  string
	Class string
	TTL   uint32
	Data  string
}

type Resolver interface {
	Lookup(ctx context.Context, domain string, queryType string) ([]Record, error)
	// WithOptions 允许更新解析器选项
	WithOptions(opts *ResolverOptions) Resolver
}

type resolverImpl struct {
	internal *internal.Resolver
	options  *ResolverOptions
}

// NewResolver 创建一个新的解析器实例
func NewResolver(server string, port int) Resolver {
	return &resolverImpl{
		internal: internal.New(server, port),
		options:  DefaultResolverOptions(),
	}
}

// WithOptions 更新解析器选项
func (r *resolverImpl) WithOptions(opts *ResolverOptions) Resolver {
	if opts != nil {
		r.options = opts
	}
	return r
}

func (r *resolverImpl) Lookup(ctx context.Context, domain string, queryType string) ([]Record, error) {
	var records []Record
	var lastErr error

	for attempt := 0; attempt <= r.options.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(r.options.RetryInterval):
			}
		}

		queryCtx, cancel := context.WithTimeout(ctx, r.options.QueryTimeout)
		msg, err := r.internal.Query(queryCtx, domain, queryType)
		cancel()

		if err == nil {
			return parseResponse(msg), nil
		}

		lastErr = err
		// 如果是格式错误或者配置错误，不需要重试
		if _, ok := err.(*internal.DNSFormatError); ok {
			return nil, err
		}
	}

	return records, lastErr
}

func parseResponse(msg *dnsmessage.Message) []Record {

	var data string
	var records []Record
	for _, answer := range msg.Answers {

		// 判断Type
		switch answer.Header.Type {
		case dnsmessage.TypeA:
			if a, ok := answer.Body.(*dnsmessage.AResource); ok {
				data = net.IP(a.A[:]).String()
			}
		case dnsmessage.TypeAAAA:
			if aaaa, ok := answer.Body.(*dnsmessage.AAAAResource); ok {
				data = net.IP(aaaa.AAAA[:]).String()
			}
		case dnsmessage.TypeCNAME:
			if cname, ok := answer.Body.(*dnsmessage.CNAMEResource); ok {
				data = cname.CNAME.String()
			}
		case dnsmessage.TypeMX:
			if mx, ok := answer.Body.(*dnsmessage.MXResource); ok {
				data = mx.MX.String()
			}
		case dnsmessage.TypeNS:
			if ns, ok := answer.Body.(*dnsmessage.NSResource); ok {
				data = ns.NS.String()
			}
		case dnsmessage.TypeTXT:
			if txt, ok := answer.Body.(*dnsmessage.TXTResource); ok {
				// logs.Info(txt.TXT)
				data = strings.Join(txt.TXT, "\n")
			}
		}

		records = append(records, Record{
			Name:  answer.Header.Name.String(),
			Type:  answer.Header.Type.String(),
			Class: answer.Header.Class.String(),
			TTL:   answer.Header.TTL,
			Data:  data,
		})
	}
	return records
}
