package middleware

import (
	"ecom-appz/internal/auth"
	"ecom-appz/internal/handlers"
	"net/http"
	"strings"
)

func Auth(requiredRole string) func(http.Handler) http.Handler{
	return func(next http.Handler)http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

			header :=r.Header.Get("Authorization")

			if header ==""{
				handlers.RespondError(w, http.StatusUnauthorized, "missing token")
				return
			}
			
			token := strings.TrimPrefix(header, "Bearer")
			claims, err := auth.ParseToken(token)

			if err !=nil{
				handlers.RespondError(w, http.StatusUnauthorized,"invalid token")
				return 
			}

			if requiredRole != "" && claims.Role != requiredRole{
				handlers.RespondError(w,http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}