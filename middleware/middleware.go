package middleware

import (
	"context"
	"encoding/json"
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
			sendErrorResponse(w, http.StatusUnauthorized, "Missing token")
			return
		}

		token := strings.Split(authHeader, "Bearer ")[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type UserInfo struct {
	UserID   int64
	Username string
}

func GetUserInfoByContext(ctx context.Context) UserInfo {
	userID := ctx.Value(ContextKeyUserID).(int64)
	userName := ctx.Value(ContextKeyUsername).(string)
	return UserInfo{
		UserID:   userID,
		Username: userName,
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	}

	json.NewEncoder(w).Encode(response)
}
