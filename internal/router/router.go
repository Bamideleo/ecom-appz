package router

import (
	"ecom-appz/internal/handlers"
	"ecom-appz/internal/logger"
	"ecom-appz/internal/middleware"
	"net/http"
)

func New(log *logger.Logger) http.Handler{
	mux := http.NewServeMux()

	// public routes
	mux.HandleFunc("/health", handlers.Health)

	// Versioned API

	v1 := http.NewServeMux()
	v1.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			handlers.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		w.Write([]byte("users endpoint"))
	})

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// Middleware chain
	handler := middleware.Recovery(log)(
		middleware.RequestLogger(log)(
			mux,
		),
	)

	return handler
}