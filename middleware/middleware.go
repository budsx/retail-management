package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/budsx/retail-management/utils"
)

type ContextKey string

const (
	ContextKeyUserID   = ContextKey("user_id")
	ContextKeyUsername = ContextKey("username")
)


func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		token := strings.Split(authHeader, "Bearer ")[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

