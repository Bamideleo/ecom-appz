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

	product := repositories.NewProductRepository(db)
	productHandler := &handlers.ProductHandler{
		Repo: product,
	}
	category := repositories.NewCategoryRepository(db)
	categoryHandler := &handlers.CategoryHandler{
		Repo: category,
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

v1.Handle("/auth/refresh", Method(http.MethodPost,
		http.HandlerFunc(authHandler.Refresh),
	),
)
v1.Handle("/auth/logout", Method(http.MethodPost,
		http.HandlerFunc(authHandler.Logout),
	),
)
v1.Handle("/products", Method(http.MethodPost,
		http.HandlerFunc(productHandler.GetAll),
	),
)

v1.Handle("/products/{id}", Method(http.MethodPost,
		http.HandlerFunc(productHandler.GetByID),
	),
)

v1.Handle("/category", Method(http.MethodPost,
		http.HandlerFunc(categoryHandler.GetAll),
	),
)


//Admin-only route
// product section
v1.Handle(
	"/products/create",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(productHandler.Create)),
			),
		),
	),
)

v1.Handle(
	"/products/update/{id}",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(productHandler.Update)),
			),
		),
	),
)

v1.Handle(
	"/products/delete/{id}",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(productHandler.Delete)),
			),
		),
	),
)

// categories section
v1.Handle(
	"/category/create",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(categoryHandler.Create)),
			),
		),
	),
)

v1.Handle(
	"/category/update/{id}",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(categoryHandler.Update)),
			),
		),
	),
)

v1.Handle(
	"/category/delete/{id}",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(categoryHandler.Delete)),
			),
		),
	),
)

v1.Handle(
	"/category/attach/{productID}",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(categoryHandler.AttachProduct)),
			),
		),
	),
)

v1.Handle(
	"/category/detach/{productID}",
	Method(http.MethodPost,
		middleware.Auth(
			middleware.Authorize(models.RoleAdmin)(
				Method(http.MethodPost, http.HandlerFunc(categoryHandler.DetachProduct)),
			),
		),
	),
)




// End only admin

// User + Admin route
v1.Handle(
    "/profile",
    middleware.Auth(
        middleware.Authorize(models.RoleUser, models.RoleAdmin)(
            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                switch r.Method {
                case http.MethodGet:
                    profileHandler.GetProfile(w, r)
                case http.MethodPost:
                    profileHandler.UpdateProfile(w, r)
                default:
                    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                }
            }),
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


func Method(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
