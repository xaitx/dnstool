package dns

import (
	"context"
	"net"
	"strings"

	"github.com/xaitx/dnstool/internal"
	"golang.org/x/net/dns/dnsmessage"
)

type Record struct {
	Name  string
	Type  string
	Class string
	TTL   uint32
	Data  string
}

type Resolver interface {
	Lookup(ctx context.Context, domain string, queryType string) ([]Record, error)
}

type resolverImpl struct {
	internal *internal.Resolver
}

func NewResolver(server string, port int) Resolver {
	return &resolverImpl{
		internal: internal.New(server, port),
	}
}

func (r *resolverImpl) Lookup(ctx context.Context, domain string, queryType string) ([]Record, error) {
	msg, err := r.internal.Query(ctx, domain, queryType)
	if err != nil {
		return nil, err
	}
	return parseResponse(msg), nil
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
