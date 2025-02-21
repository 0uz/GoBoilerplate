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

					// Log entry oluştur
					entry := logger.WithFields(map[string]any{
						"error":  err,
						"stack":  stack,
						"path":   r.URL.EscapedPath(),
						"method": r.Method,
					})

					entry.Error("panic recovered")

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
