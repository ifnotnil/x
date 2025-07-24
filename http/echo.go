package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func EchoHandler(logger *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			err := DrainAndCloseRequest(r)
			if err != nil {
				logger.ErrorContext(ctx, "error during draining and closing request", slog.String("error", err.Error()))
			}
		}(r.Context())

		b, _ := UnpackRequest(r)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(b)
		if err != nil {
			logger.ErrorContext(r.Context(), "error during json encoding", slog.String("error", err.Error()))
		}
	}
}

func UnpackRequest(r *http.Request) (map[string]any, error) {
	body, bodyErr := UnpackRequestBody(r.Header, r.Body)
	exp := map[string]any{
		"method":           r.Method,
		"proto":            r.Proto,
		"protoMajor":       r.ProtoMajor,
		"protoMinor":       r.ProtoMinor,
		"url":              UnpackURL(r.URL),
		"headers":          UnpackHeaders(r.Header),
		"body":             body,
		"contentLength":    r.ContentLength,
		"transferEncoding": r.TransferEncoding,
		"host":             r.Host,
		"trailer":          UnpackHeaders(r.Trailer),
		"remoteAddr":       r.RemoteAddr,
		"requestURI":       r.RequestURI,
	}

	return exp, bodyErr
}

func UnpackURL(u *url.URL) map[string]any {
	if u == nil {
		return nil
	}

	user := map[string]any(nil)
	if u.User != nil {
		user = map[string]any{}
		p, set := u.User.Password()
		user["username"] = u.User.Username()
		if set {
			user["password"] = p
		}
	}

	return map[string]any{
		"scheme":      u.Scheme,
		"opaque":      u.Opaque,
		"user":        user,
		"host":        u.Host,
		"path":        u.Path,
		"rawPath":     u.RawPath,
		"omitHost":    u.OmitHost,
		"forceQuery":  u.ForceQuery,
		"rawQuery":    u.RawQuery,
		"fragment":    u.Fragment,
		"rawFragment": u.RawFragment,
	}
}

func UnpackHeaders(h http.Header) map[string]string {
	if len(h) == 0 {
		return nil
	}

	exp := map[string]string{}
	for k := range h {
		exp[k] = h.Get(k)
	}

	return exp
}

func UnpackRequestBody(h http.Header, body io.ReadCloser) (any, error) {
	ct := h.Get("Content-Type")

	if strings.HasPrefix(ct, "application/json") {
		b := map[string]any{}
		return b, json.NewDecoder(body).Decode(&b)
	}

	if strings.HasPrefix(ct, "text/") {
		b, err := io.ReadAll(body)
		return string(b), err
	}

	s := strings.Builder{}
	enc := base64.NewEncoder(base64.RawURLEncoding, &s)
	_, err := io.Copy(enc, body)

	return s.String(), err
}
