package middleware

import (
	"context"
	"ecom-appz/internal/auth"
	"ecom-appz/internal/handlers"
	"ecom-appz/internal/helper"
	"net/http"
	"strings"
)



func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == "" {
			handlers.RespondError(w, http.StatusUnauthorized, "missing token")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			handlers.RespondError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			handlers.RespondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), helper.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
