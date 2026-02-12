package middleware

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/galchammat/kadeem/internal/logging"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ctxUserID    contextKey = "user_id"
	ctxUserEmail contextKey = "user_email"
	ctxUserRole  contextKey = "user_role"
)

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

// JWKS represents a JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWKSClient fetches and caches JWKS from Supabase
type JWKSClient struct {
	url     string
	mu      sync.RWMutex
	jwks    *JWKS
	expires time.Time
}

// NewJWKSClient creates a new JWKS client
func NewJWKSClient(url string) *JWKSClient {
	return &JWKSClient{
		url: url,
	}
}

// Fetch retrieves the JWKS from the URL with caching
func (c *JWKSClient) Fetch() (*JWKS, error) {
	c.mu.RLock()
	if c.jwks != nil && time.Now().Before(c.expires) {
		jwks := c.jwks
		c.mu.RUnlock()
		return jwks, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check after acquiring write lock
	if c.jwks != nil && time.Now().Before(c.expires) {
		return c.jwks, nil
	}

	resp, err := http.Get(c.url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS fetch failed with status: %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	c.jwks = &jwks
	c.expires = time.Now().Add(10 * time.Minute)
	return c.jwks, nil
}

// GetKey retrieves a key by its kid
func (c *JWKSClient) GetKey(kid string) (*ecdsa.PublicKey, error) {
	jwks, err := c.Fetch()
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Alg == "ES256" && key.Kty == "EC" {
			return parseECKey(key.X, key.Y)
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

// parseECKey converts base64url-encoded x,y coordinates to an ECDSA public key
func parseECKey(xB64, yB64 string) (*ecdsa.PublicKey, error) {
	xBytes, err := base64.RawURLEncoding.DecodeString(xB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode x: %w", err)
	}

	yBytes, err := base64.RawURLEncoding.DecodeString(yB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode y: %w", err)
	}

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	curve := elliptic.P256()
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

// CustomClaims represents Supabase JWT claims
type CustomClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Validate validates the claims (audience and issuer)
func (c *CustomClaims) Validate() error {
	// Validate issuer
	if c.Issuer != "https://seijlvqsunpbzwuydvze.supabase.co/auth/v1" {
		return fmt.Errorf("invalid issuer: %s", c.Issuer)
	}

	// Validate audience
	hasValidAud := false
	for _, aud := range c.Audience {
		if aud == "authenticated" {
			hasValidAud = true
			break
		}
	}
	if !hasValidAud {
		return fmt.Errorf("invalid audience")
	}

	return nil
}

func writeUnauthorizedResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"}); err != nil {
		logging.Error("Failed to encode unauthorized response", "error", err)
	}
}

// AuthMiddleware creates JWT validation middleware using JWKS
func AuthMiddleware(jwksURL string) func(http.Handler) http.Handler {
	client := NewJWKSClient(jwksURL)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract Bearer token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logging.Warn("Missing Authorization header", "path", r.URL.Path)
				writeUnauthorizedResponse(w)
				return
			}

			// Parse "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logging.Warn("Invalid Authorization header format", "path", r.URL.Path)
				writeUnauthorizedResponse(w)
				return
			}

			tokenString := parts[1]

			// Parse token without verification first to get the kid
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &CustomClaims{})
			if err != nil {
				logging.Warn("JWT parsing failed", "path", r.URL.Path, "error", err.Error())
				writeUnauthorizedResponse(w)
				return
			}

			// Get the key ID from the header
			kid, ok := token.Header["kid"].(string)
			if !ok || kid == "" {
				logging.Warn("Missing kid in JWT header", "path", r.URL.Path)
				writeUnauthorizedResponse(w)
				return
			}

			// Fetch the public key from JWKS
			publicKey, err := client.GetKey(kid)
			if err != nil {
				logging.Warn("Failed to get JWKS key", "path", r.URL.Path, "kid", kid, "error", err.Error())
				writeUnauthorizedResponse(w)
				return
			}

			// Parse and validate JWT with ES256
			claims := &CustomClaims{}
			validatedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method is ES256
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			}, jwt.WithValidMethods([]string{"ES256"}))

			if err != nil {
				logging.Warn("JWT validation failed", "path", r.URL.Path, "error", err.Error())
				writeUnauthorizedResponse(w)
				return
			}

			// Validate token is valid
			if !validatedToken.Valid {
				logging.Warn("JWT validation failed", "path", r.URL.Path)
				writeUnauthorizedResponse(w)
				return
			}

			// Validate claims (issuer and audience)
			if err := claims.Validate(); err != nil {
				logging.Warn("JWT claims validation failed", "path", r.URL.Path, "error", err.Error())
				writeUnauthorizedResponse(w)
				return
			}

			// Extract user info from claims
			userID := claims.Sub
			if userID == "" {
				logging.Warn("Missing sub claim in JWT", "path", r.URL.Path)
				writeUnauthorizedResponse(w)
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
