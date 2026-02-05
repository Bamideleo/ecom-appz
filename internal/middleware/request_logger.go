package middleware

import (
	"ecom-appz/internal/logger"
	"net/http"
	"time"
)

func RequestLogger(log *logger.Logger) func(http.Handler) http.Handler{
return func(next http.Handler) http.Handler {
	return  http.HandlerFunc(func( w http.ResponseWriter, r *http.Request){
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Info(
			r.Method + " " +
				r.URL.Path + " completed in " +
				time.Since(start).String(),
		)
	})
}
}