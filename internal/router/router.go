package router

import (
	"database/sql"
	"ecom-appz/internal/handlers"
	"ecom-appz/internal/logger"
	"ecom-appz/internal/middleware"
	"ecom-appz/internal/repositories"
	"net/http"
)

func New(log *logger.Logger, db *sql.DB) http.Handler{
	
	mux := http.NewServeMux()
	userRepo := repositories.NewUserRepository(db)
	authHandler := &handlers.AuthHandler{
		UserRepo: userRepo,
	}
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


	v1.HandleFunc("/auth/register", Method(http.MethodPost, authHandler.Register))
	v1.HandleFunc("/auth/login", Method(http.MethodPost,authHandler.Login))


// Protected route
	v1.Handle("/admin",
	middleware.Auth("admin")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("admin access granted"))
		}),
	),
)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// Middleware chain
	handler := middleware.Recovery(log)(
		middleware.RequestLogger(log)(
			mux,
		),
	)

	return handler
}

func Method(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			handlers.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h(w, r)
	}
}