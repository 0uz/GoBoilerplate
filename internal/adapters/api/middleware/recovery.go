package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/ouz/goauthboilerplate/internal/config"
)

func Recovery(logger *config.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := string(debug.Stack())

					entry := logger.With(
						"error", err,
						"stack", stack,
						"path", r.URL.EscapedPath(),
						"method", r.Method,
					)

					entry.Error("panic recovered")

					w.WriteHeader(http.StatusInternalServerError)

					if _, writeErr := w.Write([]byte("Internal Server Error")); writeErr != nil {
						logger.Error("Failed to write error response", "error", writeErr)
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
