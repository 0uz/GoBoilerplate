package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func Logging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapper := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(wrapper, r)

			var logFn func(context.Context, string, ...any)
			switch {
			case wrapper.status >= 500:
				logFn = logger.ErrorContext
			case wrapper.status >= 400:
				logFn = logger.WarnContext
			default:
				logFn = logger.InfoContext
			}

			logFn(r.Context(),
				"Request completed",
				"method", r.Method,
				"status", wrapper.status,
				"path", r.URL.EscapedPath(),
				"ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"duration", time.Since(start),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
