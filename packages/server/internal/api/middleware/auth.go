package middleware

import (
	"context"
	"net/http"

	"github.com/galchammat/kadeem/internal/logging"
)

type contextKey string

const (
	ctxUserID    contextKey = "user_id"
	ctxUserEmail contextKey = "user_email"
	ctxUserRole  contextKey = "user_role"
)

// AuthMiddleware extracts user information from headers set by oauth2-proxy
// Headers: X-Auth-User-Id, X-Auth-Email, X-Auth-Role (from nginx)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-Auth-User-Id")

		if userID == "" {
			logging.Warn("Missing X-Auth-User-Id header", "path", r.URL.Path)
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		email := r.Header.Get("X-Auth-Email")
		role := r.Header.Get("X-Auth-Role")

		ctx := context.WithValue(r.Context(), ctxUserID, userID)
		ctx = context.WithValue(ctx, ctxUserEmail, email)
		ctx = context.WithValue(ctx, ctxUserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts user_id from request context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(ctxUserID).(string)
	return userID, ok
}

// GetUserRole extracts user_role from request context
func GetUserRole(r *http.Request) string {
	role, _ := r.Context().Value(ctxUserRole).(string)
	return role
}

// IsAdmin checks if user has admin role
func IsAdmin(r *http.Request) bool {
	return GetUserRole(r) == "admin"
}
