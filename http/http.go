package http

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

// NewHTTPClient returns a new default http client.
func NewHTTPClient(globalTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: globalTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// DrainAndCloseRequest can be used (most probably with defer) from the client side to ensure that the http request body is consumed til the end and closed.
func DrainAndCloseRequest(req *http.Request) error {
	if req == nil || req.Body == nil || req.Body == http.NoBody {
		return nil
	}

	_, discardErr := io.Copy(io.Discard, req.Body)

	return errors.Join(discardErr, req.Body.Close())
}

// DrainAndCloseResponse can be used (most probably with defer) from the client side to ensure that the http response body is consumed til the end and closed.
func DrainAndCloseResponse(res *http.Response) error {
	if res == nil || res.Body == nil || res.Body == http.NoBody {
		return nil
	}

	_, discardErr := io.Copy(io.Discard, res.Body)

	return errors.Join(discardErr, res.Body.Close())
}
