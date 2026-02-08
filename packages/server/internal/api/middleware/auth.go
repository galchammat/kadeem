package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ctxUserID    contextKey = "user_id"
	ctxUserEmail contextKey = "user_email"
	ctxUserRole  contextKey = "user_role"
)

// CustomClaims represents Supabase JWT claims
type CustomClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates JWT validation middleware for Supabase tokens
// jwtSecret should be the SUPABASE_JWT_SECRET environment variable
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract Bearer token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logging.Warn("Missing Authorization header", "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}

			// Parse "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logging.Warn("Invalid Authorization header format", "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}

			tokenString := parts[1]

			// Parse and validate JWT
			claims := &CustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method is HS256
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			}, jwt.WithValidMethods([]string{"HS256"}))

			if err != nil {
				logging.Warn("JWT parsing failed", "path", r.URL.Path, "error", err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}

			// Validate token is valid
			if !token.Valid {
				logging.Warn("JWT validation failed", "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}

			// Extract user info from claims
			userID := claims.Sub
			if userID == "" {
				logging.Warn("Missing sub claim in JWT", "path", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				return
			}

			// Set context values
			ctx := context.WithValue(r.Context(), ctxUserID, userID)
			ctx = context.WithValue(ctx, ctxUserEmail, claims.Email)
			ctx = context.WithValue(ctx, ctxUserRole, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
