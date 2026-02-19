package qrlrpc

import (
	"io"
	"net/http"
)

type httpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type logger interface {
	Println(v ...interface{})
}

// WithHttpClient set custom http client
func WithHttpClient(client httpClient) func(rpc *QRLRPC) {
	return func(rpc *QRLRPC) {
		rpc.client = client
	}
}

// WithLogger set custom logger
func WithLogger(l logger) func(rpc *QRLRPC) {
	return func(rpc *QRLRPC) {
		rpc.log = l
	}
}

// WithDebug set debug flag
func WithDebug(enabled bool) func(rpc *QRLRPC) {
	return func(rpc *QRLRPC) {
		rpc.Debug = enabled
	}
}
