package http

import (
	"errors"
	"io"
	"net/http"
)

// DrainAndCloseRequest can be used (most probably with defer) from the client side to ensure that the http request body is consumed til the end and closed.
func DrainAndCloseRequest(r *http.Request) error {
	if r == nil {
		return nil
	}

	return DrainAndCloseBody(r.Body)
}

// DrainAndCloseResponse can be used (most probably with defer) from the client side to ensure that the http response body is consumed til the end and closed.
func DrainAndCloseResponse(r *http.Response) error {
	if r == nil {
		return nil
	}

	return DrainAndCloseBody(r.Body)
}

func DrainAndCloseBody(body io.ReadCloser) error {
	if body == nil || body == http.NoBody {
		return nil
	}

	_, discardErr := io.Copy(io.Discard, body)

	return errors.Join(discardErr, body.Close())
}
