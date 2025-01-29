package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(logger *logrus.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(wrapper, r)

			// Log entry oluştur
			entry := logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"status":     wrapper.status,
				"path":       r.URL.EscapedPath(),
				"ip":         r.RemoteAddr,
				"user_agent": r.UserAgent(),
				"duration":   time.Since(start),
			})

			switch {
			case wrapper.status >= 500:
				entry.Error("Request completed")
			case wrapper.status >= 400:
				entry.Warn("Request completed")
			default:
				entry.Info("Request completed")
			}
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
