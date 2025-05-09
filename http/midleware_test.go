package http

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecoverer(t *testing.T) {
	tests := map[string]struct {
		Handler         func(http.ResponseWriter, *http.Request)
		ExpectedStatus  int
		ExpectedBodyStr string
		ExpectedLogStr  string
	}{
		"panic string": {
			Handler:         func(_ http.ResponseWriter, _ *http.Request) { panic("panic") },
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedBodyStr: ``,
			ExpectedLogStr:  `{"level":"ERROR", "msg":"recovered from panic", "recover":"panic", "recover_type":"string"}`,
		},
		"panic error": {
			Handler:         func(_ http.ResponseWriter, _ *http.Request) { panic(errors.New("error")) },
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedBodyStr: ``,
			ExpectedLogStr:  `{"level":"ERROR", "msg":"recovered from panic", "recover":"error", "recover_type":"*errors.errorString"}`,
		},
		"http.ErrAbortHandler": {
			Handler:         func(_ http.ResponseWriter, _ *http.Request) { panic(http.ErrAbortHandler) },
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedBodyStr: ``,
			ExpectedLogStr:  ``,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			rr := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			require.NoError(t, err)

			logger, buf := testLogger()

			// subject
			subjectHTTPHandler := Recoverer(logger)(http.HandlerFunc(tc.Handler))

			// run
			subjectHTTPHandler.ServeHTTP(rr, req)

			// verify
			assert.Equal(t, tc.ExpectedStatus, rr.Code)
			if tc.ExpectedBodyStr == "" {
				assert.Equal(t, 0, rr.Body.Len())
			} else {
				assert.JSONEq(t, tc.ExpectedBodyStr, rr.Body.String())
			}

			// verify logs
			if tc.ExpectedLogStr == "" {
				assert.Equal(t, 0, buf.Len())
			} else {
				assert.JSONEq(t, tc.ExpectedLogStr, buf.String())
			}
		})
	}
}

func testLogger() (*slog.Logger, *bytes.Buffer) {
	b := &bytes.Buffer{}
	h := slog.NewJSONHandler(b, &slog.HandlerOptions{AddSource: false, Level: slog.LevelError, ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{}
		}

		return a
	}})
	l := slog.New(h)

	return l, b
}
