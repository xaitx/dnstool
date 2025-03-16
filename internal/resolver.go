package internal

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

type Resolver struct {
	server string
	port   int
}

func New(server string, port int) *Resolver {
	return &Resolver{
		server: server,
		port:   port,
	}
}

func (r *Resolver) Query(ctx context.Context, domain string, queryType string) (*dnsmessage.Message, error) {
	qtype, err := parseQueryType(queryType)
	if err != nil {
		return nil, &DNSFormatError{Message: err.Error()}
	}

	// 判断domain是否以.结尾
	if domain[len(domain)-1] != '.' {
		domain += "."
	}

	name, err := dnsmessage.NewName(domain)
	if err != nil {
		return nil, &DNSFormatError{Message: fmt.Sprintf("invalid domain name: %v", err)}
	}

	msg := dnsmessage.Message{
		Header: dnsmessage.Header{
			RecursionDesired: true,
			ID:               uint16(rand.IntN(65535)),
		},
		Questions: []dnsmessage.Question{
			{
				Name:  name,
				Type:  qtype,
				Class: dnsmessage.ClassINET,
			},
		},
	}

	packed, err := msg.Pack()
	if err != nil {
		return nil, &DNSFormatError{Message: fmt.Sprintf("failed to pack DNS message: %v", err)}
	}

	var msgData []byte

	if strings.HasPrefix(r.server, "https://") {
		msgData, err = queryDoH(r.server, packed)
		if err != nil {
			return nil, &DNSNetworkError{
				Message: "DoH query failed",
				Cause:   err,
			}
		}
	} else {
		conn, err := net.Dial("udp", net.JoinHostPort(r.server, strconv.Itoa(r.port)))
		if err != nil {
			return nil, &DNSNetworkError{
				Message: "failed to connect to DNS server",
				Cause:   err,
			}
		}
		defer conn.Close()

		if deadline, ok := ctx.Deadline(); ok {
			if err := conn.SetDeadline(deadline); err != nil {
				return nil, &DNSNetworkError{
					Message: "failed to set connection deadline",
					Cause:   err,
				}
			}
		}

		if _, err := conn.Write(packed); err != nil {
			return nil, &DNSNetworkError{
				Message: "failed to send DNS query",
				Cause:   err,
			}
		}

		response := make([]byte, 512)
		n, err := conn.Read(response)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, &DNSTimeoutError{Message: "query timed out"}
			}
			return nil, &DNSNetworkError{
				Message: "failed to receive DNS response",
				Cause:   err,
			}
		}
		msgData = response[:n]
	}

	var result dnsmessage.Message
	if err := result.Unpack(msgData); err != nil {
		return nil, &DNSFormatError{Message: fmt.Sprintf("failed to unpack DNS response: %v", err)}
	}

	return &result, nil
}

func parseQueryType(qt string) (dnsmessage.Type, error) {
	qt = strings.ToUpper(qt)

	switch qt {
	case "A":
		return dnsmessage.TypeA, nil
	case "AAAA":
		return dnsmessage.TypeAAAA, nil
	case "MX":
		return dnsmessage.TypeMX, nil
	case "NS":
		return dnsmessage.TypeNS, nil
	case "TXT":
		return dnsmessage.TypeTXT, nil
	case "CNAME":
		return dnsmessage.TypeCNAME, nil
	case "PTR":
		return dnsmessage.TypePTR, nil
	case "ALL":
		return dnsmessage.TypeALL, nil
	default:
		return 0, fmt.Errorf("unsupported query type: %s", qt)
	}
}
