package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					slog.ErrorContext(r.Context(),
						"panic recovered",
						"error", err,
						"stack", debug.Stack(),
						"path", r.URL.EscapedPath(),
						"method", r.Method,
					)

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
