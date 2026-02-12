package middleware

import (
	"ecom-appz/internal/auth"
	"ecom-appz/internal/handlers"
	"net/http"
)

func Authorize(requiredRoles ...string) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				handlers.RespondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			for _, role := range requiredRoles{
				if claims.Role == role{
					next.ServeHTTP(w, r)
					return 
				}
			}

			handlers.RespondError(w, http.StatusForbidden, "forbidden")
		}) 
	}
}