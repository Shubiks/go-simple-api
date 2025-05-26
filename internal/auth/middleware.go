package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey = contextKey("userEmail")
const ContextUserIDKey = contextKey("user_id")
const ContextUserNameKey = contextKey("user_name")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/users") || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		userID := int(claims.UserID) // assuming ValidateToken returns a struct
		fmt.Printf("JWT Middleware: userID extracted from token: %d", userID)
		name := claims.Name

		if userID == 0 {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
		ctx = context.WithValue(ctx, ContextUserNameKey, name)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
