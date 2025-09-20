package middleware

import (
	"net/http"
	"time"

	"github.com/ouz/goauthboilerplate/internal/config"
)

func Logging(logger *config.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapper, r)
			duration := time.Since(start)

			entry := logger.With(
				"method", r.Method,
				"status", wrapper.status,
				"path", r.URL.EscapedPath(),
				"ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"duration", duration.String(),
				"size", wrapper.written,
				"referer", r.Referer(),
			)

			switch {
			case wrapper.status >= 500:
				entry.Error("Server error occurred")
			case wrapper.status >= 400:
				entry.Warn("Client error occurred")
			case wrapper.status >= 300:
				entry.Info("Redirection occurred")
			default:
				entry.Info("Request completed successfully")
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	written int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}
