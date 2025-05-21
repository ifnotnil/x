package http

import (
	"net"
	"net/http"
	"time"
)

func newDefaultDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
}

func newDefaultTransport() *http.Transport {
	if t, is := http.DefaultTransport.(*http.Transport); is {
		return &http.Transport{
			Proxy:                 t.Proxy,
			DialContext:           newDefaultDialer().DialContext,
			ForceAttemptHTTP2:     t.ForceAttemptHTTP2,
			MaxIdleConns:          t.MaxIdleConns,
			IdleConnTimeout:       t.IdleConnTimeout,
			TLSHandshakeTimeout:   t.TLSHandshakeTimeout,
			ExpectContinueTimeout: t.ExpectContinueTimeout,
		}
	}

	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           newDefaultDialer().DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewHTTPClient returns a new default http client.
func NewHTTPClient(globalTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   globalTimeout,
		Transport: newDefaultTransport(),
	}
}
