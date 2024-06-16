package tools

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
)

var jwtSecret string

// Initialize the JWT middleware with a secret key
func InitAuth(secret string) {
	jwtSecret = secret
}

// CustomClaims defines custom JWT claims
type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware verifies the JWT token
// AuthMiddleware verifies the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Store the user ID from the token into the context
		ctx := context.WithValue(r.Context(), "user", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
