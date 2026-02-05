package middleware

import (
	"ecom-appz/internal/logger"
	"net/http"
	"runtime/debug"
)

func Recovery(log *logger.Logger) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			defer func(){
				if err := recover(); err != nil{
					log.Error(
						"panic recoverd: " +
						 	string(debug.Stack()),
					)

					http.Error(
						w,
						"Internal Server Error",
						http.StatusInternalServerError,
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}