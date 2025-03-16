package internal

import "fmt"

// DNSError 定义DNS错误的基础接口
type DNSError interface {
	error
	IsDNSError() bool
}

// DNSFormatError 表示DNS消息格式错误
type DNSFormatError struct {
	Message string
}

func (e *DNSFormatError) Error() string {
	return fmt.Sprintf("DNS format error: %s", e.Message)
}

func (e *DNSFormatError) IsDNSError() bool {
	return true
}

// DNSNetworkError 表示网络相关错误
type DNSNetworkError struct {
	Message string
	Cause   error
}

func (e *DNSNetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("DNS network error: %s (cause: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("DNS network error: %s", e.Message)
}

func (e *DNSNetworkError) IsDNSError() bool {
	return true
}

// DNSTimeoutError 表示查询超时错误
type DNSTimeoutError struct {
	Message string
}

func (e *DNSTimeoutError) Error() string {
	return fmt.Sprintf("DNS timeout error: %s", e.Message)
}

func (e *DNSTimeoutError) IsDNSError() bool {
	return true
}
