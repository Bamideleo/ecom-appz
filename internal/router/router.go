package router

import (
	"database/sql"
	"ecom-appz/internal/handlers"
	"ecom-appz/internal/logger"
	"ecom-appz/internal/middleware"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"net/http"
)

func New(log *logger.Logger, db *sql.DB) http.Handler{
	
	mux := http.NewServeMux()
	userRepo := repositories.NewUserRepository(db)
	authHandler := &handlers.AuthHandler{
		UserRepo: userRepo,
	}

	profileHandler := &handlers.ProfileHandler{
		UserRepo: *userRepo,
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

	// Public routes
	v1.Handle("/auth/register", Method(http.MethodPost, 
		http.HandlerFunc(authHandler.Register),
		),
	)
	v1.Handle("/auth/login", Method(http.MethodPost,
		http.HandlerFunc(authHandler.Login),
	),
)


//Admin-only route
v1.Handle(
	"/products/create",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				// http.HandlerFunc(productHandler.CreateProduct),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("admin access granted"))
		}),
			),
		),
	),
)


// User + Admin route
v1.Handle(
	"/profile",
	Method(http.MethodPost,
	middleware.Auth(
		middleware.Authorize(models.RoleUser, models.RoleAdmin)(
			// http.HandlerFunc(orderHandler.CreateOrder),
			http.HandlerFunc(profileHandler.GetProfile),
		),
	),
),
)

v1.Handle(
	"/profile",
	Method(http.MethodPost,
	middleware.Auth(
		middleware.Authorize(models.RoleUser, models.RoleAdmin)(
			// http.HandlerFunc(orderHandler.CreateOrder),
			http.HandlerFunc(profileHandler.UpdateProfile),
		),
	),
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

// func Method(method string, h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != method {
// 			handlers.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
// 			return
// 		}
// 		h(w, r)
// 	}
// }


func Method(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
