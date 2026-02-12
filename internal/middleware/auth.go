package middleware

import (
	"context"
	"ecom-appz/internal/auth"
	"ecom-appz/internal/handlers"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(next http.Handler) http.Handler{
	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == ""{
			handlers.RespondError(w, http.StatusUnauthorized, "missing token")
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer")

		claims, err := auth.ParseToken(tokenString)

		if err !=nil{
			handlers.RespondError(w, http.StatusUnauthorized, "invalid token")
			return 
		}
		// Inject into request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}