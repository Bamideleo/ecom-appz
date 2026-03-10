package router

import (
	"database/sql"
	"ecom-appz/internal/handlers"
	"ecom-appz/internal/logger"
	"ecom-appz/internal/middleware"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"net/http"
	"time"
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
	cache := utils.NewInMemoryCache(5 *time.Minute)
	product := repositories.NewProductRepository(db)
	productHandler := &handlers.ProductHandler{
		Repo: product,
		Cache: cache,
	}
	category := repositories.NewCategoryRepository(db)
	categoryHandler := &handlers.CategoryHandler{
		Repo: category,
	}
	cart := repositories.NewCartRepository(db)
	cartHandler := &handlers.CartHandler{
		Repo: cart,
	}
	// checkOut:= repositories.NewCartRepository(db)
	// checkoutHandler:= &handlers.CheckoutHandler{
	// 	Rep
	// } 
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
v1.Handle("/products", Method(http.MethodGet,
		http.HandlerFunc(productHandler.GetAll),
	),
)

v1.Handle("/product/{id}", Method(http.MethodGet,
		http.HandlerFunc(productHandler.GetByID),
	),
)

v1.Handle("/products/list", Method(http.MethodGet,
		http.HandlerFunc(productHandler.List),
	),
)

v1.Handle("/category", Method(http.MethodGet,
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

// cart section
v1.Handle(
	"/cart/add",
	Method(http.MethodPost,
		middleware.Auth(
				Method(http.MethodPost, http.HandlerFunc(cartHandler.AddToCart)),
		),
	),
)


v1.Handle(
	"/cart/update",
	Method(http.MethodPut,
		middleware.Auth(
				Method(http.MethodPut, http.HandlerFunc(cartHandler.UpdateQuantity)),
		),
	),
)

v1.Handle(
	"/cart/remove/{product_id}",
	Method(http.MethodPut,
		middleware.Auth(
				Method(http.MethodPut, http.HandlerFunc(cartHandler.RemoveItem)),
		),
	),
)

// v1.Handle(
// 	"/api/v1/checkout",
// 	Method(http.MethodPost,
// 		middleware.Auth(
// 				Method(http.MethodPost, http.HandlerFunc(checkoutHandler.Checkout)),
// 		),
// 	),
// )

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
