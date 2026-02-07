package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/galchammat/kadeem/internal/logging"
)

// RecoveryMiddleware recovers from panics and returns 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logging.Error("Panic recovered",
					"error", err,
					"path", r.URL.Path,
					"stack", string(stack),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
