package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
)

// Recoverer returns a recoverer http middleware that logs every panic into the provided slog logger.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(ctx context.Context) {
				if rvr := recover(); rvr != nil {
					// if the handler has already written the header this call is noop.
					w.WriteHeader(http.StatusInternalServerError)

					var recoverType slog.Value
					var recoverVal slog.Value

					recoverType = slog.StringValue(reflect.TypeOf(rvr).String())
					if err, is := rvr.(error); is {
						if errors.Is(err, http.ErrAbortHandler) {
							return
						}
						recoverVal = slog.StringValue(err.Error())
					} else {
						recoverVal = slog.AnyValue(rvr)
					}

					logger.ErrorContext(
						ctx,
						"recovered from panic",
						slog.Attr{Key: "recover_type", Value: recoverType},
						slog.Attr{Key: "recover", Value: recoverVal},
					)
				}

				// Maybe: explore hijacking the connection.
				// if hijacker, ok := w.(http.Hijacker); !ok {
				// 	hijacker.Hijack()
				// }
			}(r.Context())

			next.ServeHTTP(w, r)
		})
	}
}
