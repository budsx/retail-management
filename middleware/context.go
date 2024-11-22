package middleware

import (
	"context"

	"github.com/budsx/retail-management/model"
)

type contextKey string

const (
	userInfoKey contextKey = "userInfo"
)

// SetUserInfoToContext adds user information to the context
func SetUserInfoToContext(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, userInfoKey, user)
}

// HasUserInfo checks if the context contains valid user information
func HasUserInfo(ctx context.Context) bool {
	user := GetUserInfoByContext(ctx)
	return user.UserID != 0 && user.Username != ""
}
