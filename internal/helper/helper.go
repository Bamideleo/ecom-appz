package helper

import (
	"context"
	"ecom-appz/internal/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func GetUserClaims(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*auth.Claims)
	return claims, ok
}