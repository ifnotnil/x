package http

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnpackRequest(t *testing.T) {
	tests := map[string]struct {
		req           func() (*http.Request, error)
		expected      map[string]any
		expectedError require.ErrorAssertionFunc
	}{
		"simple json get request": {
			req: func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "https://domain.test/get", strings.NewReader(`{"test":"ok"}`))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/json")

				return r, nil
			},
			expected: map[string]any{
				"body": map[string]any{
					"test": "ok",
				},
				"contentLength": int64(13),
				"headers": map[string]string{
					"Content-Type": "application/json",
				},
				"host":             "domain.test",
				"method":           "GET",
				"proto":            "HTTP/1.1",
				"protoMajor":       1,
				"protoMinor":       1,
				"remoteAddr":       "",
				"requestURI":       "",
				"trailer":          map[string]string(nil),
				"transferEncoding": []string(nil),
				"url": map[string]any{
					"forceQuery":  false,
					"fragment":    "",
					"host":        "domain.test",
					"omitHost":    false,
					"opaque":      "",
					"path":        "/get",
					"rawFragment": "",
					"rawPath":     "",
					"rawQuery":    "",
					"scheme":      "https",
					"user":        map[string]any(nil),
				},
			},
			expectedError: require.NoError,
		},
		"simple json post request with url username": {
			req: func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "https://domain.test/post?param1=abc#fag1=123", strings.NewReader(`{"test":"ok"}`))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/json")
				r.URL.User = url.User("user-a")

				return r, nil
			},
			expected: map[string]any{
				"body": map[string]any{
					"test": "ok",
				},
				"contentLength": int64(13),
				"headers": map[string]string{
					"Content-Type": "application/json",
				},
				"host":             "domain.test",
				"method":           "POST",
				"proto":            "HTTP/1.1",
				"protoMajor":       1,
				"protoMinor":       1,
				"remoteAddr":       "",
				"requestURI":       "",
				"trailer":          map[string]string(nil),
				"transferEncoding": []string(nil),
				"url": map[string]any{
					"forceQuery":  false,
					"fragment":    "fag1=123",
					"host":        "domain.test",
					"omitHost":    false,
					"opaque":      "",
					"path":        "/post",
					"rawFragment": "",
					"rawPath":     "",
					"rawQuery":    "param1=abc",
					"scheme":      "https",
					"user": map[string]any{
						"username": "user-a",
					},
				},
			},
			expectedError: require.NoError,
		},
		"simple json post request with url username and password": {
			req: func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "https://domain.test/post?param1=abc#fag1=123", strings.NewReader(`{"test":"ok"}`))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/json; charset=utf-8")
				r.URL.User = url.UserPassword("user-a", "abc-123")

				return r, nil
			},
			expected: map[string]any{
				"body": map[string]any{
					"test": "ok",
				},
				"contentLength": int64(13),
				"headers": map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				"host":             "domain.test",
				"method":           "POST",
				"proto":            "HTTP/1.1",
				"protoMajor":       1,
				"protoMinor":       1,
				"remoteAddr":       "",
				"requestURI":       "",
				"trailer":          map[string]string(nil),
				"transferEncoding": []string(nil),
				"url": map[string]any{
					"forceQuery":  false,
					"fragment":    "fag1=123",
					"host":        "domain.test",
					"omitHost":    false,
					"opaque":      "",
					"path":        "/post",
					"rawFragment": "",
					"rawPath":     "",
					"rawQuery":    "param1=abc",
					"scheme":      "https",
					"user": map[string]any{
						"username": "user-a",
						"password": "abc-123",
					},
				},
			},
			expectedError: require.NoError,
		},
		"simple request text plain": {
			req: func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "https://domain.test/post", strings.NewReader(`abc`))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "text/plain")

				return r, nil
			},
			expected: map[string]any{
				"body":          "abc",
				"contentLength": int64(3),
				"headers": map[string]string{
					"Content-Type": "text/plain",
				},
				"host":             "domain.test",
				"method":           "POST",
				"proto":            "HTTP/1.1",
				"protoMajor":       1,
				"protoMinor":       1,
				"remoteAddr":       "",
				"requestURI":       "",
				"trailer":          map[string]string(nil),
				"transferEncoding": []string(nil),
				"url": map[string]any{
					"forceQuery":  false,
					"fragment":    "",
					"host":        "domain.test",
					"omitHost":    false,
					"opaque":      "",
					"path":        "/post",
					"rawFragment": "",
					"rawPath":     "",
					"rawQuery":    "",
					"scheme":      "https",
					"user":        map[string]any(nil),
				},
			},
			expectedError: require.NoError,
		},
		"simple request binary produces base64": {
			req: func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "https://domain.test/post", bytes.NewReader([]byte{0x1, 0x2, 0x3, 0x4}))
				if err != nil {
					return nil, err
				}

				return r, nil
			},
			expected: map[string]any{
				"body":             "AQID",
				"contentLength":    int64(4),
				"headers":          map[string]string(nil),
				"host":             "domain.test",
				"method":           "POST",
				"proto":            "HTTP/1.1",
				"protoMajor":       1,
				"protoMinor":       1,
				"remoteAddr":       "",
				"requestURI":       "",
				"trailer":          map[string]string(nil),
				"transferEncoding": []string(nil),
				"url": map[string]any{
					"forceQuery":  false,
					"fragment":    "",
					"host":        "domain.test",
					"omitHost":    false,
					"opaque":      "",
					"path":        "/post",
					"rawFragment": "",
					"rawPath":     "",
					"rawQuery":    "",
					"scheme":      "https",
					"user":        map[string]any(nil),
				},
			},
			expectedError: require.NoError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := tc.req()
			require.NoError(t, err)

			b, err := UnpackRequest(r)

			tc.expectedError(t, err)
			assert.Equal(t, tc.expected, b)
		})
	}
}
